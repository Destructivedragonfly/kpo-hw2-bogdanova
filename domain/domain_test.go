package domain

import (
	"errors"
	"testing"
	"time"
)

func testStudent(t *testing.T) Student {
	t.Helper()

	student, err := NewStudent("STU-1", "Alex")
	if err != nil {
		t.Fatal(err)
	}
	return student
}

func testRubric(t *testing.T) (Rubric, Criterion, Criterion) {
	t.Helper()

	rubric := NewRubric()

	exam, err := NewCriterion(
		"CRT-EXAM", "Exam", 100, 40,
	)
	if err != nil {
		t.Fatal(err)
	}

	project, err := NewCriterion(
		"CRT-PROJ", "Project", 50, 60,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(exam); err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(project); err != nil {
		t.Fatal(err)
	}
	return rubric, exam, project
}

// Балл не выходит за допустимые границы
func TestNewScoreGuardsCriterionRange(t *testing.T) {
	max, err := NewMaxScore(20)
	if err != nil {
		t.Fatal(err)
	}

	for _, value := range []float64{-0.01, 20.01} {
		if _, err := NewScore(value, max); !errors.Is(err, ErrScoreOutOfRange) {
			t.Errorf("NewScore(%v): expected ErrScoreOutOfRange, got %v", value, err)
		}
	}

	if _, err := NewScore(0, max); err != nil {
		t.Errorf("zero should be valid: %v", err)
	}

	if _, err := NewScore(20, max); err != nil {
		t.Errorf("maximum should be valid: %v", err)
	}
}

// Нельзя оценивать по чужому критерию
func TestGradeFactoryRejectsCriterionOutsideRubric(t *testing.T) {
	rubric, _, _ := testRubric(t)
	student := testStudent(t)

	_, err := NewGradeFactory().Create(
		"GRD-1",
		student,
		rubric,
		"CRT-UNKNOWN",
		10,
		time.Now(),
		"",
	)

	if !errors.Is(err, ErrCriterionNotInRubric) {
		t.Fatalf("expected ErrCriterionNotInRubric, got %v", err)
	}
}

// Factory отклоняет некорректный балл
func TestGradeFactoryEnforcesScoreRange(t *testing.T) {
	rubric, exam, _ := testRubric(t)
	student := testStudent(t)

	_, err := NewGradeFactory().Create(
		"GRD-1",
		student,
		rubric,
		exam.ID(),
		101,
		time.Now(),
		"",
	)

	if !errors.Is(err, ErrScoreOutOfRange) {
		t.Fatalf("expected ErrScoreOutOfRange, got %v", err)
	}
}

// Правильность итоговой оценки
func TestWeightedFinalGradeUsesNormalizedScoresAndWeights(t *testing.T) {
	rubric, exam, project := testRubric(t)
	student := testStudent(t)

	factory := NewGradeFactory()
	at := time.Date(
		2026, 11, 16, 10, 0, 0, 0, time.UTC,
	)

	examGrade, err := factory.Create(
		"GRD-1", student, rubric, exam.ID(), 80, at, "",
	)
	if err != nil {
		t.Fatal(err)
	}

	projectGrade, err := factory.Create(
		"GRD-2", student, rubric, project.ID(), 45, at, "",
	)
	if err != nil {
		t.Fatal(err)
	}

	service, err := NewFinalGradeService(
		WeightedAverageStrategy{},
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := service.Calculate(
		student.ID(),
		rubric,
		[]Grade{examGrade, projectGrade},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result.Percent() != 86 {
		t.Fatalf("expected 86.00%%, got %.2f%%", result.Percent())
	}
}

// Наличие оценок и отсутствие дубликатов
func TestFinalGradeRequiresGradesAndRejectsDuplicateCriterion(t *testing.T) {
	rubric, exam, _ := testRubric(t)
	student := testStudent(t)

	factory := NewGradeFactory()

	service, err := NewFinalGradeService(
		WeightedAverageStrategy{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.Calculate(
		student.ID(), rubric, nil,
	); !errors.Is(err, ErrNoGrades) {
		t.Fatalf("expected ErrNoGrades, got %v", err)
	}

	first, err := factory.Create(
		"GRD-1", student, rubric, exam.ID(), 70, time.Now(), "",
	)
	if err != nil {
		t.Fatal(err)
	}

	second, err := factory.Create(
		"GRD-2", student, rubric, exam.ID(), 80, time.Now(), "",
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.Calculate(
		student.ID(), rubric, []Grade{first, second},
	); !errors.Is(err, ErrDuplicateGrade) {
		t.Fatalf("expected ErrDuplicateGrade, got %v", err)
	}
}

// Запрет одинаковых названий критериев
func TestRubricRejectsDuplicateCriterionName(t *testing.T) {
	rubric := NewRubric()

	first, err := NewCriterion(
		"CRT-1", "Final exam", 100, 50,
	)
	if err != nil {
		t.Fatal(err)
	}

	second, err := NewCriterion(
		"CRT-2", " final EXAM ", 100, 50,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(first); err != nil {
		t.Fatal(err)
	}

	if err := rubric.AddCriterion(second); !errors.Is(err, ErrDuplicateCriterion) {
		t.Fatalf("expected ErrDuplicateCriterion, got %v", err)
	}
}

package application

import (
	"context"
	"fmt"
	"time"

	"example.com/pos/domain"
)

// Application Layer
// POSFacade координирует пользовательские сценарии
// GoF Pattern: Facade (Structural).
type POSFacade struct {
	students     StudentRepository
	rubric       RubricRepository
	grades       GradeRepository
	ids          IDGenerator
	gradeFactory domain.GradeFactory
	finalGrades  domain.FinalGradeService
}

func NewPOSFacade(
	students StudentRepository,
	rubric RubricRepository,
	grades GradeRepository,
	ids IDGenerator,
	gradeFactory domain.GradeFactory,
	finalGrades domain.FinalGradeService,
) *POSFacade {

	return &POSFacade{
		students:     students,
		rubric:       rubric,
		grades:       grades,
		ids:          ids,
		gradeFactory: gradeFactory,
		finalGrades:  finalGrades,
	}
}

// добавить студента
func (p *POSFacade) AddStudent(
	ctx context.Context,
	name string,
) (domain.Student, error) {

	student, err := domain.NewStudent(
		domain.StudentID(p.ids.NewID("STU")),
		name,
	)

	if err != nil {
		return domain.Student{}, err
	}

	if err := p.students.AddStudent(ctx, student); err != nil {
		return domain.Student{}, fmt.Errorf("save student: %w", err)
	}
	return student, nil
}

// посмотреть студентов
func (p *POSFacade) ListStudents(
	ctx context.Context,
) ([]domain.Student, error) {
	return p.students.Students(ctx)
}

// добавить критерий
func (p *POSFacade) AddCriterion(
	ctx context.Context,
	name string,
	maxScore, weight float64,
) (domain.Criterion, error) {

	rubric, err := p.rubric.LoadRubric(ctx)
	if err != nil {
		return domain.Criterion{}, fmt.Errorf("load rubric: %w", err)
	}

	criterion, err := domain.NewCriterion(
		domain.CriterionID(p.ids.NewID("CRT")),
		name,
		maxScore,
		weight,
	)

	if err != nil {
		return domain.Criterion{}, err
	}

	if err := rubric.AddCriterion(criterion); err != nil {
		return domain.Criterion{}, err
	}

	if err := p.rubric.SaveRubric(ctx, rubric); err != nil {
		return domain.Criterion{}, fmt.Errorf("save rubric: %w", err)
	}
	return criterion, nil
}

// посмотреть критерии
func (p *POSFacade) ListCriteria(
	ctx context.Context,
) ([]domain.Criterion, error) {

	rubric, err := p.rubric.LoadRubric(ctx)
	if err != nil {
		return nil, err
	}
	return rubric.Criteria(), nil
}

// выставить или исправить оценку
func (p *POSFacade) AssignGrade(
	ctx context.Context,
	studentID, criterionID string,
	score float64,
	comment string,
) (domain.Grade, bool, error) {

	student, err := p.students.StudentByID(
		ctx, domain.StudentID(studentID),
	)
	if err != nil {
		return domain.Grade{}, false, err
	}

	rubric, err := p.rubric.LoadRubric(ctx)
	if err != nil {
		return domain.Grade{}, false, err
	}

	existing, err := p.grades.GradesByStudent(ctx, student.ID())
	if err != nil {
		return domain.Grade{}, false, err
	}

	gradeID := domain.GradeID(p.ids.NewID("GRD"))
	corrected := false

	for _, grade := range existing {
		if grade.CriterionID() == domain.CriterionID(criterionID) {
			gradeID = grade.ID()
			corrected = true
			break
		}
	}

	grade, err := p.gradeFactory.Create(
		gradeID,
		student,
		rubric,
		domain.CriterionID(criterionID),
		score,
		time.Now().UTC(),
		comment,
	)

	if err != nil {
		return domain.Grade{}, false, err
	}

	if err := p.grades.SaveGrade(ctx, grade); err != nil {
		return domain.Grade{}, false, fmt.Errorf("save grade: %w", err)
	}
	return grade, corrected, nil
}

// посмотреть оценки студента
func (p *POSFacade) ListGrades(
	ctx context.Context,
	studentID string,
) ([]domain.Grade, error) {

	if _, err := p.students.StudentByID(
		ctx, domain.StudentID(studentID),
	); err != nil {
		return nil, err
	}

	return p.grades.GradesByStudent(
		ctx, domain.StudentID(studentID),
	)
}

// вычислить итоговую оценку
func (p *POSFacade) CalculateFinalGrade(
	ctx context.Context,
	studentID string,
) (domain.FinalGrade, error) {

	student, err := p.students.StudentByID(
		ctx, domain.StudentID(studentID),
	)
	if err != nil {
		return domain.FinalGrade{}, err
	}

	rubric, err := p.rubric.LoadRubric(ctx)
	if err != nil {
		return domain.FinalGrade{}, err
	}

	grades, err := p.grades.GradesByStudent(ctx, student.ID())
	if err != nil {
		return domain.FinalGrade{}, err
	}

	return p.finalGrades.Calculate(
		student.ID(),
		rubric,
		grades,
	)
}

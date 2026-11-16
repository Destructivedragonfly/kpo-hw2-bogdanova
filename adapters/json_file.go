package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"example.com/pos/application"
	"example.com/pos/domain"
)

// JSON хранение и выгрузка
type JSONFileStore struct {
	mx           sync.Mutex
	studentsFile string
	rubricFile   string
	gradesFile   string
}

func NewJSONFileStore(directory string) (*JSONFileStore, error) {
	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, err
	}

	return &JSONFileStore{
		studentsFile: filepath.Join(directory, "students.json"),
		rubricFile:   filepath.Join(directory, "rubric.json"),
		gradesFile:   filepath.Join(directory, "grades.json"),
	}, nil
}

// DTO исключительно для JSON

type studentDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type criterionDTO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	MaxScore float64 `json:"max_score"`
	Weight   float64 `json:"weight"`
}

type gradeDTO struct {
	ID          string  `json:"id"`
	StudentID   string  `json:"student_id"`
	CriterionID string  `json:"criterion_id"`
	Score       float64 `json:"score"`
	Date        string  `json:"date"`
	Comment     string  `json:"comment"`
}

func readJSON[T any](file string, result *T) error {
	data, err := os.ReadFile(file)

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func writeJSON[T any](file string, value T) error {
	data, err := json.MarshalIndent(value, "", "  ")

	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

// реализуем интерфейсы из StudentRepository
func (s *JSONFileStore) AddStudent(
	ctx context.Context,
	student domain.Student,
) error {

	s.mx.Lock()
	defer s.mx.Unlock()
	var students []studentDTO

	if err := readJSON(s.studentsFile, &students); err != nil {
		return err
	}

	students = append(students, studentDTO{
		ID:   string(student.ID()),
		Name: student.Name(),
	})
	return writeJSON(s.studentsFile, students)
}

func (s *JSONFileStore) StudentByID(
	ctx context.Context,
	id domain.StudentID,
) (domain.Student, error) {

	s.mx.Lock()
	defer s.mx.Unlock()
	var students []studentDTO

	if err := readJSON(s.studentsFile, &students); err != nil {
		return domain.Student{}, err
	}

	for _, item := range students {
		if item.ID == string(id) {
			return domain.NewStudent(
				domain.StudentID(item.ID),
				item.Name,
			)
		}
	}
	return domain.Student{}, application.ErrStudentNotFound
}

func (s *JSONFileStore) Students(
	ctx context.Context,
) ([]domain.Student, error) {

	s.mx.Lock()
	defer s.mx.Unlock()
	var stored []studentDTO

	if err := readJSON(s.studentsFile, &stored); err != nil {
		return nil, err
	}

	result := make([]domain.Student, 0, len(stored))

	for _, item := range stored {
		student, err := domain.NewStudent(
			domain.StudentID(item.ID),
			item.Name,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, student)
	}
	return result, nil
}

// реализуем интерфейсы из RubricRepository
func (s *JSONFileStore) LoadRubric(
	ctx context.Context,
) (domain.Rubric, error) {

	s.mx.Lock()
	defer s.mx.Unlock()
	var stored []criterionDTO

	if err := readJSON(s.rubricFile, &stored); err != nil {
		return domain.Rubric{}, err
	}

	rubric := domain.NewRubric()

	for _, item := range stored {
		criterion, err := domain.NewCriterion(
			domain.CriterionID(item.ID),
			item.Name,
			item.MaxScore,
			item.Weight,
		)

		if err != nil {
			return domain.Rubric{}, err
		}

		if err := rubric.AddCriterion(criterion); err != nil {
			return domain.Rubric{}, err
		}
	}
	return rubric, nil
}

func (s *JSONFileStore) SaveRubric(
	ctx context.Context,
	rubric domain.Rubric,
) error {

	s.mx.Lock()
	defer s.mx.Unlock()
	criteria := rubric.Criteria()

	stored := make([]criterionDTO, 0, len(criteria))

	for _, criterion := range criteria {
		stored = append(stored, criterionDTO{
			ID:       string(criterion.ID()),
			Name:     criterion.Name(),
			MaxScore: criterion.MaxScore(),
			Weight:   criterion.Weight(),
		})
	}
	return writeJSON(s.rubricFile, stored)
}

// реализуем интерфейсы из GradeRepository
func (s *JSONFileStore) GradesByStudent(
	ctx context.Context,
	studentID domain.StudentID,
) ([]domain.Grade, error) {

	s.mx.Lock()
	defer s.mx.Unlock()
	var stored []gradeDTO

	if err := readJSON(s.gradesFile, &stored); err != nil {
		return nil, err
	}

	result := make([]domain.Grade, 0)

	for _, item := range stored {
		if item.StudentID != string(studentID) {
			continue
		}

		grade, err := s.restoreGrade(item)
		if err != nil {
			return nil, err
		}

		result = append(result, grade)
	}
	return result, nil
}

func (s *JSONFileStore) SaveGrade(
	ctx context.Context,
	grade domain.Grade,
) error {

	s.mx.Lock()
	defer s.mx.Unlock()
	var stored []gradeDTO

	if err := readJSON(s.gradesFile, &stored); err != nil {
		return err
	}

	found := false

	for i, item := range stored {
		if item.ID == string(grade.ID()) {
			stored[i] = s.gradeToDTO(grade)
			found = true
			break
		}
	}

	if !found {
		stored = append(stored, s.gradeToDTO(grade))
	}
	return writeJSON(s.gradesFile, stored)
}

func (s *JSONFileStore) gradeToDTO(
	grade domain.Grade,
) gradeDTO {

	return gradeDTO{
		ID:          string(grade.ID()),
		StudentID:   string(grade.StudentID()),
		CriterionID: string(grade.CriterionID()),
		Score:       grade.Score(),
		Date:        grade.Date().Format(time.RFC3339),
		Comment:     grade.Comment(),
	}
}

func (s *JSONFileStore) restoreGrade(
	item gradeDTO,
) (domain.Grade, error) {

	rubric, err := s.LoadRubricUnlocked()
	if err != nil {
		return domain.Grade{}, err
	}

	student, err := domain.NewStudent(
		domain.StudentID(item.StudentID),
		"Restored student",
	)
	if err != nil {
		return domain.Grade{}, err
	}

	parsedDate, err := time.Parse(
		time.RFC3339,
		item.Date,
	)
	if err != nil {
		return domain.Grade{}, err
	}

	return domain.NewGradeFactory().Create(
		domain.GradeID(item.ID),
		student,
		rubric,
		domain.CriterionID(item.CriterionID),
		item.Score,
		parsedDate,
		item.Comment,
	)
}

func (s *JSONFileStore) LoadRubricUnlocked() (domain.Rubric, error) {
	var stored []criterionDTO

	if err := readJSON(s.rubricFile, &stored); err != nil {
		return domain.Rubric{}, err
	}

	rubric := domain.NewRubric()

	for _, item := range stored {
		criterion, err := domain.NewCriterion(
			domain.CriterionID(item.ID),
			item.Name,
			item.MaxScore,
			item.Weight,
		)

		if err != nil {
			return domain.Rubric{}, err
		}

		if err := rubric.AddCriterion(criterion); err != nil {
			return domain.Rubric{}, err
		}
	}
	return rubric, nil
}

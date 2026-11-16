package application

import (
	"context"
	"errors"

	"example.com/pos/domain"
)

var ErrStudentNotFound = errors.New("student was not found")

// Ports
// интерфейсы для работы с данными
type StudentRepository interface {
	AddStudent(context.Context, domain.Student) error

	StudentByID(
		context.Context,
		domain.StudentID,
	) (domain.Student, error)

	Students(context.Context) ([]domain.Student, error)
}

type RubricRepository interface {
	LoadRubric(context.Context) (domain.Rubric, error)
	SaveRubric(context.Context, domain.Rubric) error
}

type GradeRepository interface {
	GradesByStudent(
		context.Context,
		domain.StudentID,
	) ([]domain.Grade, error)

	SaveGrade(context.Context, domain.Grade) error
}

type IDGenerator interface {
	NewID(prefix string) string
}

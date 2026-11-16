package domain

import "time"

// Entity: GradeFactory
// создаёт оценку для конкретного AssessmentCriterion, проверяет, что он есть в Rubric
// внутренний конструктор оценки же проверяет допустимый диапазон баллов
type GradeFactory struct{}

func NewGradeFactory() GradeFactory {
	return GradeFactory{}
}

// GoF Pattern: Factory (Creational)
func (GradeFactory) Create(
	id GradeID,
	student Student,
	rubric Rubric,
	criterionID CriterionID,
	rawScore float64,
	gradedAt time.Time,
	comment string,
) (Grade, error) {

	criterion, exists := rubric.FindCriterion(criterionID)
	if !exists {
		return Grade{}, ErrCriterionNotInRubric
	}

	return newGrade(
		id,
		student.ID(),
		criterion,
		rawScore,
		gradedAt,
		comment,
	)
}

package domain

import (
	"fmt"
	"math"
)

// Entity: FinalGrade
// FinalGrade вычисляется из существующих оценок и rubric
type FinalGrade struct {
	percent float64
}

func newFinalGrade(percent float64) (FinalGrade, error) {
	if math.IsNaN(percent) ||
		math.IsInf(percent, 0) ||
		percent < 0 ||
		percent > 100 {
		return FinalGrade{}, ErrInvalidFinalGrade
	}

	return FinalGrade{
		percent: math.Round(percent*100) / 100,
	}, nil
}

func (f FinalGrade) Percent() float64 {
	return f.percent
}

// FinalGradeStrategy позволяет задавать метод вычисления FinalGrade
// GoF Pattern: Strategy (Behavioral)
type FinalGradeStrategy interface {
	Calculate(grades []Grade, rubric Rubric) (float64, error)
}

// WeightedAverageStrategy вычисляет взвешенное среднее
type WeightedAverageStrategy struct{}

func (WeightedAverageStrategy) Calculate(
	grades []Grade,
	rubric Rubric,
) (float64, error) {

	var weightedTotal float64
	var totalWeight float64

	for _, grade := range grades {
		criterion, ok := rubric.FindCriterion(grade.CriterionID())
		if !ok {
			return 0, ErrCriterionNotInRubric
		}

		normalizedScore := grade.Score() / criterion.MaxScore()
		weightedTotal += normalizedScore * criterion.Weight()
		totalWeight += criterion.Weight()
	}

	if totalWeight == 0 {
		return 0, ErrNoGrades
	}
	return (weightedTotal / totalWeight) * 100, nil
}

// Domain Service: FinalGradeService
// так как формирование оценки требует сразу несколько Entity
type FinalGradeService struct {
	strategy FinalGradeStrategy
}

func NewFinalGradeService(
	strategy FinalGradeStrategy,
) (FinalGradeService, error) {

	if strategy == nil {
		return FinalGradeService{},
			fmt.Errorf("final grade strategy is required")
	}

	return FinalGradeService{
		strategy: strategy,
	}, nil
}

func (s FinalGradeService) Calculate(
	studentID StudentID,
	rubric Rubric,
	grades []Grade,
) (FinalGrade, error) {

	if len(grades) == 0 {
		return FinalGrade{}, ErrNoGrades
	}

	seen := make(map[CriterionID]struct{})

	for _, grade := range grades {
		if grade.StudentID() != studentID {
			return FinalGrade{}, ErrGradeForAnotherPerson
		}

		if _, exists := seen[grade.CriterionID()]; exists {
			return FinalGrade{}, ErrDuplicateGrade
		}

		seen[grade.CriterionID()] = struct{}{}

		criterion, ok := rubric.FindCriterion(
			grade.CriterionID(),
		)

		if !ok {
			return FinalGrade{}, ErrCriterionNotInRubric
		}

		if grade.Score() < 0 ||
			grade.Score() > criterion.MaxScore() {
			return FinalGrade{}, ErrScoreOutOfRange
		}
	}
	percent, err := s.strategy.Calculate(grades, rubric)
	if err != nil {
		return FinalGrade{}, err
	}

	return newFinalGrade(percent)
}

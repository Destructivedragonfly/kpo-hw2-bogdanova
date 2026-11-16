package domain

import (
	"math"
	"strings"
)

type CriterionID string

// Entity: AssessmentCriterion
// критерий оценивания
type AssessmentCriterion struct {
	id       CriterionID
	name     string
	maxScore float64 // макс баллы
	weight   float64 // вес критерия
}

func NewAssessmentCriterion(
	id CriterionID,
	name string,
	maxScore float64,
	weight float64,
) (AssessmentCriterion, error) {

	if strings.TrimSpace(string(id)) == "" {
		return AssessmentCriterion{}, ErrInvalidID
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return AssessmentCriterion{}, ErrInvalidName
	}

	if math.IsNaN(maxScore) || math.IsInf(maxScore, 0) || maxScore <= 0 {
		return AssessmentCriterion{}, ErrInvalidMaxScore
	}

	if math.IsNaN(weight) || math.IsInf(weight, 0) || weight <= 0 {
		return AssessmentCriterion{}, ErrInvalidWeight
	}

	return AssessmentCriterion{
		id:       CriterionID(strings.TrimSpace(string(id))),
		name:     name,
		maxScore: maxScore,
		weight:   weight,
	}, nil
}

func (c AssessmentCriterion) ID() CriterionID   { return c.id }
func (c AssessmentCriterion) Name() string      { return c.name }
func (c AssessmentCriterion) MaxScore() float64 { return c.maxScore }
func (c AssessmentCriterion) Weight() float64   { return c.weight }

// Criterion is the assessment criterion used by the application.
type Criterion = AssessmentCriterion

func NewCriterion(
	id CriterionID,
	name string,
	maxScore float64,
	weight float64,
) (Criterion, error) {
	return NewAssessmentCriterion(id, name, maxScore, weight)
}

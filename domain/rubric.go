package domain

import "strings"

// Entity: Rubric
// хранит все AssessmentCriterion в map
type Rubric struct {
	criteria map[CriterionID]AssessmentCriterion
}

func NewRubric() Rubric {
	return Rubric{criteria: make(map[CriterionID]AssessmentCriterion)}
}

func (r *Rubric) AddCriterion(
	criterion AssessmentCriterion,
) error {

	if _, exists := r.criteria[criterion.ID()]; exists {
		return ErrDuplicateCriterion
	}
	for _, existing := range r.criteria {
		if strings.EqualFold(strings.TrimSpace(existing.Name()), strings.TrimSpace(criterion.Name())) {
			return ErrDuplicateCriterion
		}
	}
	r.criteria[criterion.ID()] = criterion
	return nil
}

func (r Rubric) FindCriterion(
	id CriterionID,
) (AssessmentCriterion, bool) {

	criterion, exists := r.criteria[id]
	return criterion, exists
}

func (r Rubric) Criteria() []AssessmentCriterion {
	result := make([]AssessmentCriterion, 0, len(r.criteria))

	for _, criterion := range r.criteria {
		result = append(result, criterion)
	}
	return result
}

package domain

import "errors"

var (
	ErrInvalidID = errors.New(
		"id is empry",
	)

	ErrInvalidName = errors.New(
		"name is empty",
	)

	ErrInvalidMaxScore = errors.New(
		"maximum score should be >0",
	)

	ErrInvalidWeight = errors.New(
		"weight should be >0",
	)

	ErrInvalidScore = errors.New(
		"score should be a valid number",
	)

	ErrScoreOutOfRange = errors.New(
		"score should be in [0, criterion max]",
	)

	ErrCriterionNotInRubric = errors.New(
		"criterion is not part of the rubric",
	)

	ErrDuplicateCriterion = errors.New(
		"criterion already exists in the rubric",
	)

	ErrNoGrades = errors.New(
		"the student without grades",
	)

	ErrGradeForAnotherPerson = errors.New(
		"grade belongs to another student",
	)

	ErrDuplicateGrade = errors.New(
		"more than one current grade exists for a criterion",
	)

	ErrInvalidFinalGrade = errors.New(
		"final grade should be in [0, 100]",
	)
)

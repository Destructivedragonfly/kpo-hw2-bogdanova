package domain

import (
	"math"
	"strings"
	"time"
)

type GradeID string

// Entity: Grade
// конкретная оценка конкретного студента по какому-то критерию
type Grade struct {
	id          GradeID
	studentID   StudentID   // связка со студентом
	criterionID CriterionID // связка с критерием
	score       float64
	date        time.Time
	comment     string
}

func NewGrade(
	id GradeID,
	studentID StudentID,
	criterionID CriterionID,
	score float64,
	comment string,
) (Grade, error) {

	if id == "" || studentID == "" || criterionID == "" {
		return Grade{}, ErrInvalidID
	}

	if math.IsNaN(score) || math.IsInf(score, 0) {
		return Grade{}, ErrInvalidScore
	}

	if score < 0 {
		return Grade{}, ErrInvalidScore
	}

	return Grade{
		id:          id,
		studentID:   studentID,
		criterionID: criterionID,
		score:       score,
		date:        time.Now(),
		comment:     strings.TrimSpace(comment),
	}, nil
}

func (g Grade) ID() GradeID              { return g.id }
func (g Grade) StudentID() StudentID     { return g.studentID }
func (g Grade) CriterionID() CriterionID { return g.criterionID }
func (g Grade) Score() float64           { return g.score }
func (g Grade) Date() time.Time          { return g.date }
func (g Grade) Comment() string          { return g.comment }

func newGrade(
	id GradeID,
	studentID StudentID,
	criterion AssessmentCriterion,
	score float64,
	gradedAt time.Time,
	comment string,
) (Grade, error) {
	if strings.TrimSpace(string(id)) == "" ||
		strings.TrimSpace(string(studentID)) == "" {
		return Grade{}, ErrInvalidID
	}

	if math.IsNaN(score) || math.IsInf(score, 0) {
		return Grade{}, ErrInvalidScore
	}

	if score < 0 || score > criterion.MaxScore() {
		return Grade{}, ErrScoreOutOfRange
	}

	if gradedAt.IsZero() {
		gradedAt = time.Now()
	}

	return Grade{
		id:          id,
		studentID:   studentID,
		criterionID: criterion.ID(),
		score:       score,
		date:        gradedAt,
		comment:     strings.TrimSpace(comment),
	}, nil
}

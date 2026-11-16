package domain

import "math"

func NewMaxScore(value float64) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value <= 0 {
		return 0, ErrInvalidMaxScore
	}
	return value, nil
}

func NewScore(value, maxScore float64) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, ErrInvalidScore
	}
	if value < 0 || value > maxScore {
		return 0, ErrScoreOutOfRange
	}
	return value, nil
}

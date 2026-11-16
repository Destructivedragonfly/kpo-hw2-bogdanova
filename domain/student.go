package domain

import "strings"

// Entity: Student
// StudentID остаётся одним и тем же, остальные детали меняются
type Student struct {
	id   StudentID
	name string
}

type StudentID string

func NewStudent(id StudentID, name string) (Student, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Student{}, ErrInvalidID
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return Student{}, ErrInvalidName
	}
	return Student{id: StudentID(strings.TrimSpace(string(id))), name: name}, nil
}

func (s Student) ID() StudentID { return s.id }
func (s Student) Name() string  { return s.name }

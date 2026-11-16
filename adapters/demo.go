package adapters

import (
	"context"
	"fmt"
	"io"

	"example.com/pos/application"
)

func RunDemo(ctx context.Context, pos *application.POSFacade, output io.Writer) error {
	student, err := pos.AddStudent(ctx, "Alex Morgan")
	if err != nil {
		return err
	}

	exam, err := pos.AddCriterion(ctx, "Exam", 100, 40)
	if err != nil {
		return err
	}

	project, err := pos.AddCriterion(ctx, "Project", 50, 60)
	if err != nil {
		return err
	}

	if _, _, err := pos.AssignGrade(ctx, string(student.ID()), string(exam.ID()), 80, "Exam result"); err != nil {
		return err
	}
	if _, _, err := pos.AssignGrade(ctx, string(student.ID()), string(project.ID()), 45, "Project result"); err != nil {
		return err
	}

	finalGrade, err := pos.CalculateFinalGrade(ctx, string(student.ID()))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(output,
		"Student: %s (ID: %s)\nFinal grade: %.2f%%\n",
		student.Name(), student.ID(), finalGrade.Percent(),
	)
	return err
}

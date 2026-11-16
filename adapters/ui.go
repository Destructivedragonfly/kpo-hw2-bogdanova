package adapters

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.com/pos/application"
)

// ConsoleUI
// управляет взаимодействием с пользователем через консоль
type ConsoleUI struct {
	pos    *application.POSFacade
	reader *bufio.Reader
}

func NewConsoleUI(pos *application.POSFacade) *ConsoleUI {
	return &ConsoleUI{
		pos:    pos,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (ui *ConsoleUI) Run(ctx context.Context) {
	for {
		fmt.Println("\n---------> POS <--------")
		fmt.Println("1. Add student")
		fmt.Println("2. View students")
		fmt.Println("3. Add assessment criterion")
		fmt.Println("4. View criteria")
		fmt.Println("5. Assign or correct grade")
		fmt.Println("6. View student grades")
		fmt.Println("7. Calculate final grade")
		fmt.Println("0. Exit")
		fmt.Print("Choose an option: ")

		choice := ui.readLine()
		var err error

		switch choice {
		case "1":
			err = ui.addStudent(ctx)
		case "2":
			err = ui.listStudents(ctx)
		case "3":
			err = ui.addCriterion(ctx)
		case "4":
			err = ui.listCriteria(ctx)
		case "5":
			err = ui.assignGrade(ctx)
		case "6":
			err = ui.listGrades(ctx)
		case "7":
			err = ui.calculateFinalGrade(ctx)
		case "0":
			fmt.Println("The end!")
			return
		default:
			fmt.Println("Unknown option.")
		}

		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func (ui *ConsoleUI) addStudent(ctx context.Context) error {
	fmt.Print("Student name: ")
	name := ui.readLine()

	student, err := ui.pos.AddStudent(ctx, name)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Student added successfully: %s (ID: %s)\n",
		student.Name(),
		student.ID(),
	)
	return nil
}

func (ui *ConsoleUI) listStudents(ctx context.Context) error {
	students, err := ui.pos.ListStudents(ctx)
	if err != nil {
		return err
	}

	if len(students) == 0 {
		fmt.Println("No students found.")
		return nil
	}

	fmt.Println("\n--> STUDENTS <--")

	for _, student := range students {
		fmt.Printf(
			"ID: %s, Name: %s\n",
			student.ID(),
			student.Name(),
		)
	}
	return nil
}

func (ui *ConsoleUI) addCriterion(ctx context.Context) error {
	fmt.Print("Criterion name: ")
	name := ui.readLine()

	fmt.Print("Maximum score: ")
	maxScore, err := ui.readFloat()
	if err != nil {
		return err
	}

	fmt.Print("Weight: ")
	weight, err := ui.readFloat()
	if err != nil {
		return err
	}

	criterion, err := ui.pos.AddCriterion(
		ctx,
		name,
		maxScore,
		weight,
	)

	if err != nil {
		return err
	}

	fmt.Printf(
		"Criterion added: %s, Maximum: %.2f, Weight: %.2f\n",
		criterion.Name(),
		criterion.MaxScore(),
		criterion.Weight(),
	)
	return nil
}

func (ui *ConsoleUI) listCriteria(ctx context.Context) error {
	criteria, err := ui.pos.ListCriteria(ctx)
	if err != nil {
		return err
	}

	if len(criteria) == 0 {
		fmt.Println("No criteria found.")
		return nil
	}

	fmt.Println("\n--> ASSESSMENT CRITERIA <--")

	for _, criterion := range criteria {
		fmt.Printf(
			"ID: %s, Name: %s, Maximum: %.2f, Weight: %.2f\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore(),
			criterion.Weight(),
		)
	}
	return nil
}

func (ui *ConsoleUI) assignGrade(ctx context.Context) error {
	fmt.Print("Student ID: ")
	studentID := ui.readLine()

	fmt.Print("Criterion ID: ")
	criterionID := ui.readLine()

	fmt.Print("Score: ")
	score, err := ui.readFloat()
	if err != nil {
		return err
	}

	fmt.Print("Comment: ")
	comment := ui.readLine()

	grade, corrected, err := ui.pos.AssignGrade(
		ctx,
		studentID,
		criterionID,
		score,
		comment,
	)

	if err != nil {
		return err
	}

	if corrected {
		fmt.Println("Grade corrected successfully.")
	} else {
		fmt.Println("Grade assigned successfully.")
	}

	fmt.Printf(
		"Grade ID: %s, Score: %.2f\n",
		grade.ID(),
		grade.Score(),
	)
	return nil
}

func (ui *ConsoleUI) listGrades(ctx context.Context) error {
	fmt.Print("Student ID: ")
	studentID := ui.readLine()

	grades, err := ui.pos.ListGrades(ctx, studentID)
	if err != nil {
		return err
	}

	if len(grades) == 0 {
		fmt.Println("This student has no grades yet.")
		return nil
	}

	fmt.Println("\n--> STUDENT GRADES <--")

	for _, grade := range grades {
		fmt.Printf(
			"Criterion ID: %s, Score: %.2f, Date: %s, Comment: %s\n",
			grade.CriterionID(),
			grade.Score(),
			grade.Date().Format("2006-01-02"),
			grade.Comment(),
		)
	}
	return nil
}

func (ui *ConsoleUI) calculateFinalGrade(ctx context.Context) error {
	fmt.Print("Student ID: ")
	studentID := ui.readLine()

	finalGrade, err := ui.pos.CalculateFinalGrade(
		ctx,
		studentID,
	)

	if err != nil {
		return err
	}

	fmt.Printf(
		"\nFinal grade: %.2f%%\n",
		finalGrade.Percent(),
	)
	return nil
}

func (ui *ConsoleUI) readLine() string {
	value, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(value)
}

func (ui *ConsoleUI) readFloat() (float64, error) {
	value := ui.readLine()

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"please enter a correct number: %w",
			err,
		)
	}
	return result, nil
}

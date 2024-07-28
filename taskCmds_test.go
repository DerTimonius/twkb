package main

import (
	"errors"
	"strings"
	"testing"
)

type formTest struct {
	expectedErr error
	name        string
	expected    string
	form        Form
}

type taskTest struct {
	expectedErr error
	name        string
	expected    string
	task        Task
}

func TestAddCmd(t *testing.T) {
	testForm1 := newDefaultForm()
	testForm1.description.SetValue("test the add command")

	testForm2 := newDefaultForm()
	testForm2.description.SetValue("test the add command")
	testForm2.project.SetValue("twkb")

	testForm3 := newDefaultForm()
	testForm3.description.SetValue("test the add command")
	testForm3.project.SetValue("twkb")
	testForm3.label.SetValue("go tui")

	testForm4 := newDefaultForm()
	testForm4.description.SetValue("test the add command")
	testForm4.project.SetValue("twkb")
	testForm4.label.SetValue("go tui")
	testForm4.due.SetValue("7d")

	testForm5 := newDefaultForm()
	testForm5.description.SetValue("test the add command")
	testForm5.label.SetValue("go tui")

	testForm6 := newDefaultForm()
	testForm6.description.SetValue("test the add command")
	testForm6.due.SetValue("eod")

	validTests := []formTest{
		{
			nil,
			"Basic task creation with no label, project or due date",
			"task add test the add command",
			*testForm1,
		},
		{
			nil,
			"Task creation with a project",
			"task add test the add command project:twkb",
			*testForm2,
		},
		{
			nil,
			"Task creation with a project and two tags",
			"task add test the add command project:twkb +go +tui",
			*testForm3,
		},
		{
			nil,
			"Task creation with a project, two tags and a due date",
			"task add test the add command project:twkb due:7d +go +tui",
			*testForm4,
		},
		{
			nil,
			"Task creation only with labels",
			"task add test the add command +go +tui",
			*testForm5,
		},
		{
			nil,
			"Task creation only with due date",
			"task add test the add command due:eod",
			*testForm6,
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := AddCmd(tt.form)
			if result != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}

	errorTestForm := newDefaultForm()
	errorTests := []formTest{
		{
			errors.New("cannot create a task without a description"),
			"Form without description",
			"",
			*errorTestForm,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AddCmd(tt.form)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestStartCmd(t *testing.T) {
	validTests := []taskTest{
		{
			nil,
			"Basic task",
			"task 42 start",
			Task{id: 42, description: "a basic task"},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := StartCmd(&tt.task)
			if strings.Join(result, " ") != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.task, result, tt.expected)
			}
		})
	}

	errorTests := []taskTest{
		{
			errors.New("cannot start a task with ID 0"),
			"Zero ID task",
			"",
			Task{id: 0, description: "an invalid task"},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StartCmd(&tt.task)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestStopCmd(t *testing.T) {
	validTests := []taskTest{
		{
			nil,
			"Basic task",
			"task 42 stop",
			Task{id: 42, description: "a basic task"},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := StopCmd(&tt.task)
			if strings.Join(result, " ") != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.task, result, tt.expected)
			}
		})
	}

	errorTests := []taskTest{
		{
			errors.New("cannot stop a task with ID 0"),
			"Zero ID task",
			"",
			Task{id: 0, description: "an invalid task"},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StopCmd(&tt.task)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestDoneCmd(t *testing.T) {
	validTests := []taskTest{
		{
			nil,
			"Basic task",
			"task 42 done",
			Task{id: 42, description: "a basic task"},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := DoneCmd(&tt.task)
			if strings.Join(result, " ") != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.task, result, tt.expected)
			}
		})
	}

	errorTests := []taskTest{
		{
			errors.New("cannot finish a task with ID 0"),
			"Zero ID task",
			"",
			Task{id: 0, description: "an invalid task"},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DoneCmd(&tt.task)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestDeleteCmd(t *testing.T) {
	validTests := []taskTest{
		{
			nil,
			"Basic task",
			"task rc.confirmation=no 42 delete",
			Task{id: 42, description: "a basic task"},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := DeleteCmd(&tt.task)
			if strings.Join(result, " ") != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.task, result, tt.expected)
			}
		})
	}

	errorTests := []taskTest{
		{
			errors.New("cannot delete a task with ID 0"),
			"Zero ID task",
			"",
			Task{id: 0, description: "an invalid task"},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeleteCmd(&tt.task)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

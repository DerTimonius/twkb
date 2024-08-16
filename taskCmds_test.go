package main

import (
	"errors"
	"strings"
	"testing"
)

type modifyTest struct {
	expectedErr error
	name        string
	expected    string
	task        Task
	form        TaskForm
}

type formTest struct {
	expectedErr error
	name        string
	expected    string
	form        TaskForm
}

type taskTest struct {
	expectedErr error
	name        string
	expected    string
	task        Task
}

type blockTest struct {
	expectedErr error
	name        string
	expected    string
	blocked     []Task
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

	testForm7 := newDefaultForm()
	testForm7.description.SetValue("test the add command")
	testForm7.due.SetValue("eow")
	testForm7.recur.SetValue("monthly")
	testForm7.until.SetValue("now+1yr")

	testForm8 := newDefaultForm()
	testForm8.description.SetValue("test the add command")
	testForm8.label.SetValue("go tui")
	testForm8.due.SetValue("eow")
	testForm8.recur.SetValue("monthly")
	testForm8.until.SetValue("now+1yr")

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
		{
			nil,
			"Task creation with recur and until",
			"task add test the add command due:eow recur:monthly until:now+1yr",
			*testForm7,
		},
		{
			nil,
			"Task creation with recur and until with tags",
			"task add test the add command due:eow +go +tui recur:monthly until:now+1yr",
			*testForm8,
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := AddCmd(tt.form)
			if strings.Join(result, " ") != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}

	errorTestForm := newDefaultForm()

	errorTestForm2 := newDefaultForm()
	errorTestForm2.description.SetValue("invalid test")
	errorTestForm2.recur.SetValue("monthly")

	errorTests := []formTest{
		{
			errors.New("cannot create a task without a description"),
			"Form without description",
			"",
			*errorTestForm,
		},
		{
			errors.New("cannot create a recurring task without a due date"),
			"Recurring task without a due date",
			"",
			*errorTestForm2,
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
			"task rc.confirmation=no 42 done",
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

func TestModifyCmd(t *testing.T) {
	testForm1 := newDefaultForm()
	testForm1.description.SetValue("test the modify command")
	testForm1.project.SetValue("task-gui")

	testForm2 := newDefaultForm()
	testForm2.project.SetValue("twkb")

	testForm3 := newDefaultForm()
	testForm3.label.SetValue("go tui")

	testForm4 := newDefaultForm()
	testForm4.description.SetValue("test the modify command")
	testForm4.project.SetValue("twkb")
	testForm4.label.SetValue("go tui")
	testForm4.due.SetValue("7d")

	testForm5 := newDefaultForm()
	testForm5.description.SetValue("test the modify command")
	testForm5.label.SetValue("go tui")
	testForm5.project.SetValue("task-gui")

	testForm6 := newDefaultForm()
	testForm6.due.SetValue("eow")
	testForm6.project.SetValue("task-gui")

	baseTask := Task{id: 42, description: "basic task", project: "task-gui", tags: []string{"rust", "cli"}, due: "eod"}

	validTests := []modifyTest{
		{
			nil,
			"Modify only the description",
			"task rc.confirmation=no 42 modify test the modify command",
			baseTask,
			*testForm1,
		},
		{
			nil,
			"Modify only the project",
			"task rc.confirmation=no 42 modify project:twkb",
			baseTask,
			*testForm2,
		},
		{
			nil,
			"Remove the project and modify the labels",
			"task rc.confirmation=no 42 modify project: +go +tui -rust -cli",
			baseTask,
			*testForm3,
		},
		{
			nil,
			"Modify every aspect of the task",
			"task rc.confirmation=no 42 modify test the modify command project:twkb due:7d +go +tui -rust -cli",
			baseTask,
			*testForm4,
		},
		{
			nil,
			"Modify the description and the labels",
			"task rc.confirmation=no 42 modify test the modify command +go +tui -rust -cli",
			baseTask,
			*testForm5,
		},
		{
			nil,
			"Modify only the due date",
			"task rc.confirmation=no 42 modify due:eow",
			baseTask,
			*testForm6,
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := ModifyCmd(tt.task, &tt.form)
			if strings.Join(result, " ") != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}

	errorTestForm := newDefaultForm()
	errorTests := []modifyTest{
		{
			errors.New("cannot modify a task with ID 0"),
			"Form without description",
			"",
			Task{id: 0},
			*errorTestForm,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ModifyCmd(tt.task, &tt.form)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestBlockCmd(t *testing.T) {
	validTests := []blockTest{
		{
			nil,
			"Block single task",
			"task 23 modify depends:42",
			[]Task{{id: 23, description: "a blocked task"}},
			Task{id: 42, description: "a basic task"},
		},
		{
			nil,
			"Block multiple tasks",
			"task 23,4,8 modify depends:42",
			[]Task{{id: 23, description: "a blocked task"}, {id: 4, description: "a blocked task"}, {id: 8, description: "a blocked task"}},
			Task{id: 42, description: "a basic task"},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := BlockCmd(&tt.task, &tt.blocked)
			if strings.Join(result, " ") != tt.expected {
				t.Errorf("StartCmd(%v) = %q, want %q", tt.task, result, tt.expected)
			}
		})
	}

	errorTests := []blockTest{
		{
			errors.New("blocking task cannot have ID 0"),
			"Blocking task has ID 0",
			"",
			[]Task{{id: 23, description: "a blocked task"}, {id: 4, description: "a blocked task"}, {id: 8, description: "a blocked task"}},
			Task{id: 0, description: "an invalid task"},
		},
		{
			errors.New("cannot block a task with ID 0"),
			"Blocked task has ID 0",
			"",
			[]Task{{id: 23, description: "a blocked task"}, {id: 0, description: "an invalid task"}, {id: 8, description: "a blocked task"}},
			Task{id: 42, description: "a basic task"},
		},
		{
			errors.New("cannot block a task with same ID"),
			"Blocked task has same ID as blocking task",
			"",
			[]Task{{id: 23, description: "a blocked task"}, {id: 42, description: "an invalid task"}, {id: 8, description: "a blocked task"}},
			Task{id: 42, description: "a basic task"},
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BlockCmd(&tt.task, &tt.blocked)
			if err == nil {
				t.Fatal("Expected an error, but got nil")
			}
			if err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

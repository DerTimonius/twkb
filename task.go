package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Task struct {
	description string
	uuid        string
	start       string
	modified    string
	project     string
	due         string
	tags        []string
	status      status
	id          int
	urgency     float64
	blocked     bool
	recurring   bool
}

func (t *Task) StartStop() {
	// don't start the task when it is blocked
	if t.blocked {
		return
	}
	if t.status == inProgress {
		cmdStr, err := StopCmd(t)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		t.status = todo
	} else {
		cmdStr, err := StartCmd(t)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		t.status = inProgress
	}
}

func (t *Task) Finish() {
	cmdStr, err := DoneCmd(t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	t.status = done
}

func (t *Task) Delete() {
	cmdStr, err := DeleteCmd(t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	t.status = never
}

func (t Task) ModifyTask(f *TaskForm) Task {
	cmdStr, err := ModifyCmd(t, f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	if f.description.Value() != "" && f.description.Value() != t.description {
		t.description = f.description.Value()
	}

	if f.due.Value() != "" && f.due.Value() != t.due {
		t.due = f.due.Value()
	}

	if f.project.Value() != t.project {
		t.project = f.project.Value()
	}

	if f.label.Value() != "" {
		formValues := strings.Split(f.label.Value(), " ")
		addedLabels := []string{}
		currLabels := t.tags
		for _, label := range formValues {
			idx := slices.Index(currLabels, label)
			if idx == -1 {
				addedLabels = append(addedLabels, label)
			} else {
				currLabels = slices.Delete(currLabels, idx, idx+1)
			}
		}
		t.tags = addedLabels
	}

	return t
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.description
}

func (t Task) Title() string {
	var addMsg string
	if t.blocked {
		addMsg = "[BLOCKED]"
	}
	if t.recurring {
		addMsg = "[RECURRING]"
	}
	return fmt.Sprintf("%s %s", t.description, addMsg)
}

func (t Task) Description() string {
	var tagsMsg string
	if len(t.tags) > 0 {
		tagsMsg = fmt.Sprintf("Tags: %s, ", t.tags)
	}
	var projectMsg string
	if t.project != "" {
		projectMsg = fmt.Sprintf("Project: %s, ", t.project)
	}
	var dueMsg string
	if t.due != "" {
		dueMsg = fmt.Sprintf("Due: %s, ", t.due)
	}
	return fmt.Sprintf("%s%s%sUrgency: %.1f", projectMsg, tagsMsg, dueMsg, t.urgency)
}

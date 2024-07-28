package main

import (
	"errors"
	"fmt"
	"strings"
)

func AddCmd(f Form) (string, error) {
	if f.description.Value() == "" {
		return "", errors.New("cannot create a task without a description")
	}
	var due string
	var project string
	var tags string

	if f.due.Value() != "" {
		due = fmt.Sprintf("due:%s ", f.due.Value())
	}

	if f.project.Value() != "" {
		project = fmt.Sprintf("project:%s ", f.project.Value())
	}

	if f.label.Value() != "" {
		labels := strings.Split(f.label.Value(), " ")
		var labelStrings []string
		for _, l := range labels {
			labelStrings = append(labelStrings, fmt.Sprintf("+%s", l))
		}
		tags = strings.Join(labelStrings, " ")
	}

	str := fmt.Sprintf("task add %s %s%s%s", f.description.Value(), project, due, tags)
	return strings.TrimSuffix(str, " "), nil
}

func StartCmd(t *Task) ([]string, error) {
	if t.id == 0 {
		return []string{}, errors.New("cannot start a task with ID 0")
	}
	return []string{"task", fmt.Sprint(t.id), "start"}, nil
}

func StopCmd(t *Task) ([]string, error) {
	if t.id == 0 {
		return []string{}, errors.New("cannot stop a task with ID 0")
	}
	return []string{"task", fmt.Sprint(t.id), "stop"}, nil
}

func DoneCmd(t *Task) ([]string, error) {
	if t.id == 0 {
		return []string{}, errors.New("cannot finish a task with ID 0")
	}
	return []string{"task", fmt.Sprint(t.id), "done"}, nil
}

func DeleteCmd(t *Task) ([]string, error) {
	if t.id == 0 {
		return []string{}, errors.New("cannot delete a task with ID 0")
	}
	return []string{"task", "rc.confirmation=no", fmt.Sprint(t.id), "delete"}, nil
}

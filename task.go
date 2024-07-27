package main

import "fmt"

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
}

func (t *Task) StartStop() {
	if t.status == inProgress {
		t.status = todo
	} else {
		t.status = inProgress
	}
}

func (t *Task) Finish() {
	t.status = done
}

func (t *Task) Delete() {
	t.status = never
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.description
}

func (t Task) Title() string {
	return t.description
}

func (t Task) Description() string {
	return fmt.Sprintf("Project: %s, Tags: %s, Due: %s, Urgency: %.1f", t.project, t.tags, t.due, t.urgency)
}

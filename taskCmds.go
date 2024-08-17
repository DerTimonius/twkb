package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func AddCmd(f TaskForm) ([]string, error) {
	if f.description.Value() == "" {
		return []string{}, errors.New("cannot create a task without a description")
	}
	var due string
	var project string
	var tags string
	var recur string
	var until string

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
		labelStrings = append(labelStrings, "")
		tags = strings.Join(labelStrings, " ")
	}

	if f.recur.Value() != "" {
		recur = fmt.Sprintf("recur:%s ", f.recur.Value())
	}

	if f.until.Value() != "" {
		until = fmt.Sprintf("until:%s ", f.until.Value())
	}

	if recur != "" && due == "" {
		return []string{}, errors.New("cannot create a recurring task without a due date")
	}

	str := fmt.Sprintf("task add %s %s%s%s%s%s", f.description.Value(), project, due, tags, recur, until)
	return strings.Split(strings.TrimSuffix(str, " "), " "), nil
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
	return []string{"task", "rc.confirmation=no", fmt.Sprint(t.id), "done"}, nil
}

func DeleteCmd(t *Task) ([]string, error) {
	if t.id == 0 {
		return []string{}, errors.New("cannot delete a task with ID 0")
	}
	return []string{"task", "rc.confirmation=no", fmt.Sprint(t.id), "delete"}, nil
}

func ModifyCmd(t Task, f *TaskForm) ([]string, error) {
	if t.id == 0 {
		return []string{}, errors.New("cannot modify a task with ID 0")
	}

	var changedLabels []string
	var changedDescription string
	var changedDue string
	var changedProject string

	if f.description.Value() != "" && f.description.Value() != t.description {
		changedDescription = f.description.Value()
	}

	if f.due.Value() != "" && f.due.Value() != t.due {
		changedDue = fmt.Sprintf("due:%s ", f.due.Value())
	}

	if f.project.Value() != t.project {
		changedProject = fmt.Sprintf("project:%s ", f.project.Value())
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

		for _, label := range addedLabels {
			changedLabels = append(changedLabels, fmt.Sprintf("+%s", label))
		}

		for _, label := range currLabels {
			changedLabels = append(changedLabels, fmt.Sprintf("-%s", label))
		}
	}

	str := fmt.Sprintf(
		"task rc.confirmation=no %d modify %s %s%s%s",
		t.id,
		changedDescription,
		changedProject,
		changedDue,
		strings.Join(changedLabels, " "),
	)
	cmdArgs := []string{}
	// remove the nil values if there are any present
	for _, arg := range strings.Split(strings.TrimSuffix(str, " "), " ") {
		if arg == "" {
			continue
		}
		cmdArgs = append(cmdArgs, arg)
	}
	return cmdArgs, nil
}

func BlockCmd(t *Task, blocked *[]Task) ([]string, error) {
	if len(*blocked) == 0 {
		return []string{}, errors.New("need to select at least 1 task")
	}

	if t.id == 0 {
		return []string{}, errors.New("blocking task cannot have ID 0")
	}

	cmd := []string{"task"}

	var taskIds []string
	for _, task := range *blocked {
		if task.id == 0 {
			return []string{}, errors.New("cannot block a task with ID 0")
		}

		if task.id == t.id {
			return []string{}, errors.New("cannot block a task with same ID")
		}

		taskIds = append(taskIds, fmt.Sprint(task.id))
	}
	cmd = append(cmd, strings.Join(taskIds, ","))
	cmd = append(cmd, "modify")
	cmd = append(cmd, fmt.Sprintf("depends:%d", t.id))
	return cmd, nil
}

func UnblockCmd(t *Task) ([]string, error) {
	if t.id == 0 {
		return []string{}, errors.New("cannot unblock task with ID 0")
	}
	if !t.blocked {
		return []string{}, errors.New("cannot unblock a task that is not blocked")
	}

	return []string{"task", fmt.Sprint(t.id), "modify", "depends:"}, nil
}

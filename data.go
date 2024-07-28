package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

func (b *Board) initLists() {
	tasks := getFromTW()
	var todoTasks []Task
	var doingTasks []Task
	var doneTasks []Task
	// var neverTasks []TwTask

	for _, t := range tasks {
		switch t.status {
		case done:
			doneTasks = append(doneTasks, t)
		case inProgress:
			doingTasks = append(doingTasks, t)
		case todo:
			todoTasks = append(todoTasks, t)
			// default:
			// 	neverTasks = append(neverTasks, t)
		}
	}

	sortTasks(todoTasks)
	sortTasks(doingTasks)

	// TODO: add a never column
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}

	// Init To Do
	b.cols[todo].list.Title = "To Do"
	b.cols[todo].list.SetItems(convertToListItems(todoTasks))
	// Init in progress
	b.cols[inProgress].list.Title = "In Progress"
	b.cols[inProgress].list.SetItems(convertToListItems(doingTasks))
	// Init done
	b.cols[done].list.Title = "Done"
	b.cols[done].list.SetItems(convertToListItems(doneTasks))
}

func getFromTW() []Task {
	cmd := exec.Command("task", "export")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	output := out.String()
	var result []map[string]interface{}
	err = json.Unmarshal([]byte(output), &result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var tasks []Task

	for _, v := range result {
		task := Task{}
		if start, ok := v["start"].(string); ok {
			task.start = start
		}
		if modified, ok := v["modified"].(string); ok {
			task.modified = modified
		}
		if uuid, ok := v["uuid"].(string); ok {
			task.uuid = uuid
		}
		if description, ok := v["description"].(string); ok {
			task.description = description
		}
		if status, ok := v["status"].(string); ok {
			if status == "completed" {
				task.status = done
			} else if status == "pending" && task.start != "" {
				task.status = inProgress
			} else if status == "deleted" {
				task.status = never
			} else {
				task.status = todo
			}
		}
		if due, ok := v["due"].(string); ok {
			t, _ := time.Parse("20060102T150405Z", due)
			task.due = fmt.Sprintf("%.1fd", time.Until(t).Hours()/24)
		}
		if project, ok := v["project"].(string); ok {
			task.project = project
		}
		if id, ok := v["id"].(float64); ok {
			task.id = int(id)
		}
		if _, ok := v["depends"].([]interface{}); ok {
			task.blocked = true
		}
		if urgency, ok := v["urgency"].(float64); ok {
			task.urgency = urgency
		}
		if tags, ok := v["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if t, ok := tag.(string); ok {
					task.tags = append(task.tags, t)
				}
			}
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func convertToListItems(tasks []Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	return items
}

func sortTasks(tasks []Task) {
	slices.SortFunc(tasks, func(a, b Task) int {
		return cmp.Compare(a.urgency, b.urgency) * -1
	})
}

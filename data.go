package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
)

// Provides the mock data to fill the kanban board

func (b *Board) initLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}
	// Init To Do
	b.cols[todo].list.Title = "To Do"
	b.cols[todo].list.SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	b.cols[inProgress].list.Title = "In Progress"
	b.cols[inProgress].list.SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init done
	b.cols[done].list.Title = "Done"
	b.cols[done].list.SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
}

type twstatus int

const (
	pending twstatus = iota
	completed
)

type TwTask struct {
	description string
	uuid        string
	start       string
	modified    string
	status      string
	project     string
	due         string
	tags        []string
	id          int
	urgency     float64
}

func (b *Board) getFromTW() {
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

	var tasks []TwTask

	for _, v := range result {
		task := TwTask{}
		if uuid, ok := v["uuid"].(string); ok {
			task.uuid = uuid
		}
		if description, ok := v["description"].(string); ok {
			task.description = description
		}
		if status, ok := v["status"].(string); ok {
			task.status = status
		}
		if modified, ok := v["modified"].(string); ok {
			task.modified = modified
		}
		if due, ok := v["due"].(string); ok {
			task.due = due
		}
		if start, ok := v["start"].(string); ok {
			task.start = start
		}
		if project, ok := v["project"].(string); ok {
			task.project = project
		}
		if id, ok := v["id"].(float64); ok {
			task.id = int(id)
		}
		if urgency, ok := v["urgency"].(float64); ok {
			task.urgency = urgency
		}
		if tags, ok := v["tags"].([]string); ok {
			task.tags = tags
		}
		tasks = append(tasks, task)
	}

	var todoTasks []TwTask
	var doingTasks []TwTask
	var doneTasks []TwTask

	for _, t := range tasks {
		if t.status == "completed" || t.status == "deleted" {
			doneTasks = append(doneTasks, t)
			continue
		}
		if t.status == "pending" && t.start != "" {
			doingTasks = append(doingTasks, t)
			continue
		}
		if t.status == "pending" {
			todoTasks = append(todoTasks, t)
			continue
		}
	}

	fmt.Printf("todoTasks: %v\n", todoTasks)
	fmt.Printf("doingTasks: %v\n", doingTasks)
	fmt.Printf("doneTasks: %v\n", doneTasks)
}

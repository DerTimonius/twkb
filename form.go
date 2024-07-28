package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	help        help.Model
	description textinput.Model
	project     textinput.Model
	label       textinput.Model
	due         textinput.Model
	col         column
	index       int
}

type EditForm struct {
	relatedTask Task
	form        *Form
}

func newDefaultForm() *Form {
	return NewForm("task name", "project", "labels", "due")
}

func NewForm(description, project, label, due string) *Form {
	form := Form{
		help:        help.New(),
		description: textinput.New(),
		project:     textinput.New(),
		label:       textinput.New(),
		due:         textinput.New(),
	}
	form.description.Placeholder = description
	form.project.Placeholder = project
	form.label.Placeholder = label
	form.due.Placeholder = due
	form.description.Focus()
	return &form
}

func (f Form) CreateTask() Task {
	cmdStr, err := AddCmd(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return Task{status: todo, description: f.description.Value(), project: f.project.Value(), tags: strings.Split(f.label.Value(), " ")}
}

func NewEditForm(t Task) EditForm {
	form := Form{
		help:        help.New(),
		description: textinput.New(),
		project:     textinput.New(),
		label:       textinput.New(),
		due:         textinput.New(),
	}
	form.description.SetValue(t.description)
	form.project.SetValue(t.project)
	form.label.SetValue(strings.Join(t.tags, " "))
	form.due.SetValue(t.due)
	form.description.Focus()
	return EditForm{form: &form, relatedTask: t}
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case column:
		f.col = msg
		f.col.list.Index()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return f, tea.Quit

		case key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Enter):
			return board.Update(f)
		case key.Matches(msg, keys.Tab):
			if f.description.Focused() {
				f.description.Blur()
				f.project.Focus()
				return f, textarea.Blink
			}
			if f.project.Focused() {
				f.project.Blur()
				f.label.Focus()
				return f, textarea.Blink
			}
			if f.label.Focused() {
				f.label.Blur()
				f.due.Focus()
				return f, textarea.Blink
			}
			if f.due.Focused() {
				f.due.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}
		}
	}
	if f.description.Focused() {
		f.description, cmd = f.description.Update(msg)
		return f, cmd
	}
	if f.project.Focused() {
		f.project, cmd = f.project.Update(msg)
		return f, cmd
	}
	if f.label.Focused() {
		f.label, cmd = f.label.Update(msg)
		return f, cmd
	}
	f.due, cmd = f.due.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create a new task",
		f.description.View(),
		f.project.View(),
		f.label.View(),
		f.due.View(),
		f.help.View(keys))
}

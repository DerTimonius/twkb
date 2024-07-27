package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	description textinput.Model
	project     textinput.Model
	label       textinput.Model
	due         textinput.Model
	help        help.Model
	col         column
	index       int
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
	return Task{status: todo, description: f.description.Value(), project: f.project.Value()}
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
			// Return the completed form as a message.
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

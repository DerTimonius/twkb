package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/DerTimonius/twkb/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskForm struct {
	help        help.Model
	description textinput.Model
	project     textinput.Model
	label       textinput.Model
	due         textinput.Model
	recur       textinput.Model
	until       textinput.Model
	col         column
	relatedTask Task
	index       int
	isEdit      bool
}

func newDefaultForm() *TaskForm {
	return NewForm("task name", "project (no spaces)", "labels (space separted list)", "due (e.g. eod, 2d)", "recur (e.g. monthly)", "until (e.g. now+1yr)")
}

func NewForm(description, project, label, due, recur, until string) *TaskForm {
	form := TaskForm{
		help:        help.New(),
		description: textinput.New(),
		project:     textinput.New(),
		label:       textinput.New(),
		due:         textinput.New(),
		recur:       textinput.New(),
		until:       textinput.New(),
	}
	form.description.Placeholder = description
	form.project.Placeholder = project
	form.label.Placeholder = label
	form.due.Placeholder = due
	form.recur.Placeholder = recur
	form.until.Placeholder = until
	form.description.Focus()
	return &form
}

func (f TaskForm) CreateTask() Task {
	cmdStr, err := AddCmd(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdStr[0], cmdStr[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get the id from the output of task
	output := out.String()
	re := regexp.MustCompile(`\d+`)
	matchId := re.FindString(output)
	id, e := strconv.Atoi(matchId)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	task := Task{id: id, status: todo, description: f.description.Value(), project: f.project.Value(), tags: strings.Split(f.label.Value(), " ")}
	task.UpdateUrgency()
	return task
}

func NewEditForm(t Task) *TaskForm {
	form := TaskForm{
		help:        help.New(),
		description: textinput.New(),
		project:     textinput.New(),
		label:       textinput.New(),
		due:         textinput.New(),
		isEdit:      true,
		relatedTask: t,
	}
	form.description.SetValue(t.description)
	form.project.SetValue(t.project)
	form.label.SetValue(strings.Join(t.tags, " "))
	form.due.SetValue(t.due)
	form.description.Focus()
	return &form
}

func (f TaskForm) Init() tea.Cmd {
	return nil
}

func (f TaskForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				f.recur.Focus()
				return f, textarea.Blink
			}
			if f.recur.Focused() {
				f.recur.Blur()
				f.until.Focus()
				return f, textarea.Blink
			}
			if f.until.Focused() {
				f.until.Blur()
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
	if f.recur.Focused() {
		f.recur, cmd = f.recur.Update(msg)
		return f, cmd
	}
	if f.until.Focused() {
		f.until, cmd = f.until.Update(msg)
		return f, cmd
	}
	f.due, cmd = f.due.Update(msg)
	return f, cmd
}

func (f TaskForm) View() string {
	title := styles.TitleStyle.Render("Create or update a Task")

	fieldStyle := styles.FieldStyle
	inputStyle := styles.InputStyle
	inputs := lipgloss.JoinVertical(
		lipgloss.Left,
		fieldStyle.Render(inputStyle.Render("Description: "+f.description.View())),
		fieldStyle.Render(inputStyle.Render("Project:     "+f.project.View())),
		fieldStyle.Render(inputStyle.Render("Label:       "+f.label.View())),
		fieldStyle.Render(inputStyle.Render("Due:         "+f.due.View())),
	)

	if !f.isEdit {
		inputs = lipgloss.JoinVertical(lipgloss.Left, inputs,
			fieldStyle.Render(inputStyle.Render("Recur:       "+f.recur.View())),
			fieldStyle.Render(inputStyle.Render("Until:       "+f.until.View())),
		)
	}

	help := f.help.View(keys)

	return styles.FormStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			inputs,
			strings.Repeat("─", 63), // Separator line
			help,
		),
	)
}

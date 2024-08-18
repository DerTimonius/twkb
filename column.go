package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/DerTimonius/twkb/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

type column struct {
	list   list.Model
	status status
	height int
	width  int
	focus  bool
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Blur() {
	c.focus = false
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(status status) column {
	var focus bool
	if status == todo {
		focus = true
	}
	defaultDelegate := list.NewDefaultDelegate()
	defaultDelegate.Styles.SelectedTitle = styles.DefaultSelectedTitleStyle
	defaultDelegate.Styles.SelectedDesc = styles.DefaultSelectedDesc
	defaultList := list.New([]list.Item{}, defaultDelegate, 0, 0)
	defaultList.SetShowHelp(false)
	defaultList.Styles.Title = styles.DefaultListTitleStyle
	return column{focus: focus, status: status, list: defaultList}
}

func (c column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height-8)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Edit):
			if len(c.list.VisibleItems()) != 0 {
				task := c.list.SelectedItem().(Task)
				f := NewEditForm(task)
				f.index = c.list.Index()
				f.col = c
				return f.Update(nil)
			}
		case key.Matches(msg, keys.New):
			f := newDefaultForm()
			f.index = APPEND
			f.col = c
			return f.Update(nil)
		case key.Matches(msg, keys.Unblock):
			task := c.list.SelectedItem().(Task)
			conf := NewConfirmation(fmt.Sprintf("Are you sure you want to unblock the task '%s'?", task.description), c.Unblock)
			conf.index = APPEND
			conf.column = c
			return conf.Update(nil)
		case key.Matches(msg, keys.Delete):
			task := c.list.SelectedItem().(Task)
			conf := NewConfirmation(fmt.Sprintf("Are you sure you want to delete the task '%s'?", task.description), c.DeleteCurrent)
			conf.index = APPEND
			conf.column = c
			return conf.Update(nil)
		case key.Matches(msg, keys.Block):
			task := c.list.SelectedItem().(Task)
			todoTasks := board.cols[todo].list.Items()
			b := NewBlockForm(task, todoTasks)
			b.index = APPEND
			b.column = c
			return b.Update(nil)
		case key.Matches(msg, keys.Space):
			return c, c.MoveToNext()
		case key.Matches(msg, keys.Enter):
			return c, c.MoveToDone()
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() tea.Cmd {
	var task Task
	var ok bool

	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}

	if len(c.list.VisibleItems()) > 0 {
		c.list.RemoveItem(c.list.Index())
	}

	cmdStr, err := DeleteCmd(&task)
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

	var cm tea.Cmd
	c.list, cm = c.list.Update(nil)
	return cm
}

func (c *column) Unblock() tea.Cmd {
	var task Task
	var ok bool

	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}

	task.UnblockTask()
	task.blocked = false

	var cm tea.Cmd
	c.list, cm = c.list.Update(nil)
	return cm
}

func (c *column) Set(i int, t Task) tea.Cmd {
	if i != APPEND {
		return c.list.SetItem(i, t)
	}
	return c.list.InsertItem(APPEND, t)
}

func (c *column) setSize(width, height int) {
	c.width = width / margin
	c.height = height - (margin * 2)
}

func (c *column) getStyle() lipgloss.Style {
	baseColumnStyle := styles.ColumnBaseStyle
	if c.Focused() {
		return baseColumnStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(styles.Blue)).
			Height(c.height).
			Width(c.width)
	}
	return baseColumnStyle.
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

type moveMsg struct {
	Task
}

func (c *column) MoveToNext() tea.Cmd {
	var task Task
	var ok bool
	// If nothing is selected, the SelectedItem will return Nil.
	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}

	// Don't move the task if it is in the done column
	if task.status == done {
		return nil
	}

	// move item
	c.list.RemoveItem(c.list.Index())
	task.StartStop()

	// refresh list
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{task} })
}

func (c *column) MoveToDone() tea.Cmd {
	var task Task
	var ok bool
	// If nothing is selected, the SelectedItem will return Nil.
	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}

	c.list.RemoveItem(c.list.Index())
	task.Finish()

	// refresh list
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{task} })
}

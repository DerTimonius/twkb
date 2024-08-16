package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type itemDelegate struct{}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("42"))
)

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Task)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.description)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Block struct {
	todoTasks     list.Model
	selectedTasks []Task
	column        column
	blocking      Task
	index         int
}

func (b Block) Init() tea.Cmd {
	return nil
}

func NewBlockForm(t Task, todos []list.Item) *Block {
	var filteredTodos []list.Item

	for _, td := range todos {
		if td.(Task).uuid == t.uuid || td.(Task).blocked || td.(Task).recurring {
			continue
		}

		filteredTodos = append(filteredTodos, td)
	}

	l := list.New([]list.Item{}, itemDelegate{}, 65, 45)
	l.Title = fmt.Sprintf("What other tasks does '%s' block?", t.description)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)

	l.SetItems(filteredTodos)

	return &Block{
		blocking:      t,
		todoTasks:     l,
		selectedTasks: []Task{},
		column:        column{},
		index:         0,
	}
}

func (b Block) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.column.setSize(msg.Width, msg.Height)
		b.column.list.SetSize(msg.Width/margin, msg.Height-8)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			return board.Update(b)
		case key.Matches(msg, keys.Space):
			selected := b.todoTasks.SelectedItem().(Task)
			b.selectedTasks = append(b.selectedTasks, selected)
			return b, nil
		case key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Quit):
			return b, tea.Quit
		}
	}
	list, cmd := b.todoTasks.Update(msg)
	b.todoTasks = list
	return b, cmd
}

func (b Block) View() string {
	return b.column.getStyle().Render(b.todoTasks.View())
}

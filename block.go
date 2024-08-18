package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/DerTimonius/twkb/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type blockItemDelegate struct {
	selectedTasks map[string]bool
}

func (d blockItemDelegate) Height() int                             { return 1 }
func (d blockItemDelegate) Spacing() int                            { return 0 }
func (d blockItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d blockItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Task)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.description)

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.BlockSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	if d.selectedTasks[i.uuid] {
		fn = func(s ...string) string {
			return styles.BlockSelectedItemStyle.
				Copy().
				Foreground(lipgloss.Color(styles.Green)).
				Render("* " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type Block struct {
	todoTaskList  list.Model
	selectedTasks map[string]bool
	help          help.Model
	todoTasks     []Task
	column        column
	blocking      Task
	index         int
}

func (b Block) Init() tea.Cmd {
	return nil
}

func NewBlockForm(t Task, todos []list.Item, height, width int) *Block {
	var filteredTodos []list.Item
	var filteredTasks []Task

	for _, td := range todos {
		if td.(Task).uuid == t.uuid || td.(Task).blocked || td.(Task).recurring {
			continue
		}

		filteredTodos = append(filteredTodos, td)
		filteredTasks = append(filteredTasks, td.(Task))
	}

	l := list.New([]list.Item{}, blockItemDelegate{}, width, height)

	l.Title = fmt.Sprintf("'%s' blocks?", t.description)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	l.SetItems(filteredTodos)
	b := Block{
		blocking:      t,
		todoTaskList:  l,
		todoTasks:     filteredTasks,
		selectedTasks: map[string]bool{},
		column:        column{},
		index:         0,
		help:          help.New(),
	}
	l.SetDelegate(blockItemDelegate{selectedTasks: b.selectedTasks})

	return &b
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
			selected := b.todoTaskList.SelectedItem().(Task)
			if b.selectedTasks[selected.uuid] {
				delete(b.selectedTasks, selected.uuid)
			} else {
				b.selectedTasks[selected.uuid] = true
			}
			return b, nil
		case key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Quit):
			return b, tea.Quit
		}
	}
	list, cmd := b.todoTaskList.Update(msg)
	b.todoTaskList = list
	return b, cmd
}

func (b Block) View() string {
	helpView := b.help.ShortHelpView(keys.BlockHelp())

	b.todoTaskList.SetDelegate(blockItemDelegate{selectedTasks: b.selectedTasks})
	content := b.todoTaskList.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		b.column.getStyle().Render(content),
		helpView,
	)
}

func (b Block) GetSelectedTasks() []Task {
	var tasks []Task
	for _, t := range b.todoTasks {
		if b.selectedTasks[t.uuid] {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

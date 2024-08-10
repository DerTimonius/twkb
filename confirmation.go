package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var confirmationStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("#96CDFB")).Padding(1).Width(75).AlignHorizontal(lipgloss.Center).Foreground(lipgloss.Color("#FAE3B0"))

type Confirmation struct {
	confirm func() tea.Cmd
	message string
	column  column
	index   int
}

func (c Confirmation) Init() tea.Cmd {
	return nil
}

func NewConfirmation(message string, confirm func() tea.Cmd) *Confirmation {
	return &Confirmation{
		confirm: confirm,
		message: message,
		column:  column{},
		index:   0,
	}
}

func (c Confirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Yes):
			return board.Update(c)
			// "n" is used for both 'no' and 'new'
		case key.Matches(msg, keys.New), key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Quit):
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c Confirmation) View() string {
	return confirmationStyle.Render(fmt.Sprintf("%s (y/n)", c.message))
}

package main

import (
	"fmt"

	"github.com/DerTimonius/twkb/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

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
		case key.Matches(msg, keys.No), key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Quit):
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c Confirmation) View() string {
	return styles.ConfirmationStyle.Render(fmt.Sprintf("%s (y/n)", c.message))
}

package main

import "github.com/charmbracelet/bubbles/key"

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Submit, k.Back, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Left, k.Right},
		{k.Space, k.Enter},
		{k.New, k.Edit},
		{k.Block, k.Unblock},
		{k.Filter, k.Quit},
	}
}

func (k keyMap) BlockHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.BlockSelect, k.BlockSubmit, k.Back}
}

type keyMap struct {
	New         key.Binding
	Edit        key.Binding
	Delete      key.Binding
	Up          key.Binding
	Down        key.Binding
	Right       key.Binding
	Left        key.Binding
	Enter       key.Binding
	Space       key.Binding
	Help        key.Binding
	Quit        key.Binding
	Back        key.Binding
	Tab         key.Binding
	Submit      key.Binding
	Filter      key.Binding
	Yes         key.Binding
	No          key.Binding
	Unblock     key.Binding
	Block       key.Binding
	BlockSelect key.Binding
	BlockSubmit key.Binding
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "add new task"),
	),
	Edit: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "modify focused task"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/l", "move left"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "start/stop task"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "finish task"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "b"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Goto next input"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit task"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter tasks"),
	),
	No: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "No"),
	),
	Yes: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "Yes"),
	),
	Unblock: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "unblock task"),
	),
	Block: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "block tasks"),
	),
	BlockSelect: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "select task"),
	),
	BlockSubmit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}

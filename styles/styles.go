package styles

import "github.com/charmbracelet/lipgloss"

const (
	Blue      = "#89b4fa"
	Flamingo  = "#f2cdcd"
	Gray      = "#45475a"
	Green     = "#a6e3a1"
	LightBlue = "#96CDFB"
	Maroon    = "#eba0ac"
	Mauve     = "#cba6f7"
	Peach     = "#ffb796"
	Red       = "#f38ba8"
	Sapphire  = "#74c7ec"
	Yellow    = "#f9e2af"
	Pink      = "#f896ad"
)

var (
	FormStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(Blue)).
			Padding(1).
			Width(65)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Blue)).
			Padding(0, 1).
			MarginBottom(1)

	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(Blue)).
			PaddingLeft(1)

	FieldStyle = lipgloss.NewStyle().
			MarginBottom(1)

	ColumnBaseStyle = lipgloss.NewStyle().Padding(1, 2)

	ConfirmationStyle = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color(Blue)).Padding(1).Width(75).AlignHorizontal(lipgloss.Center).Foreground(lipgloss.Color(Red))

	ItemStyle              = lipgloss.NewStyle().PaddingLeft(4)
	BlockSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color(LightBlue))

	DefaultSelectedTitleStyle = lipgloss.NewStyle().
					Border(lipgloss.NormalBorder(), false, false, false, true).
					BorderForeground(lipgloss.Color(Pink)).
					Foreground(lipgloss.Color(Pink)).
					Padding(0, 0, 0, 1)

	DefaultSelectedDesc = DefaultSelectedTitleStyle.Copy().
				Foreground(lipgloss.Color(Pink))

	DefaultListTitleStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(Blue)).
				Foreground(lipgloss.Color(Gray)).
				Padding(0, 1)
)

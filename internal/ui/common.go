package ui

import "github.com/charmbracelet/lipgloss"

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	PromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	HintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2)
)

// Messages
type NextScreenMsg struct {
	Screen string
}

type ErrMsg struct {
	Err error
}

func (e ErrMsg) Error() string { return e.Err.Error() }

// Session holds data during the installation flow
type Session struct {
	ArchiveURL  string
	Version     string
	ConfigCmd   string
	ProjectName string
	ConfigURL   string
	ConfigToken string
	InstallPath string
}

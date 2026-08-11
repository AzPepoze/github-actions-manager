package shared

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor = lipgloss.AdaptiveColor{Light: "#111111", Dark: "#F5F5F5"}
	mutedColor   = lipgloss.AdaptiveColor{Light: "#6B6B6B", Dark: "#9A9A9A"}
	dividerColor = lipgloss.AdaptiveColor{Light: "#B8B8B8", Dark: "#4A4A4A"}
	selectedText = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"}
	selectedFill = lipgloss.AdaptiveColor{Light: "#111111", Dark: "#F5F5F5"}
	successColor = lipgloss.AdaptiveColor{Light: "#18794E", Dark: "#5DDB9B"}
	errorColor   = lipgloss.AdaptiveColor{Light: "#B42318", Dark: "#FF7B72"}
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	PromptStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	HintStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	SubPromptStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			MarginBottom(1)

	CodeStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	LogStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2)

	TableStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(dividerColor)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	RowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	SelectedRowStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(selectedText).
				Background(selectedFill).
				Padding(0, 1)

	InputPromptStyle  = lipgloss.NewStyle().Foreground(primaryColor)
	DividerStyle      = lipgloss.NewStyle().Foreground(dividerColor)
	LabelStyle        = lipgloss.NewStyle().Foreground(mutedColor)
	FocusedLabelStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
)

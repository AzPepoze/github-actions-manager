package shared

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderPage renders a consistent page layout with a title, body, and help section.
func RenderPage(title string, body string, help string) string {
	var out strings.Builder

	out.WriteString(TitleStyle.Render(fmt.Sprintf("⚡ GitHub Actions Manager - %s", title)))
	out.WriteString("\n\n")

	out.WriteString(body)

	if help != "" {
		out.WriteString("\n")
		out.WriteString(help)
	}

	return ContainerStyle.Render(out.String())
}

// LabelStyle is for input labels
var LabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

// FocusedLabelStyle is for focused input labels
var FocusedLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Bold(true)

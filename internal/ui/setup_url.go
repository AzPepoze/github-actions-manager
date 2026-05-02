package ui

import (
	"actions-manager/internal/ops"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SetupURLModel struct {
	session    *Session
	urlInput   textinput.Model
	isParsed   bool
}

func NewSetupURLModel(session *Session) *SetupURLModel {
	input := textinput.New()
	input.Placeholder = "Paste curl command from GitHub"
	input.Focus()
	input.Width = 80

	return &SetupURLModel{
		session:  session,
		urlInput: input,
	}
}

func (m *SetupURLModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SetupURLModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && m.isParsed {
			return m, func() tea.Msg {
				return NextScreenMsg{Screen: "setup_config"}
			}
		}
	}

	m.urlInput, cmd = m.urlInput.Update(msg)

	// Live parse the input
	if m.urlInput.Value() != "" {
		if parsed, _ := ops.ParseCurl(m.urlInput.Value()); parsed != nil {
			m.session.ArchiveURL = parsed.URL
			m.session.Version = parsed.Version
			m.isParsed = true
		} else {
			m.isParsed = false
		}
	}

	return m, cmd
}

func (m *SetupURLModel) View() string {
	var body strings.Builder

	body.WriteString(TitleStyle.Render("⚡ Actions Manager [Setup: Runner URL]"))
	body.WriteString("\n\n")
	body.WriteString(PromptStyle.Render("Paste the curl download command from GitHub:"))
	body.WriteString("\n\n")
	body.WriteString(m.urlInput.View())
	body.WriteString("\n\n")

	if m.isParsed {
		body.WriteString(SuccessStyle.Render("✓ Valid command detected"))
		body.WriteString("\n")
		body.WriteString(HintStyle.Render("[Enter] Continue"))
	}

	return ContainerStyle.Render(body.String())
}

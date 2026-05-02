package ui

import (
	"actions-manager/internal/ops"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SetupConfigModel struct {
	session     *Session
	inputs      []textinput.Model
	focusIndex  int
}

const (
	inputConfig = iota
	inputName
)

func NewSetupConfigModel(session *Session) *SetupConfigModel {
	inputs := make([]textinput.Model, 2)

	config := textinput.New()
	config.Placeholder = "Config command (./config.sh ...)"
	config.Focus()
	config.Width = 80
	inputs[inputConfig] = config

	name := textinput.New()
	name.Placeholder = "Project name"
	name.Width = 30
	inputs[inputName] = name

	return &SetupConfigModel{
		session: session,
		inputs:  inputs,
	}
}

func (m *SetupConfigModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SetupConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			} else if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)

		case "enter":
			if m.inputs[inputConfig].Value() != "" && m.inputs[inputName].Value() != "" {
				m.session.ConfigCmd = m.inputs[inputConfig].Value()
				m.session.ProjectName = m.inputs[inputName].Value()
				return m, func() tea.Msg {
					return NextScreenMsg{Screen: "progress"}
				}
			}
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// Live parse config to auto-fill name
	if m.inputs[inputConfig].Value() != "" {
		if parsed, _ := ops.ParseConfig(m.inputs[inputConfig].Value()); parsed != nil {
			m.session.ConfigURL = parsed.URL
			m.session.ConfigToken = parsed.Token
			if m.inputs[inputName].Value() == "" {
				m.inputs[inputName].SetValue(parsed.ProjectName)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *SetupConfigModel) View() string {
	var body strings.Builder

	body.WriteString(TitleStyle.Render("⚡ Actions Manager [Setup: Configure Runner]"))
	body.WriteString("\n\n")
	body.WriteString(PromptStyle.Render("Paste the configuration command from GitHub:"))
	body.WriteString("\n\n")

	body.WriteString("Configuration command:\n")
	body.WriteString(m.inputs[inputConfig].View())
	body.WriteString("\n\n")

	body.WriteString("Project Name:\n")
	body.WriteString(m.inputs[inputName].View())
	body.WriteString("\n\n")

	body.WriteString(HintStyle.Render("[Tab] Switch fields  [Enter] Start installation"))

	return ContainerStyle.Render(body.String())
}

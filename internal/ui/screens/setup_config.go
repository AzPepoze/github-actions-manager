package screens

import (
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui/shared"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SetupConfigModel struct {
	session    *shared.Session
	inputs     []textinput.Model
	focusIndex int
	help       help.Model
	keys       setupKeyMap
}

type setupKeyMap struct {
	Next   key.Binding
	Prev   key.Binding
	Submit key.Binding
	Quit   key.Binding
}

var setupKeys = setupKeyMap{
	Next: key.NewBinding(
		key.WithKeys("tab", "down"),
		key.WithHelp("tab/↓", "next field"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab/↑", "prev field"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "install"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k setupKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev, k.Submit, k.Quit}
}

func (k setupKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Next, k.Prev, k.Submit, k.Quit}}
}

const (
	inputConfig = iota
	inputName
)

func NewSetupConfigModel(session *shared.Session) *SetupConfigModel {
	inputs := make([]textinput.Model, 2)

	config := textinput.New()
	config.Placeholder = "./config.sh --url https://github.com/user/repo --token ..."
	config.Focus()
	config.Width = 500
	config.Prompt = "  > "
	config.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	inputs[inputConfig] = config

	name := textinput.New()
	name.Placeholder = "Runner Name"
	name.Width = 30
	name.Prompt = "  > "
	name.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	inputs[inputName] = name

	return &SetupConfigModel{
		session: session,
		inputs:  inputs,
		help:    help.New(),
		keys:    setupKeys,
	}
}

func (m *SetupConfigModel) Init() tea.Cmd {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}
	m.focusIndex = 0
	m.inputs[0].Focus()
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
				m.session.RunnerName = m.inputs[inputName].Value()
				return m, func() tea.Msg {
					return shared.NavigateToMsg{Screen: shared.ScreenInstaller}
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
		if parsed, _ := service.ParseConfig(m.inputs[inputConfig].Value()); parsed != nil {
			m.session.ConfigURL = parsed.URL
			m.session.ConfigToken = parsed.Token
			m.session.ProjectName = parsed.ProjectName
			if m.inputs[inputName].Value() == "" {
				m.inputs[inputName].SetValue(parsed.ProjectName)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *SetupConfigModel) View() string {
	var body strings.Builder

	configLabel := shared.LabelStyle.Render("Configuration Command")
	if m.focusIndex == inputConfig {
		configLabel = shared.FocusedLabelStyle.Render("Configuration Command")
	}
	body.WriteString(configLabel + "\n")
	body.WriteString(m.inputs[inputConfig].View())
	body.WriteString("\n\n")

	nameLabel := shared.LabelStyle.Render("Runner Name")
	if m.focusIndex == inputName {
		nameLabel = shared.FocusedLabelStyle.Render("Runner Name")
	}
	body.WriteString(nameLabel + "\n")
	body.WriteString(m.inputs[inputName].View())
	body.WriteString("\n")

	return shared.RenderPage("Configure Runner", body.String(), m.help.View(m.keys))
}

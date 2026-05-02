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

type SetupURLModel struct {
	session  *shared.Session
	urlInput textinput.Model
	isParsed bool
	help     help.Model
	keys     urlKeyMap
}

type urlKeyMap struct {
	Submit key.Binding
	Quit   key.Binding
}

var urlKeys = urlKeyMap{
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "continue"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

func (k urlKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Submit, k.Quit}
}

func (k urlKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Submit, k.Quit}}
}

func NewSetupURLModel(session *shared.Session) *SetupURLModel {
	input := textinput.New()
	input.Placeholder = "curl -o actions-runner-linux-x64-2.334.0.tar.gz -L ..."
	input.Focus()
	input.Width = 200
	input.Prompt = "  > "
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	return &SetupURLModel{
		session:  session,
		urlInput: input,
		help:     help.New(),
		keys:     urlKeys,
	}
}

func (m *SetupURLModel) Init() tea.Cmd {
	m.urlInput.Reset()
	m.isParsed = false
	return textinput.Blink
}

func (m *SetupURLModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && m.isParsed {
			return m, func() tea.Msg {
				return shared.NavigateToMsg{Screen: shared.ScreenSetupConfig}
			}
		}
	}

	m.urlInput, cmd = m.urlInput.Update(msg)

	// Live parse the input
	if m.urlInput.Value() != "" {
		if parsed, _ := service.ParseCurl(m.urlInput.Value()); parsed != nil {
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

	body.WriteString(shared.FocusedLabelStyle.Render("Download Command") + "\n")
	body.WriteString(m.urlInput.View())
	body.WriteString("\n\n")

	if m.isParsed {
		body.WriteString(shared.SuccessStyle.Render("  ✓ Command recognized"))
		body.WriteString("\n")
	} else if m.urlInput.Value() != "" {
		body.WriteString(shared.ErrorStyle.Render("  ✗ Invalid command format"))
		body.WriteString("\n")
	} else {
		body.WriteString("\n")
	}

	return shared.RenderPage("Download Runner", body.String(), m.help.View(m.keys))
}

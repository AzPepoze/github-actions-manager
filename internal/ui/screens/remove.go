package screens

import (
	"fmt"
	"github-actions-manager/internal/core"
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui/shared"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RemoveModel struct {
	store     *core.Store
	runner    core.Runner
	nameInput textinput.Model
	help      help.Model
	keys      removeKeyMap
	status    string
	err       error
	finished  bool
	isForce   bool
}

type removeKeyMap struct {
	Submit key.Binding
	Mode   key.Binding
	Cancel key.Binding
}

var removeKeys = removeKeyMap{
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm removal"),
	),
	Mode: key.NewBinding(
		key.WithKeys("up", "down", "tab", "f"),
		key.WithHelp("↑/↓/tab/f", "switch mode"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "cancel"),
	),
}

func (k removeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Submit, k.Mode, k.Cancel}
}

func (k removeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Submit, k.Mode, k.Cancel}}
}

func NewRemoveModel(store *core.Store) *RemoveModel {
	ti := textinput.New()
	ti.Placeholder = "Type project name to confirm"
	ti.Focus()
	ti.Prompt = "  > "
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	return &RemoveModel{
		store:     store,
		nameInput: ti,
		help:      help.New(),
		keys:      removeKeys,
	}
}

func (m *RemoveModel) SetRunner(r core.Runner) {
	m.runner = r
	m.nameInput.Reset()
	m.nameInput.Focus()
	m.status = ""
	m.err = nil
	m.finished = false
	m.isForce = false
}

func (m *RemoveModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RemoveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.finished {
		if _, ok := msg.(tea.KeyMsg); ok {
			return m, func() tea.Msg { return shared.NavigateToMsg{Screen: shared.ScreenDashboard} }
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return shared.NavigateToMsg{Screen: shared.ScreenDashboard} }
		case "up", "down", "tab", "f":
			if !m.finished && m.status == "" {
				m.isForce = !m.isForce
			}
			return m, nil
		case "enter":
			if m.nameInput.Value() == m.runner.ProjectName {
				return m, m.runRemoveFlow()
			}
		}

	case shared.ErrMsg:
		m.err = msg.Err
		return m, func() tea.Msg {
			return shared.NavigateToMsg{Screen: shared.ScreenDashboard}
		}

	case shared.StatusMsg:
		m.status = string(msg)
		return m, nil

	case shared.ConfigDoneMsg:
		m.finished = true
		m.status = "Runner removed successfully."
		runners, _ := service.DiscoverRunners("actions")
		m.store.SetRunners(runners)
		return m, nil
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	return m, cmd
}

func (m *RemoveModel) runRemoveFlow() tea.Cmd {
	m.status = "Running removal and cleanup script..."

	// We use a single robust script to handle:
	// 1. Service stopping/uninstallation (only if svc.sh exists)
	// 2. Runner unconfiguration (config.sh remove) - Skip if isForce is enabled
	// 3. Final directory cleanup
	removeCmd := "./config.sh remove"
	if m.isForce {
		removeCmd = "echo '> Force removal: skipping runner unconfiguration (token remove)'"
	}

	script := fmt.Sprintf(`
if [ -f "./svc.sh" ]; then
    echo "> Stopping and uninstalling service..."
    sudo ./svc.sh stop || true
    sudo ./svc.sh uninstall || true
fi
echo "> Removing runner configuration..."
%s
echo "> Cleaning up installation directory..."
TARGET_DIR=$(basename "%s")
cd ..
sudo rm -rf "$TARGET_DIR"
`, removeCmd, m.runner.InstallPath)

	wrapped := shared.WrapWithPause(script, "Press ENTER to return to Dashboard...")

	c := exec.Command("/bin/sh", "-c", wrapped)
	c.Dir = m.runner.InstallPath

	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return shared.ErrMsg{Err: err}
		}
		return shared.ConfigDoneMsg{}
	})
}

func (m *RemoveModel) View() string {
	var body strings.Builder

	fmt.Fprintf(&body, "  Runner:  %s\n", shared.SuccessStyle.Render(m.runner.RunnerName))
	fmt.Fprintf(&body, "  Project: %s\n", shared.SuccessStyle.Render(m.runner.ProjectName))
	body.WriteString("\n\n")

	if m.finished {
		body.WriteString(shared.SuccessStyle.Render("  ✓ " + m.status))
		body.WriteString("\n\n")
		body.WriteString(shared.HintStyle.Render("  Press any key to return to Dashboard"))
	} else if m.status != "" {
		body.WriteString(shared.SuccessStyle.Render("  > " + m.status))
		body.WriteString("\n\n")
	} else {
		body.WriteString(shared.FocusedLabelStyle.Render("  Select Removal Mode:") + "\n")
		
		safeIndicator := "(*) "
		forceIndicator := "( ) "
		if m.isForce {
			safeIndicator = "( ) "
			forceIndicator = "(*) "
		}
		
		fmt.Fprintf(&body, "  %sSafe Remove (Token required)\n", safeIndicator)
		fmt.Fprintf(&body, "  %s", forceIndicator)
		if m.isForce {
			body.WriteString(shared.ErrorStyle.Bold(true).Render("Force Remove (Wipe files - NO TOKEN)"))
		} else {
			body.WriteString("Force Remove (Wipe files - NO TOKEN)")
		}
		body.WriteString("\n\n")

		body.WriteString(shared.FocusedLabelStyle.Render(fmt.Sprintf("  Type '%s' to confirm removal:", m.runner.ProjectName)) + "\n")
		body.WriteString(m.nameInput.View())
		body.WriteString("\n\n")

		if m.nameInput.Value() != "" && m.nameInput.Value() != m.runner.ProjectName {
			body.WriteString(shared.ErrorStyle.Render("    Name does not match"))
			body.WriteString("\n")
		}
	}

	return shared.RenderPage("Remove Runner", body.String(), m.help.View(m.keys))
}

package screens

import (
	"fmt"
	"github-actions-manager/internal/core"
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui/shared"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type DashboardModel struct {
	store  *core.Store
	help   help.Model
	keys   keyMap
	cursor int
	err    string
}

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Start   key.Binding
	Install key.Binding
	Remove  key.Binding
	Open    key.Binding
	Refresh key.Binding
	New     key.Binding
	Quit    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start/stop"),
	),
	Install: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "install/uninstall"),
	),
	Remove: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "remove"),
	),
	Open: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "cd/open repo"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Start, k.Install, k.Remove, k.Open, k.Refresh, k.New, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Start},
		{k.Install, k.Remove, k.Open, k.Refresh, k.New, k.Quit},
	}
}

func NewDashboardModel(store *core.Store) *DashboardModel {
	return &DashboardModel{
		store: store,
		help:  help.New(),
		keys:  keys,
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.store.Runners))
	for i := range m.store.Runners {
		cmds[i] = m.refreshStatusCmd(i)
	}
	return tea.Batch(cmds...)
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.err = ""
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			m.err = ""
			if m.cursor < len(m.store.Runners)-1 {
				m.cursor++
			}
		case "n":
			return m, func() tea.Msg {
				screen := shared.ScreenSetupURL
				if shared.HasRunnerArchive() {
					screen = shared.ScreenSetupConfig
				}
				return shared.NavigateToMsg{Screen: screen}
			}
		case "q":
			return m, tea.Quit
		case "s":
			return m, m.toggleServiceCmd()
		case "i":
			return m, m.toggleInstallCmd()
		case "x":
			return m, m.removeRunnerCmd()
		case "c":
			m.err = ""
			return m, m.openRepositoryCmd()
		case "r":
			m.err = ""
			runners, _ := service.DiscoverRunners("actions")
			m.store.SetRunners(runners)
			return m, m.Init()
		}

	case repositoryShellFinishedMsg:
		if msg.err != nil {
			m.err = fmt.Sprintf("Unable to open repository shell: %v", msg.err)
		}

	case shared.StatusRefreshMsg:
		_ = m.store.UpdateInstalled(msg.Index, msg.Installed)
		status := core.StatusStopped
		if msg.Running {
			status = core.StatusRunning
		}
		_ = m.store.UpdateStatus(msg.Index, status)
	}

	return m, nil
}

type repositoryShellFinishedMsg struct {
	err error
}

func (m *DashboardModel) refreshStatusCmd(index int) tea.Cmd {
	return func() tea.Msg {
		r := m.store.Runners[index]
		installed, running, _ := service.GetServiceStatus(r.InstallPath)
		return shared.StatusRefreshMsg{Index: index, Installed: installed, Running: running, FromTask: false}
	}
}

func (m *DashboardModel) toggleServiceCmd() tea.Cmd {
	if m.cursor < 0 || m.cursor >= len(m.store.Runners) {
		return nil
	}

	r := m.store.Runners[m.cursor]
	action := "start"
	if r.Status == core.StatusRunning {
		action = "stop"
	}

	cmd := service.ServiceCmd(r.InstallPath, action)
	index := m.cursor
	caser := cases.Title(language.English)
	return func() tea.Msg {
		return shared.SwitchToTaskMsg{
			Title:   fmt.Sprintf("%s Runner", caser.String(action)),
			Message: fmt.Sprintf("%s Runner", caser.String(action)),
			Cmd:     cmd,
			OnDone: func() tea.Msg {
				installed, running, _ := service.GetServiceStatus(r.InstallPath)
				return shared.StatusRefreshMsg{Index: index, Installed: installed, Running: running, FromTask: true}
			},
		}
	}
}

func (m *DashboardModel) toggleInstallCmd() tea.Cmd {
	if m.cursor < 0 || m.cursor >= len(m.store.Runners) {
		return nil
	}

	r := m.store.Runners[m.cursor]
	action := "install"
	if r.IsInstalled {
		action = "uninstall"
	}

	cmd := service.ServiceCmd(r.InstallPath, action)
	index := m.cursor
	caser := cases.Title(language.English)
	return func() tea.Msg {
		return shared.SwitchToTaskMsg{
			Title:   fmt.Sprintf("%s Service", caser.String(action)),
			Message: fmt.Sprintf("%s Service", caser.String(action)),
			Cmd:     cmd,
			OnDone: func() tea.Msg {
				installed, running, _ := service.GetServiceStatus(r.InstallPath)
				return shared.StatusRefreshMsg{Index: index, Installed: installed, Running: running, FromTask: true}
			},
		}
	}
}

func (m *DashboardModel) removeRunnerCmd() tea.Cmd {
	if m.cursor < 0 || m.cursor >= len(m.store.Runners) {
		return nil
	}

	r := m.store.Runners[m.cursor]
	return func() tea.Msg {
		return StartRemoveMsg{Runner: r}
	}
}

func (m *DashboardModel) openRepositoryCmd() tea.Cmd {
	if m.cursor < 0 || m.cursor >= len(m.store.Runners) {
		return nil
	}

	cmd, err := service.RepositoryShellCmd(m.store.Runners[m.cursor])
	if err != nil {
		m.err = err.Error()
		return nil
	}

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return repositoryShellFinishedMsg{err: err}
	})
}

type StartRemoveMsg struct {
	Runner core.Runner
}

func (m *DashboardModel) renderCell(content string, width int, style lipgloss.Style) string {
	truncated := truncate.StringWithTail(content, uint(width), "...")
	return style.Width(width).Render(truncated)
}

func (m *DashboardModel) View() string {
	var body strings.Builder

	if len(m.store.Runners) == 0 {
		body.WriteString("No runners configured yet.\n")
	} else {

		wStatus := 15
		wStartup := 12
		wName := 30
		wPath := 45

		header := lipgloss.JoinHorizontal(lipgloss.Top,
			m.renderCell("Status", wStatus, shared.HeaderStyle),
			m.renderCell("Startup", wStartup, shared.HeaderStyle),
			m.renderCell("Project Name", wPath, shared.HeaderStyle),
			m.renderCell("Runner Name", wName, shared.HeaderStyle),
		)
		body.WriteString(header + "\n")
		body.WriteString(shared.DividerStyle.Render(strings.Repeat("─", wStatus+wStartup+wName+wPath)) + "\n")

		for i, r := range m.store.Runners {
			status := shared.ErrorStyle.Render("STOPPED")
			if r.Status == core.StatusRunning {
				status = shared.SuccessStyle.Render("RUNNING")
			}

			startup := shared.ErrorStyle.Render("NO")
			if r.IsInstalled {
				startup = shared.SuccessStyle.Render("YES")
			}

			style := shared.RowStyle
			if m.cursor == i {
				style = shared.SelectedRowStyle
			}

			row := lipgloss.JoinHorizontal(lipgloss.Top,
				m.renderCell(status, wStatus, style),
				m.renderCell(startup, wStartup, style),
				m.renderCell(r.ProjectName, wPath, style),
				m.renderCell(r.RunnerName, wName, style),
			)

			body.WriteString(row + "\n")
		}
	}

	if m.err != "" {
		body.WriteString("\n" + shared.ErrorStyle.Render(m.err) + "\n")
	}

	return shared.RenderPage("Dashboard", body.String(), m.help.View(m.keys))
}

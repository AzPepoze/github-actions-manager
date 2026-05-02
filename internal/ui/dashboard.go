package ui

import (
	"actions-manager/internal/runner"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type DashboardModel struct {
	store  *runner.Store
	cursor int
}

func NewDashboardModel(store *runner.Store) *DashboardModel {
	return &DashboardModel{
		store: store,
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	return nil
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, func() tea.Msg {
				return NextScreenMsg{Screen: "setup_url"}
			}
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.store.Runners)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m *DashboardModel) View() string {
	var body strings.Builder

	body.WriteString(TitleStyle.Render("⚡ Actions Manager [Dashboard]"))
	body.WriteString("\n\n")

	if len(m.store.Runners) == 0 {
		body.WriteString("No runners configured yet.\n")
	} else {
		for i, r := range m.store.Runners {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			status := "[STOPPED]"
			if r.Status == runner.StatusRunning {
				status = SuccessStyle.Render("[RUNNING]")
			}

			body.WriteString(fmt.Sprintf("%s %s %s (%s)\n", cursor, status, r.ProjectName, r.ProjectURL))
		}
	}

	body.WriteString("\n")
	body.WriteString(HintStyle.Render("[n] New Runner  [q] Quit"))

	return ContainerStyle.Render(body.String())
}

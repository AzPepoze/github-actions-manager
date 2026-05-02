package ui

import (
	"actions-manager/internal/runner"
	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenSetupURL
	ScreenSetupConfig
	ScreenProgress
)

type RootModel struct {
	store   *runner.Store
	screen  Screen
	session *Session

	dashboard   tea.Model
	setupURL    tea.Model
	setupConfig tea.Model
	progress    tea.Model
}

func NewRootModel(store *runner.Store) RootModel {
	m := RootModel{
		store:   store,
		session: &Session{},
	}
	
	m.dashboard = NewDashboardModel(store)
	m.setupURL = NewSetupURLModel(m.session)
	m.setupConfig = NewSetupConfigModel(m.session)
	m.progress = NewProgressModel(store, m.session)

	if len(store.Runners) > 0 {
		m.screen = ScreenDashboard
	} else {
		m.screen = ScreenSetupURL
	}

	return m
}

func (m RootModel) Init() tea.Cmd {
	switch m.screen {
	case ScreenDashboard:
		return m.dashboard.Init()
	default:
		return m.setupURL.Init()
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case NextScreenMsg:
		return m.handleNavigation(msg.Screen)
	}

	var cmd tea.Cmd
	switch m.screen {
	case ScreenDashboard:
		m.dashboard, cmd = m.dashboard.Update(msg)
	case ScreenSetupURL:
		m.setupURL, cmd = m.setupURL.Update(msg)
	case ScreenSetupConfig:
		m.setupConfig, cmd = m.setupConfig.Update(msg)
	case ScreenProgress:
		m.progress, cmd = m.progress.Update(msg)
	}

	return m, cmd
}

func (m *RootModel) handleNavigation(screenName string) (tea.Model, tea.Cmd) {
	switch screenName {
	case "setup_url":
		m.screen = ScreenSetupURL
		return m, m.setupURL.Init()
	case "setup_config":
		m.screen = ScreenSetupConfig
		return m, m.setupConfig.Init()
	case "progress":
		m.screen = ScreenProgress
		return m, m.progress.Init()
	case "dashboard":
		m.screen = ScreenDashboard
		// Re-initialize dashboard to refresh list
		m.dashboard = NewDashboardModel(m.store)
		return m, m.dashboard.Init()
	}
	return m, nil
}

func (m RootModel) View() string {
	switch m.screen {
	case ScreenDashboard:
		return m.dashboard.View()
	case ScreenSetupURL:
		return m.setupURL.View()
	case ScreenSetupConfig:
		return m.setupConfig.View()
	case ScreenProgress:
		return m.progress.View()
	default:
		return "Error: Unknown screen"
	}
}

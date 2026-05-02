package ui

import (
	"github-actions-manager/internal/core"
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui/components"
	"github-actions-manager/internal/ui/screens"
	"github-actions-manager/internal/ui/shared"

	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	store   *core.Store
	screen  shared.Screen
	session *shared.Session
	screens map[shared.Screen]tea.Model
	program *tea.Program
	Err     error
}

func NewRootModel(store *core.Store) *RootModel {
	session := &shared.Session{}
	m := &RootModel{
		store:   store,
		session: session,
		screens: make(map[shared.Screen]tea.Model),
	}

	m.screens[shared.ScreenDashboard] = screens.NewDashboardModel(store)
	m.screens[shared.ScreenSetupURL] = screens.NewSetupURLModel(session)
	m.screens[shared.ScreenSetupConfig] = screens.NewSetupConfigModel(session)
	m.screens[shared.ScreenInstaller] = screens.NewInstallerModel(store, session)
	m.screens[shared.ScreenTask] = components.NewTaskModel()
	m.screens[shared.ScreenRemove] = screens.NewRemoveModel(store)

	if len(store.Runners) > 0 {
		m.screen = shared.ScreenDashboard
	} else if shared.HasRunnerArchive() {
		m.screen = shared.ScreenSetupConfig
	} else {
		m.screen = shared.ScreenSetupURL
	}

	return m
}

func (m *RootModel) SetProgram(p *tea.Program) {
	m.program = p
	if task, ok := m.screens[shared.ScreenTask].(*components.TaskModel); ok {
		task.SetProgram(p)
	}
}

func (m *RootModel) Init() tea.Cmd {
	return m.screens[m.screen].Init()
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case shared.NavigateToMsg:
		m.screen = msg.Screen
		return m, m.screens[m.screen].Init()

	case screens.StartRemoveMsg:
		m.screen = shared.ScreenRemove
		if remove, ok := m.screens[shared.ScreenRemove].(*screens.RemoveModel); ok {
			remove.SetRunner(msg.Runner)
			return m, remove.Init()
		}

	case shared.SwitchToTaskMsg:
		m.screen = shared.ScreenTask
		if task, ok := m.screens[shared.ScreenTask].(*components.TaskModel); ok {
			task.SetTask(msg.Title, msg.Message, msg.Cmd, msg.OnDone)
			return m, task.Init()
		}

	case shared.StatusRefreshMsg:
		if msg.FromTask {
			runners, _ := service.DiscoverRunners("actions")
			m.store.SetRunners(runners)
			m.screen = shared.ScreenDashboard
			return m, m.screens[m.screen].Init()
		}

	case shared.ErrMsg:
		m.Err = msg.Err

	}

	if activeModel, ok := m.screens[m.screen]; ok {
		var cmd tea.Cmd
		m.screens[m.screen], cmd = activeModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *RootModel) View() string {
	if activeModel, ok := m.screens[m.screen]; ok {
		return activeModel.View()
	}
	return "Error: Unknown screen"
}

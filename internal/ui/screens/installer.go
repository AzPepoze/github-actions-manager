package screens

import (
	"fmt"
	"github-actions-manager/internal/core"
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui/screens/install_steps"
	"github-actions-manager/internal/ui/shared"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type InstallerModel struct {
	store        *core.Store
	session      *shared.Session
	progress     progress.Model
	status       string
	percent      float64
	progressChan chan service.DownloadProgress
	err          error
}

func NewInstallerModel(store *core.Store, session *shared.Session) *InstallerModel {
	return &InstallerModel{
		store:    store,
		session:  session,
		progress: progress.New(progress.WithDefaultGradient()),
		status:   "Initializing...",
	}
}

func (m *InstallerModel) Init() tea.Cmd {
	m.status = "Starting download..."
	m.percent = 0
	
	m.progressChan = make(chan service.DownloadProgress, 20)
	
	return tea.Batch(
		install_steps.DownloadStep(m.session.ArchiveURL, m.progressChan),
		install_steps.WaitForProgress(m.progressChan),
	)
}

func (m *InstallerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case shared.ProgressMsg:
		if msg.Total > 0 {
			m.percent = float64(msg.Downloaded) / float64(msg.Total)
		}
		return m, install_steps.WaitForProgress(m.progressChan)

	case shared.StatusMsg:
		m.status = string(msg)
		if strings.Contains(m.status, "Extracting") {
			m.percent = 0.5
			return m, install_steps.ExtractStep(m.session.RunnerName, "actions-runner.tar.gz")
		}
		return m, nil

	case shared.DoneMsg:
		m.session.InstallPath = msg.Path
		m.status = "Starting configuration..."
		m.percent = 0.8
		return m, install_steps.RunCommandStep(m.session.ConfigCmd, m.session.InstallPath, "config")

	case shared.ConfigDoneMsg:
		m.status = "Configuration finished. Installing service..."
		m.percent = 0.9
		return m, install_steps.RunCommandStep("sudo ./svc.sh install", m.session.InstallPath, "install")

	case shared.InstallDoneMsg:
		m.status = "Service installed. Starting..."
		return m, install_steps.RunCommandStep("sudo ./svc.sh start", m.session.InstallPath, "start")

	case shared.ErrMsg:
		m.err = msg.Err
		return m, func() tea.Msg {
			return shared.NavigateToMsg{Screen: shared.ScreenDashboard}
		}

	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		m.progress = newModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m *InstallerModel) View() string {
	var body strings.Builder

	fmt.Fprintf(&body, "Status: %s\n\n", m.status)
	body.WriteString(m.progress.ViewAs(m.percent) + "\n\n")
	body.WriteString(shared.HintStyle.Render("Running installation steps..."))

	return shared.RenderPage("Installing Runner", body.String(), "")
}

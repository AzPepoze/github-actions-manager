package ui

import (
	"actions-manager/internal/ops"
	"actions-manager/internal/runner"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type ProgressModel struct {
	store    *runner.Store
	session  *Session
	progress progress.Model
	status   string
	percent  float64
	err      error
}

func NewProgressModel(store *runner.Store, session *Session) *ProgressModel {
	return &ProgressModel{
		store:    store,
		session:  session,
		progress: progress.New(progress.WithDefaultGradient()),
		status:   "Initializing...",
	}
}

func (m *ProgressModel) Init() tea.Cmd {
	m.status = "Starting download..."
	m.percent = 0
	return m.downloadCmd()
}

type statusMsg string
type doneMsg struct{ path string }

func (m *ProgressModel) downloadCmd() tea.Cmd {
	return func() tea.Msg {
		archivePath := "actions-runner.tar.gz"
		if err := ops.DownloadRunner(m.session.ArchiveURL, archivePath, nil); err != nil {
			return ErrMsg{Err: err}
		}
		return statusMsg("Download complete. Extracting...")
	}
}

func (m *ProgressModel) extractCmd() tea.Cmd {
	return func() tea.Msg {
		path, err := ops.ExtractRunner(m.session.ProjectName, "actions-runner.tar.gz")
		if err != nil {
			return ErrMsg{Err: err}
		}
		return doneMsg{path: path}
	}
}

func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = string(msg)
		if strings.Contains(m.status, "Extracting") {
			m.percent = 0.5
			return m, m.extractCmd()
		}
		return m, nil

	case doneMsg:
		m.session.InstallPath = msg.path
		m.status = "Preparing interactive configuration..."
		m.percent = 0.8
		return m, m.runHandoffCmd()

	case ErrMsg:
		m.err = msg.Err
		m.status = fmt.Sprintf("Error: %v", m.err)
		return m, nil

	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		m.progress = newModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m *ProgressModel) runHandoffCmd() tea.Cmd {
	// 1. Run config.sh
	config := exec.Command("/bin/sh", "-c", m.session.ConfigCmd)
	config.Dir = m.session.InstallPath

	// 2. Install service
	installSvc := exec.Command("sudo", "./svc.sh", "install")
	installSvc.Dir = m.session.InstallPath

	// 3. Start service
	startSvc := exec.Command("sudo", "./svc.sh", "start")
	startSvc.Dir = m.session.InstallPath

	return tea.Sequence(
		tea.ExecProcess(config, func(err error) tea.Msg {
			if err != nil {
				return ErrMsg{Err: err}
			}
			return statusMsg("Config complete. Installing service...")
		}),
		tea.ExecProcess(installSvc, func(err error) tea.Msg {
			if err != nil {
				return ErrMsg{Err: err}
			}
			return statusMsg("Service installed. Starting...")
		}),
		tea.ExecProcess(startSvc, func(err error) tea.Msg {
			if err != nil {
				return ErrMsg{Err: err}
			}
			// Save to store
			_ = m.store.Add(runner.Runner{
				ProjectName: m.session.ProjectName,
				ProjectURL:  m.session.ConfigURL,
				Token:       m.session.ConfigToken,
				InstallPath: m.session.InstallPath,
				Status:      runner.StatusRunning,
			})
			return NextScreenMsg{Screen: "dashboard"}
		}),
	)
}

func (m *ProgressModel) View() string {
	var body strings.Builder

	body.WriteString(TitleStyle.Render("⚡ Actions Manager [Installing Runner]"))
	body.WriteString("\n\n")
	body.WriteString(fmt.Sprintf("Status: %s\n\n", m.status))
	
	if m.err != nil {
		body.WriteString(ErrorStyle.Render(fmt.Sprintf("✘ %v", m.err)))
		body.WriteString("\n\n")
		body.WriteString(HintStyle.Render("Press Ctrl+C to exit"))
	} else {
		body.WriteString(m.progress.ViewAs(m.percent))
		body.WriteString("\n\n")
		body.WriteString(HintStyle.Render("Running background tasks..."))
	}

	return ContainerStyle.Render(body.String())
}

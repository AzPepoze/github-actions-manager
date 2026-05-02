package shared

import (
	"github-actions-manager/internal/service"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
)

type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenSetupURL
	ScreenSetupConfig
	ScreenInstaller
	ScreenTask
	ScreenRemove
)

type StatusRefreshMsg struct {
	Index     int
	Installed bool
	Running   bool
	FromTask  bool
}

type NavigateToMsg struct {
	Screen Screen
}

type SwitchToTaskMsg struct {
	Title   string
	Message string
	Cmd     *exec.Cmd
	OnDone  func() tea.Msg
}

type ErrMsg struct {
	Err error
}

func (e ErrMsg) Error() string { return e.Err.Error() }

type StatusMsg string
type DoneMsg struct{ Path string }
type ProgressMsg service.DownloadProgress
type ConfigDoneMsg struct{}
type InstallDoneMsg struct{}

// Session holds data during the installation flow
type Session struct {
	ArchiveURL  string
	Version     string
	ConfigCmd   string
	RunnerName  string
	ProjectName string
	ConfigURL   string
	ConfigToken string
	InstallPath string
}

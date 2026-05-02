package core

type Status string

const (
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
	StatusUnknown Status = "unknown"
)

type Runner struct {
	RunnerName  string `json:"runner_name"`
	ProjectName string `json:"project_name"`
	ProjectURL  string `json:"project_url"`
	Token       string `json:"token"`
	InstallPath string `json:"install_path"`
	Status      Status `json:"status"`
	IsInstalled bool   `json:"is_installed"`
}

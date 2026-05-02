package install_steps

import (
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui/shared"

	tea "github.com/charmbracelet/bubbletea"
)

func ExtractStep(runnerName string, archivePath string) tea.Cmd {
	return func() tea.Msg {
		path, err := service.ExtractRunner(runnerName, archivePath)
		if err != nil {
			return shared.ErrMsg{Err: err}
		}
		return shared.DoneMsg{Path: path}
	}
}

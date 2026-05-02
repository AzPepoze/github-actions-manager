package install_steps

import (
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui/shared"

	tea "github.com/charmbracelet/bubbletea"
)

func DownloadStep(archiveURL string, progressChan chan service.DownloadProgress) tea.Cmd {
	return func() tea.Msg {
		archivePath := "actions-runner.tar.gz"
		if err := service.DownloadRunner(archiveURL, archivePath, progressChan); err != nil {
			return shared.ErrMsg{Err: err}
		}
		return nil
	}
}

func WaitForProgress(progressChan chan service.DownloadProgress) tea.Cmd {
	return func() tea.Msg {
		progress, ok := <-progressChan
		if !ok {
			return shared.StatusMsg("Download complete. Extracting...")
		}
		return shared.ProgressMsg(progress)
	}
}

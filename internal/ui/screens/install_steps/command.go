package install_steps

import (
	"fmt"
	"github-actions-manager/internal/ui/shared"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func RunCommandStep(commandString string, directory string, stepName string) tea.Cmd {
	successMsg := fmt.Sprintf("%s completed. Press ENTER to continue...", stepName)
	errorMsg := fmt.Sprintf("Error occurred during %s. Press ENTER to continue...", stepName)
	title := fmt.Sprintf("RUNNING %s", strings.ToUpper(stepName))
	wrappedCommand := shared.WrapCommand(title, commandString, successMsg, errorMsg, true, true)

	command := exec.Command("/bin/sh", "-c", wrappedCommand)
	command.Dir = directory

	return tea.ExecProcess(command, func(err error) tea.Msg {
		if err != nil {
			return shared.ErrMsg{Err: err}
		}
		switch stepName {
		case "config":
			return shared.ConfigDoneMsg{}
		case "install":
			return shared.InstallDoneMsg{}
		case "start":
			return shared.NavigateToMsg{Screen: shared.ScreenDashboard}
		}
		return nil
	})
}

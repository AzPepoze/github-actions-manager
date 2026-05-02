package shared

import "fmt"

const separator = "-------------------------------------------"
const titleSeparator = "==========================================="

func WrapCommand(title, cmd, successMsg, errorMsg string, pauseOnComplete, pauseOnError bool) string {
	titleBlock := fmt.Sprintf("echo '%s'\necho '%s'\necho '%s'\necho", titleSeparator, title, titleSeparator)

	pauseOnCompleteCmd := ":"
	if pauseOnComplete {
		pauseOnCompleteCmd = "read"
	}

	pauseOnErrorCmd := "exit 1"
	if pauseOnError {
		pauseOnErrorCmd = "read; exit 1"
	}

	return fmt.Sprintf(`clear
%s
%s
EXIT_CODE=$?
echo
echo '%s'
if [ $EXIT_CODE -eq 0 ]; then
    echo '%s'
    %s
else
    echo '%s'
    %s
fi
exit $EXIT_CODE`, titleBlock, cmd, separator, successMsg, pauseOnCompleteCmd, errorMsg, pauseOnErrorCmd)
}

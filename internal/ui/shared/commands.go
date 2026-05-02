package shared

import "fmt"

const separator = "-------------------------------------------"

// WrapWithPause wraps a shell command to always pause and wait for ENTER after execution.
func WrapWithPause(cmd string, pauseMsg string) string {
	return fmt.Sprintf("%s\necho; echo '%s'; echo '%s'; read", cmd, separator, pauseMsg)
}

// WrapWithPauseOnError wraps a shell command to pause only if the command fails.
func WrapWithPauseOnError(cmd string, pauseMsg string) string {
	return fmt.Sprintf("%s || { echo; echo '%s'; echo '%s'; read; exit 1; }", cmd, separator, pauseMsg)
}

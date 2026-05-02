package shared

import "os"

func HasRunnerArchive() bool {
	_, err := os.Stat("actions-runner.tar.gz")
	return err == nil
}

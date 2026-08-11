package service

import (
	"fmt"
	"github-actions-manager/internal/core"
	"os"
	"os/exec"
	"path/filepath"
)

const defaultShell = "/bin/sh"

// RepositoryWorkDir returns the working checkout created by a GitHub Actions
// runner for the supplied repository.
func RepositoryWorkDir(runner core.Runner) string {
	return filepath.Join(runner.InstallPath, "_work", runner.ProjectName, runner.ProjectName)
}

// RepositoryShellCmd validates the runner checkout and returns an interactive
// shell command rooted in it.
func RepositoryShellCmd(runner core.Runner) (*exec.Cmd, error) {
	workDir := RepositoryWorkDir(runner)
	info, err := os.Stat(workDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("repository checkout not found: %s", workDir)
		}
		return nil, fmt.Errorf("inspect repository checkout: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("repository checkout is not a directory: %s", workDir)
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = defaultShell
	}

	cmd := exec.Command(shell)
	cmd.Dir = workDir
	return cmd, nil
}

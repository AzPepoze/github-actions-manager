package service

import (
	"github-actions-manager/internal/core"
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryWorkDir(t *testing.T) {
	tests := []struct {
		name        string
		installPath string
		projectName string
		want        string
	}{
		{
			name:        "relative install path",
			installPath: "actions/runner-1",
			projectName: "repo",
			want:        filepath.Join("actions", "runner-1", "_work", "repo", "repo"),
		},
		{
			name:        "absolute install path",
			installPath: "/srv/actions/runner-1",
			projectName: "repo",
			want:        "/srv/actions/runner-1/_work/repo/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RepositoryWorkDir(core.Runner{InstallPath: tt.installPath, ProjectName: tt.projectName})
			if got != tt.want {
				t.Fatalf("RepositoryWorkDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRepositoryShellCmd(t *testing.T) {
	t.Setenv("SHELL", "/bin/sh")
	installPath := t.TempDir()
	runner := core.Runner{InstallPath: installPath, ProjectName: "repo"}

	t.Run("existing checkout", func(t *testing.T) {
		workDir := RepositoryWorkDir(runner)
		if err := os.MkdirAll(workDir, 0o755); err != nil {
			t.Fatal(err)
		}

		cmd, err := RepositoryShellCmd(runner)
		if err != nil {
			t.Fatal(err)
		}
		if cmd.Path != "/bin/sh" {
			t.Fatalf("command path = %q, want /bin/sh", cmd.Path)
		}
		if cmd.Dir != workDir {
			t.Fatalf("command directory = %q, want %q", cmd.Dir, workDir)
		}
	})

	t.Run("missing checkout", func(t *testing.T) {
		missing := core.Runner{InstallPath: t.TempDir(), ProjectName: "missing"}
		if _, err := RepositoryShellCmd(missing); err == nil {
			t.Fatal("RepositoryShellCmd() error = nil, want error")
		}
	})

	t.Run("checkout is a file", func(t *testing.T) {
		fileRunner := core.Runner{InstallPath: t.TempDir(), ProjectName: "file"}
		workDir := RepositoryWorkDir(fileRunner)
		if err := os.MkdirAll(filepath.Dir(workDir), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(workDir, []byte("not a directory"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := RepositoryShellCmd(fileRunner); err == nil {
			t.Fatal("RepositoryShellCmd() error = nil, want error")
		}
	})
}

func TestRepositoryShellCmdUsesDefaultShell(t *testing.T) {
	t.Setenv("SHELL", "")
	installPath := t.TempDir()
	runner := core.Runner{InstallPath: installPath, ProjectName: "repo"}
	if err := os.MkdirAll(RepositoryWorkDir(runner), 0o755); err != nil {
		t.Fatal(err)
	}

	cmd, err := RepositoryShellCmd(runner)
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Path != defaultShell {
		t.Fatalf("command path = %q, want %q", cmd.Path, defaultShell)
	}
}

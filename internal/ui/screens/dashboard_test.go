package screens

import (
	"github-actions-manager/internal/core"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDashboardOpenRepositoryWithoutSelection(t *testing.T) {
	store := core.NewStore()
	model := NewDashboardModel(store)

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	if updated != model {
		t.Fatal("Update() returned a different dashboard model")
	}
	if cmd != nil {
		t.Fatal("Update() returned a command without a selected runner")
	}
}

func TestDashboardOpenRepositoryShowsResolutionError(t *testing.T) {
	store := core.NewStore()
	store.SetRunners([]core.Runner{{InstallPath: t.TempDir(), ProjectName: "missing"}})
	model := NewDashboardModel(store)

	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	if model.err == "" {
		t.Fatal("expected repository resolution error")
	}
	if got := model.View(); got == "" {
		t.Fatal("View() returned empty output")
	}
}

func TestDashboardShellCompletionDisplaysError(t *testing.T) {
	store := core.NewStore()
	model := NewDashboardModel(store)

	_, _ = model.Update(repositoryShellFinishedMsg{err: errTest{}})
	if model.err == "" {
		t.Fatal("expected shell completion error")
	}
}

type errTest struct{}

func (errTest) Error() string { return "shell failed" }

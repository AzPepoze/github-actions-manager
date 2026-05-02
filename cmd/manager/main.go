package main

import (
	"fmt"
	"github-actions-manager/internal/core"
	"github-actions-manager/internal/service"
	"github-actions-manager/internal/ui"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	_ = exec.Command("sudo", "-v").Run()

	// Initialize runner store
	store := core.NewStore()
	runners, err := service.DiscoverRunners("actions")
	if err != nil {
		fmt.Printf("Error discovering runners: %v\n", err)
		os.Exit(1)
	}
	store.SetRunners(runners)

	// Create root model
	model := ui.NewRootModel(store)

	// Run the program
	program := tea.NewProgram(model, tea.WithAltScreen())
	model.SetProgram(program)
	
	m, err := program.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	finalModel := m.(*ui.RootModel)
	if finalModel.Err != nil {
		os.Exit(1)
	}
}

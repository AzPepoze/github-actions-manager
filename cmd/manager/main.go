package main

import (
	"actions-manager/internal/runner"
	"actions-manager/internal/ui"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize runner store
	store, err := runner.NewStore("runners.json")
	if err != nil {
		fmt.Printf("Error initializing store: %v\n", err)
		os.Exit(1)
	}

	// Create root model
	model := ui.NewRootModel(store)

	// Run the program
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

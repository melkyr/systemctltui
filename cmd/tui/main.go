// package main
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"systemctltui/internal/tui" // Import your main TUI package
)

func main() {
	// Create a new instance of your TUI model
	initialModel := tui.NewModel()

	// Create and start the bubbletea program
	p := tea.NewProgram(
		initialModel,
		tea.WithAltScreen(), // Use the alternate screen buffer
		// tea.WithMouseCellMotion(), // Optional: enable mouse support
	)

	// Run the program
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

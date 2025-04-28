// package tui

package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"systemctltui/internal/listui"
	"systemctltui/internal/constants"
	// No need to import lipgloss or styles here
)

// AppState defines the different states of the application view.
type AppState int

const (
	StateBrowse AppState = iota // Browse tabs and lists
	StatePreview                  // Showing command preview
	StateOutput                   // Showing command output/error
)

// model represents the main state of the TUI application.
type model struct {
	activeTab int
	tabs      []string
	lists     []list.Model
	showHelp  bool
	width, height int

	// State for command preview and execution
	state           AppState // Current application state
	selectedCommand string   // The command string (e.g., "status")
	selectedUnit    string   // The unit string (e.g., "nginx.service")
	previewCommand  string   // The full command to be previewed/executed
	commandOutput   string   // Output/error from command execution
}

// NewModel initializes the main application model.
func NewModel() model {
	tabs := []string{"Global Options", "Commands", "Units"}
	lists := listui.NewLists()

	return model{
		activeTab: constants.TabOptions, // Use constant for initial tab
		tabs:      tabs,
		lists:     lists,
		showHelp:  false,
        // width/height are zero initially, set by WindowSizeMsg
		state:     StateBrowse, // Start in Browse state
	}
}

// Init performs initial setup for the bubbletea model.
func (m model) Init() tea.Cmd {
	// nil is sufficient here, as WindowSizeMsg is sent automatically on startup.
    return nil
}

// Update and View methods will be defined in update.go and view.go

// package tui
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"systemctltui/internal/constants"
	"systemctltui/internal/listui"
	"systemctltui/internal/system"
	"systemctltui/internal/messages" // <--- Import the CORRECT messages package
	"fmt"
	"strings"
)

// Update handles messages and updates the model state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle messages based on the current state
	switch m.state {
	case StateBrowse:
		return updateBrowse(m, msg)
	case StatePreview:
		return updatePreview(m, msg)
	case StateOutput:
		return updateOutput(m, msg)
	default:
		// Should not happen
		return m, nil
	}
}

// updateBrowse handles messages when the user is Browse tabs and lists.
func updateBrowse(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	// If help is shown, any keypress hides it
	if m.showHelp {
		if _, ok := msg.(tea.KeyMsg); ok {
			m.showHelp = false
		}
		return m, nil
	}

	// Handle key presses specific to Browse
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
            // Optional: Reset filter or scroll when switching tabs
            // m.lists[m.activeTab].ResetFilter()
            // m.lists[m.activeTab].GotoTop()
		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
             // Optional: Reset filter or scroll when switching tabs
            // m.lists[m.activeTab].ResetFilter()
            // m.lists[m.activeTab].GotoTop()
		case "f1":
			m.showHelp = true
		case "enter":
			// Handle selection based on the active tab
			switch m.activeTab {
			case constants.TabOptions:
				selectedItem := m.lists[m.activeTab].SelectedItem()
                if selectedItem != nil {
                     optItem := selectedItem.(listui.ListItem) // Type assertion
                     // For --version, could execute and show output
                     if optItem.Title() == "--version" {
                          m.selectedCommand = "--version" // Store command name
                          m.previewCommand = "systemctl --version" // Full command string for output view
                          m.selectedUnit = "" // No unit
                          m.state = StateOutput // Go directly to output state
                          m.commandOutput = "Running 'systemctl --version'..." // Initial running message
                          // Execute the command
                          return m, system.SystemctlCommand(m.selectedCommand, m.selectedUnit) // Use the helper
                     } else if optItem.Title() == "-h" || optItem.Title() == "--help" {
                         // For help, could execute systemctl --help and show output
                         m.selectedCommand = "--help"
                         m.previewCommand = "systemctl --help"
                         m.selectedUnit = ""
                         m.state = StateOutput
                         m.commandOutput = "Running 'systemctl --help'..."
                         return m, system.SystemctlCommand(m.selectedCommand, m.selectedUnit)
                     }
                     // For other options, maybe show description or error?
                     m.commandOutput = fmt.Sprintf("Info for option: %s - %s (Execution not implemented)", optItem.Title(), optItem.Description())
                     m.state = StateOutput
                     return m, nil // Stay in output state
                }
                // If no item selected in Options tab
                m.commandOutput = "Select an option to see info or execute."
                m.state = StateOutput
                return m, nil


			case constants.TabCommands:
				selectedItem := m.lists[m.activeTab].SelectedItem()
				if selectedItem != nil {
                    cmdItem := selectedItem.(listui.ListItem) // Type assertion
                    m.selectedCommand = cmdItem.Title() // Store the selected command

                    // Construct the command to preview
                    // We need the selected unit from the Units tab.
                    // First, check if the command *requires* a unit.
                    needsUnit := true
                    switch m.selectedCommand {
                        case "status", "list-units", "list-timers", "list-sockets", "list-jobs", "list-dependencies", "list-unit-files", "is-active", "is-enabled", "is-failed":
                            needsUnit = false
                        default:
                            needsUnit = true // Assume most other commands need a unit
                    }


                    if needsUnit && m.selectedUnit == "" {
                        // If the command requires a unit but none is selected
                        m.commandOutput = fmt.Sprintf("Command '%s' requires a unit. Please select a unit first in the Units tab.", m.selectedCommand)
                        m.showHelp = false // Hide help if it was showing
                        m.state = StateOutput // Go directly to output state
                        return m, nil // Stay in output state until keypress
                    }

                    // Construct the preview command string
                    m.previewCommand = "systemctl " + m.selectedCommand
                    if m.selectedUnit != "" && needsUnit { // Only append unit if needed and available
                         m.previewCommand += " " + m.selectedUnit
                    } else if m.selectedUnit != "" && !needsUnit && (m.selectedCommand == "status" || m.selectedCommand == "is-active" || m.selectedCommand == "is-enabled" || m.selectedCommand == "is-failed") {
                        // Special case: status/is-active etc. *can* take a unit, so include it if selected
                         m.previewCommand += " " + m.selectedUnit
                    }
                    // else: command doesn't need a unit, selected unit ignored for preview string


                    m.state = StatePreview // Change state to preview
                    return m, nil // Stay in preview state

				}
                 // If no item selected in Commands tab
                m.commandOutput = "Select a command to preview or execute."
                m.state = StateOutput
                return m, nil


			case constants.TabUnits:
				selectedItem := m.lists[m.activeTab].SelectedItem()
				if selectedItem != nil {
                    unitItem := selectedItem.(listui.ListItem) // Type assertion
                    m.selectedUnit = unitItem.Title() // Store the selected unit name

                    // Optional: Show confirmation or indicate unit is selected
                    m.commandOutput = fmt.Sprintf("Unit '%s' selected.", m.selectedUnit)
                    m.state = StateOutput // Go to output state briefly
                    return m, nil // Stay in output state
				}
                 // If no item selected in Units tab
                m.commandOutput = "Select a unit using Enter."
                m.state = StateOutput
                return m, nil
			}

		}

	case tea.WindowSizeMsg:
        // Update window size and resize lists
		m.width, m.height = msg.Width, msg.Height
		// Recalculate and set list sizes based on available space
		overheadHeight := 4 // Adjust this as needed

		listItemsViewportHeight := m.height - overheadHeight
		if listItemsViewportHeight < 0 { listItemsViewportHeight = 0 }

		listTotalWidth := m.width

		for i := range m.lists {
			m.lists[i].SetSize(listTotalWidth, listItemsViewportHeight)
		}
		return m, nil // No command needed for size change

	// This case handles the message when a command finishes executing asynchronously.
    // It can arrive while we are in StateOutput.
	case messages.CommandFinishedMsg: // <--- Use messages.CommandFinishedMsg from the correct package
        // This message is handled in updateOutput, but included here for robustness
        // if the state flow was different. With current flow, it's handled in updateOutput.
        return updateOutput(m, msg)

	} // <--- Closing brace for the switch on msg type

	// Delegate key presses (and potentially other messages the list understands)
    // to the currently active list *only if* it wasn't a message handled above.
	var cmd tea.Cmd
	m.lists[m.activeTab], cmd = m.lists[m.activeTab].Update(msg)
	return m, cmd
}


// updatePreview handles messages when the command preview is shown.
func updatePreview(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// User confirms execution
            if m.previewCommand == "" {
                 // Should not happen if previewCommand was set correctly in updateBrowse
                 m.state = StateOutput
                 m.commandOutput = "Error: Cannot execute empty command."
                 return m, nil
            }

			m.state = StateOutput // Change state to show output view
            m.commandOutput = "Running '" + m.previewCommand + "'..." // Show a running message immediately

            // Execute the command asynchronously
            // Need to parse m.previewCommand into command and args
            // Use strings.Fields for basic parsing (handles spaces)
            parts := strings.Fields(m.previewCommand)
            if len(parts) == 0 {
                 m.commandOutput = "Error: No command to execute."
                 return m, nil // Stay in output state with error
            }
            cmd := parts[0] // e.g., "systemctl"
            args := parts[1:] // e.g., ["status", "nginx.service"]

            // Call the async command execution function
			return m, system.ExecuteCommandAsync(cmd, args...) // Return the async command


		case "esc":
			// User cancels preview
			m.state = StateBrowse // Go back to Browse state
			m.selectedCommand = "" // Clear command state
			m.previewCommand = ""
			// Keep selectedUnit
			return m, nil
		}
	    // <--- Add closing brace for the switch on msg.String()

    // Handle WindowSizeMsg in preview state as well
     case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        // No list to resize, but model dimensions are updated for rendering preview box
        return m, nil

	} // <--- Add closing brace for the switch on msg type
	return m, nil
}

// updateOutput handles messages when command output is shown.
func updateOutput(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle messages specific to output state
	switch msg := msg.(type) {
    // This case handles the message when the async command finishes
	case messages.CommandFinishedMsg: // <--- Use messages.CommandFinishedMsg from the correct package
        // Command has finished, update the output
		m.commandOutput = msg.Output
		if msg.Err != nil {
			// Prepend "Error:" or append error details clearly
			m.commandOutput += "\n--- ERROR ---\n" + msg.Err.Error() // Append error details
		}
        // No state change, stay in output state showing the result
		return m, nil

    case tea.KeyMsg:
        // Any keypress dismisses the output view
        m.state = StateBrowse // Go back to Browse
        m.selectedCommand = ""
        m.previewCommand = ""
        m.commandOutput = "" // Clear the output
        // Keep selectedUnit
        return m, nil

    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        // No list to resize, but model dimensions are updated for rendering output box
        return m, nil

	} // <--- Add closing brace for the switch on msg type
	// Do not delegate to list or handle other messages in output state
	return m, nil
}
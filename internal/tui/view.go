// package tui
package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"systemctltui/internal/constants"
	"systemctltui/internal/styles" // <--- Check this import path matches your module name
)

// View renders the TUI based on the current state.
func (m model) View() string {
	// If window size hasn't been set yet, render a loading message
	if m.width == 0 || m.height == 0 {
		return "Initializing..." // Or a spinner
	}

	// Render different views based on state
	switch m.state {
	case StateBrowse: // Use StateBrowse
		return renderBrowseView(m) // Use renderBrowseView
	case StatePreview:
		return renderPreviewView(m)
	case StateOutput:
		return renderOutputView(m)
	default:
		// Should not happen
		return "Unknown state."
	}
}

// renderBrowseView renders the standard tab/list view.
func renderBrowseView(m model) string { // Use renderBrowseView
	// Render the header (tabs)
	header := styles.RenderTabs(m.tabs, m.activeTab)

	// Render the body (the active list)
	// The list view should respect the size constraints set by SetSize in Update.
	body := m.lists[m.activeTab].View()

	// Render the footer (help/info)
	var footerText string
	if m.showHelp {
		footerText = "Help: Press any key to return."
	} else {
		footerText = "Tab/Shift+Tab: switch tabs | ↑/↓: navigate | F1: help | q: quit"
		if m.activeTab == constants.TabCommands {
			footerText += " | Enter: preview/run"
            // Add info about selected unit if any
            if m.selectedUnit != "" {
                footerText += fmt.Sprintf(" (Unit: %s)", m.selectedUnit)
            } else {
                 footerText += " (No unit selected - some commands may fail)"
            }

		} else if m.activeTab == constants.TabUnits {
            footerText += " | Enter: select unit" // Indicate Enter selects the unit
        } else if m.activeTab == constants.TabOptions {
             footerText += " | Enter: info/run" // Indicate Enter shows info for options
        }
         // Add currently selected unit regardless of tab (useful context)
        if m.selectedUnit != "" {
             // Avoid showing it twice if already in commands tab
             if m.activeTab != constants.TabCommands {
                 footerText += fmt.Sprintf(" | Selected Unit: %s", m.selectedUnit)
             }
        }
	}
	footer := styles.FooterStyle.Render(footerText)

	// Use lipgloss.JoinVertical to stack header, body, and footer explicitly.
	layout := lipgloss.JoinVertical(
		lipgloss.Left, // Align components to the left
		header,
		// Optional: Add a small vertical gap (1 line) between header and list
		// lipgloss.NewStyle().Height(1).Width(m.width).Render(""), // Ensure width matches terminal
		body, // The list's view output
		// Optional: Add a small vertical gap (1 line) between list and footer
		// lipgloss.NewStyle().Height(1).Width(m.width).Render(""), // Ensure width matches terminal
		footer,
	)

    // Optional: Wrap the entire layout
    // return styles.AppBoundaryStyle.Render(layout)

	return layout
}

// renderPreviewView renders the command preview screen.
func renderPreviewView(m model) string {
    // Style for the preview box
    previewStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#5A56E0")).
        Padding(1, 2)

    // Content of the preview
    previewContent := lipgloss.JoinVertical(
        lipgloss.Left,
        "Command Preview:",
        styles.TabActiveStyle.Render(m.previewCommand), // Style the command string
        "", // Empty line for spacing
        "Press Enter to Execute, Esc to Cancel",
    )

    // Render the content within the styled box
    renderedPreview := previewStyle.Render(previewContent)

    // Calculate padding needed to center the box
    previewWidth := lipgloss.Width(renderedPreview)
    previewHeight := lipgloss.Height(renderedPreview)

    hPadding := (m.width - previewWidth) / 2
    vPadding := (m.height - previewHeight) / 2

    if hPadding < 0 { hPadding = 0 } // Prevent negative padding
    if vPadding < 0 { vPadding = 0 }

    // Use lipgloss.Place to center the rendered box on the screen
    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center, // Horizontal alignment
        lipgloss.Center, // Vertical alignment
        renderedPreview,
    )
}
// renderOutputView renders the command output screen.
func renderOutputView(m model) string {
	// Style for the output box
	outputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5A56E0")).
		Padding(1, 2) // Add padding inside the border: 1 vert, 2 horiz

	// Command text line
	commandLine := styles.TabActiveStyle.Render("> " + m.previewCommand) // Style the command that was run

	// Calculate available height for the output content area within the box
	// Total screen height - (space taken by outputStyle's vertical padding/border) - commandLineHeight - footerHeight - buffer
	// --- Use hardcoded vertical overhead based on Padding(1,2) and 1-unit border ---
	// Vertical space = Top Padding (1) + Bottom Padding (1) + Top Border (1) + Bottom Border (1) = 4
	borderAndPaddingHeight := 4
	// -----------------------------------------------------------------------------

	// Use a fixed footer height estimate or calculate it if complex
	outputFooterEstimateHeight := 1 // "Press any key..." line

	availableOutputHeight := m.height - borderAndPaddingHeight - lipgloss.Height(commandLine) - outputFooterEstimateHeight - 1 // Subtract components and a small buffer

	if availableOutputHeight < 0 {
		availableOutputHeight = 0
	}

	// Style for the actual output text block - constraints its size
	outputContentStyle := lipgloss.NewStyle().
		// --- Use hardcoded horizontal overhead based on Padding(1,2) and 1-unit border ---
		// Horizontal space = Left Padding (2) + Right Padding (2) + Left Border (1) + Right Border (1) = 6
		MaxWidth(m.width - 6).
		// ---------------------------------------------------------------------------------
		MaxHeight(availableOutputHeight) // Max height based on calculation
	// Add word wrapping if needed
	// WordWrap(true) // You might need to import github.com/muesli/reflow/wordwrap for this


	// Render the output content, respecting the size constraints
	outputContent := outputContentStyle.Render(m.commandOutput)

	// Combine the command line and the output content vertically
	fullOutput := lipgloss.JoinVertical(
		lipgloss.Left, // Align left
		commandLine,
		outputContent,
	)

	// Footer instruction
	outputFooter := styles.FooterStyle.Render("Press any key to return.")

	// Stack the combined output and the footer within the box style
	boxContent := lipgloss.JoinVertical(
		lipgloss.Left, // Align left
		fullOutput,
		outputFooter,
	)

	// Render the box, which automatically wraps its content
	renderedOutputBox := outputStyle.Render(boxContent)

	// Calculate padding needed to center the box
	outputWidth := lipgloss.Width(renderedOutputBox)
	outputHeight := lipgloss.Height(renderedOutputBox)


	hPadding := (m.width - outputWidth) / 2
	vPadding := (m.height - outputHeight) / 2

	if hPadding < 0 {
		hPadding = 0
	}
	if vPadding < 0 {
		vPadding = 0
	}


	// Place the rendered box centered on the screen
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center, // Horizontal alignment
		lipgloss.Center, // Vertical alignment
		renderedOutputBox,
	)
}

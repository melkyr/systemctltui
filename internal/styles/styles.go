// package styles
package styles

import (
	"strings" // fmt is removed

	"github.com/charmbracelet/lipgloss"
)

var (
	// TabActiveStyle is the lipgloss style for the active tab.
	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")). // Example active color
			// Correctly set the bottom border characters from the ThickBorder style
			Border(lipgloss.Border{Bottom: lipgloss.ThickBorder().Bottom}).
			BorderForeground(lipgloss.Color("#7D56F4")).         // Set the border color
			UnsetBorderTop(). // Ensure no top border
			UnsetBorderLeft(). // Ensure no left border
			UnsetBorderRight(). // Ensure no right border
			PaddingBottom(0) // Adjust padding to keep text close to the border

	// TabInactiveStyle is the lipgloss style for inactive tabs.
	TabInactiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A1A1A1")) // Example inactive color

	// AppBoundaryStyle defines the overall border for the application view.
	// You could wrap the whole view in this.
	AppBoundaryStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#5A56E0"))

	// FooterStyle for the help/info text at the bottom.
	FooterStyle = lipgloss.NewStyle().
			PaddingTop(1).
			Foreground(lipgloss.Color("#A1A1A1"))
)


// RenderTabs renders the tab bar string.
func RenderTabs(tabs []string, active int) string {
	var b strings.Builder
	for i, t := range tabs {
		style := TabInactiveStyle
		if i == active {
			style = TabActiveStyle
		}
		// Added a space inside the active tab rendering for visual padding if desired
		// b.WriteString(style.Render(fmt.Sprintf(" %s ", t))) // Example with padding
		b.WriteString(style.Render(t)) // Render the text with the chosen style
		b.WriteString(" ") // Space between tabs
	}
	return b.String()
}

// You can add other rendering helper functions here
// func RenderHelp() string { ... }
// func RenderPreviewBox(content string) string { ... }

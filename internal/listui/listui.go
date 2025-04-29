// package listui
package listui

import (
	"fmt"
	"log"
	//"strings"

	"github.com/charmbracelet/bubbles/list"
	//"github.com/charmbracelet/lipgloss" // Needed for list styles
	"systemctltui/internal/system"
)

// ListItem implements list.Item and holds a system.Unit.
// Exported because it's used in tui/model.
type ListItem struct {
	Unit system.Unit // Embed the Unit data - system.Unit is already exported
}

// Title returns the unit name for the list item title.
func (i ListItem) Title() string { return i.Unit.Name }

// Description returns a formatted string of unit status and description for the list item description.
func (i ListItem) Description() string {
	return fmt.Sprintf("[%s/%s/%s] %s", i.Unit.Load, i.Unit.Active, i.Unit.Sub, i.Unit.Description)
}

// FilterValue returns the unit name for filtering.
func (i ListItem) FilterValue() string { return i.Unit.Name }


// SimpleListItem implements list.Item for static lists (Options, Commands, Filters).
// Exported because it's used in tui/model for the filter list items.
type SimpleListItem struct { // <--- Exported struct name
    TitleValue string // <--- Exported field name
    DescValue string // <--- Exported field name
}

// Implement the list.Item interface for SimpleListItem
func (i SimpleListItem) Title() string { return i.TitleValue } // <--- Uses exported field
func (i SimpleListItem) Description() string { return i.DescValue } // <--- Uses exported field
func (i SimpleListItem) FilterValue() string { return i.TitleValue } // <--- Uses exported field


// CreateList creates a new bubbletea list model from a slice of ListItems (Units).
// Exported because it's used in NewLists.
func CreateList(items []ListItem) list.Model {
	const width, height = 60, 20 // Example fixed size
	delegate := list.NewDefaultDelegate()

	listItems := convert(items)
	// Provide a placeholder if the list is empty after conversion
	if len(listItems) == 0 {
		// Use ListItem for the placeholder for consistency
		listItems = []list.Item{ListItem{Unit: system.Unit{Name: "No units found", Description: "Try refreshing or check systemctl status."}}}
	}


	l := list.New(listItems, delegate, width, height)

	l.SetShowTitle(false) // Hide the list's default title

	l.SetShowPagination(true)
	l.SetFilteringEnabled(true) // Filtering is useful for units
	l.SetShowStatusBar(true)

	// Customize list styles if needed (example)
	// l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	// l.Styles.Item.Normal = lipgloss.NewStyle().PaddingLeft(2)
    // l.Styles.Item.Selected = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Bold(true)


	return l
}

// CreateSimpleList creates a new bubbletea list model from a slice of SimpleListItems.
// This is for static lists like Options, Commands, and Filters.
// Exported because it's used in tui/model.
func CreateSimpleList(items []SimpleListItem) list.Model { // <--- Exported function name
    const width, height = 30, 10 // Smaller size for dialogs/static lists
    delegate := list.NewDefaultDelegate()

    // Convert SimpleListItem slice to a list.Item slice
    converter := func(s []SimpleListItem) []list.Item {
        out := make([]list.Item, len(s))
        for i, it := range s {
            out[i] = list.Item(it) // Convert SimpleListItem (which implements list.Item)
        }
        return out
    }

    l := list.New(converter(items), delegate, width, height)
    l.SetShowTitle(false)
    l.SetShowPagination(false) // No pagination needed for small static lists
    l.SetFilteringEnabled(true) // Filtering is useful for commands/filters
    l.SetShowStatusBar(false) // No status bar needed

    // Customize styles for simple lists if they differ (example)
    // l.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)

    return l
}


// convert converts a slice of ListItem to a slice of list.Item.
// Unexported as it's only used internally.
func convert(source []ListItem) []list.Item {
	out := make([]list.Item, len(source))
	for i, it := range source {
		out[i] = list.Item(it)
	}
	return out
}

// InitOptionsList creates the list for global options.
// Exported because it's used in NewLists.
func InitOptionsList() list.Model {
	items := []SimpleListItem{ // Use exported SimpleListItem
		{TitleValue: "-h, --help", DescValue: "Show help text"}, // Use exported fields
		{TitleValue: "--version", DescValue: "Show version"}, // Use exported fields
	}
    return CreateSimpleList(items) // Use exported CreateSimpleList
}

// InitCommandsList creates the list for systemctl commands.
// Exported because it's used in NewLists.
func InitCommandsList() list.Model {
	items := []SimpleListItem{ // Use exported SimpleListItem
		{TitleValue: "status", DescValue: "Show unit status and logs"}, // Use exported fields
		{TitleValue: "start", DescValue: "Start one or more units"}, // Use exported fields
		{TitleValue: "stop", DescValue: "Stop one or more units"}, // Use exported fields
		{TitleValue: "restart", DescValue: "Restart one or more units"}, // Use exported fields
		{TitleValue: "enable", DescValue: "Enable one or more units"}, // Use exported fields
		{TitleValue: "disable", DescValue: "Disable one or more units"}, // Use exported fields
		// Add more commands here
	}
    return CreateSimpleList(items) // Use exported CreateSimpleList
}


// InitUnitsList fetches units from the system and creates the list of ListItems (Units).
// Exported because it's used in NewLists.
func InitUnitsList() list.Model {
	// Fetch units from the system
	units, err := system.FetchUnits() // system.FetchUnits is already exported
	if err != nil {
		log.Printf("Error fetching units: %v", err)
		// Create a list item to show the error
        // Using the Unit struct for consistency, even for errors
        errorUnit := system.Unit{ // system.Unit is already exported
            Name: "Error",
            Description: fmt.Sprintf("Failed to fetch units: %v", err),
            Load: "error", Active: "error", Sub: "error", Type: "error",
        }
		return CreateList([]ListItem{{Unit: errorUnit}}) // Use exported CreateList and ListItem
	}

	// Convert fetched Units to ListItems
	items := make([]ListItem, 0, len(units)) // Use exported ListItem
	for _, unit := range units {
		items = append(items, ListItem{Unit: unit}) // Use exported ListItem
	}

	// Create and return the initial list.Model containing all units
	return CreateList(items) // Use exported CreateList
}

// NewLists initializes all the necessary lists for the application.
// It fetches units (which are stored in the model afterwards)
// and returns the initial list models for the tabs.
// Exported because it's used in tui/model.
func NewLists() []list.Model { // <--- Exported function name
    // InitUnitsList returns the initial list model for the Units tab.
    // The full list of units fetched by InitUnitsList will be extracted
    // and stored separately in the model in NewModel.
	unitsListModel := InitUnitsList() // InitUnitsList is exported

	return []list.Model{
		InitOptionsList(),    // InitOptionsList is exported
		InitCommandsList(),   // InitCommandsList is exported
		unitsListModel,       // The list.Model itself is a type from an external package
	}
}


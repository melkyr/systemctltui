// package listui
package listui

import (
	"fmt"
	"log"
	//"strings" // Import strings for parsing unit type

	"github.com/charmbracelet/bubbles/list"
	"systemctltui/internal/system" // Import the enhanced system package
)

// ListItem implements list.Item and holds a system.Unit.
type ListItem struct {
	Unit system.Unit // Embed the Unit data
}

// Title returns the unit name for the list item title.
func (i ListItem) Title() string { return i.Unit.Name }

// Description returns a formatted string of unit status and description for the list item description.
func (i ListItem) Description() string {
	// Format Load, Active, Sub, and Description nicely
	return fmt.Sprintf("[%s/%s/%s] %s", i.Unit.Load, i.Unit.Active, i.Unit.Sub, i.Unit.Description)
}

// FilterValue returns the unit name for filtering. You might add type or description later.
func (i ListItem) FilterValue() string { return i.Unit.Name }


// simpleListItem implements list.Item for static lists (Options, Commands).
// Defined at package level, unexported as it's internal to this package.
type simpleListItem struct { // <--- Moved to package level
    title, desc string
}

// Implement the list.Item interface for simpleListItem
func (i simpleListItem) Title() string { return i.title } // <--- Moved to package level
func (i simpleListItem) Description() string { return i.desc } // <--- Moved to package level
func (i simpleListItem) FilterValue() string { return i.title } // <--- Moved to package level


// CreateList creates a new bubbletea list model from a slice of ListItems.
// Update function signature - it still takes []ListItem
func CreateList(items []ListItem) list.Model {
	const width, height = 60, 20 // Example fixed size
	delegate := list.NewDefaultDelegate()

	// Check if items is empty and provide a placeholder if necessary
	listItems := convert(items)
	if len(listItems) == 0 {
		listItems = []list.Item{ListItem{Unit: system.Unit{Name: "No units found", Description: "Try refreshing or check systemctl status."}}}
	}


	l := list.New(listItems, delegate, width, height)

	l.SetShowTitle(false) // Hide the list's default title

	l.SetShowPagination(true)
	l.SetFilteringEnabled(true) // Filtering is useful for units
	l.SetShowStatusBar(true)

	// Customize list styles if needed
	// l.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF"))
	// l.Styles.Item.Normal = lipgloss.NewStyle().PaddingLeft(2)
    // l.Styles.Item.Selected = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Bold(true)


	return l
}

// convert converts a slice of ListItem to a slice of list.Item.
// Update function signature - it still takes []ListItem
func convert(source []ListItem) []list.Item {
	out := make([]list.Item, len(source))
	for i, it := range source {
		out[i] = list.Item(it) // Correctly convert ListItem (which implements list.Item)
	}
	return out
}

// InitOptionsList creates the list for global options.
func InitOptionsList() list.Model {
    // simpleListItem and its methods are now defined at package level

	items := []simpleListItem{ // simpleListItem is now visible here
		{title: "-h, --help", desc: "Show help text"},
		{title: "--version", desc: "Show version"},
		// Add more options here as you implement them
	}

    // Convert simpleListItem slice to a list.Item slice
    simpleConverter := func(s []simpleListItem) []list.Item {
        out := make([]list.Item, len(s))
        for i, it := range s {
            out[i] = list.Item(it) // Convert simpleListItem (which now implements list.Item)
        }
        return out
    }

	const width, height = 60, 20 // Match list size
	delegate := list.NewDefaultDelegate()
	l := list.New(simpleConverter(items), delegate, width, height)
    l.SetShowTitle(false)
    l.SetShowPagination(false) // No pagination for static lists
    l.SetFilteringEnabled(false) // Filtering not typically needed for these few options
    l.SetShowStatusBar(false) // No status bar

    return l
}

// InitCommandsList creates the list for systemctl commands.
func InitCommandsList() list.Model {
    // simpleListItem and its methods are now defined at package level

	items := []simpleListItem{ // simpleListItem is now visible here
		{title: "status", desc: "Show unit status and logs"},
		{title: "start", desc: "Start one or more units"},
		{title: "stop", desc: "Stop one or more units"},
		{title: "restart", desc: "Restart one or more units"},
		{title: "enable", desc: "Enable one or more units"},
		{title: "disable", desc: "Disable one or more units"},
		// Add more commands here
	}
    // Convert simpleListItem slice to a list.Item slice
    simpleConverter := func(s []simpleListItem) []list.Item {
        out := make([]list.Item, len(s))
        for i, it := range s {
            out[i] = list.Item(it) // Convert simpleListItem (which now implements list.Item)
        }
        return out
    }

	const width, height = 60, 20 // Match list size
	delegate := list.NewDefaultDelegate()
	l := list.New(simpleConverter(items), delegate, width, height)
    l.SetShowTitle(false)
    l.SetShowPagination(false) // No pagination for static lists
    l.SetFilteringEnabled(true) // Filtering IS useful for commands!
    l.SetShowStatusBar(true) // Status bar is useful with filtering

    return l
}


// InitUnitsList fetches units from the system and creates the list of ListItems.
func InitUnitsList() list.Model {
	// Call the enhanced function from the system package
	units, err := system.FetchUnits()
	if err != nil {
		log.Printf("Error fetching units: %v", err)
		// Create a list item to show the error
        // Using the Unit struct for consistency, even for errors
        errorUnit := system.Unit{
            Name: "Error",
            Description: fmt.Sprintf("Failed to fetch units: %v", err),
            Load: "error", Active: "error", Sub: "error", Type: "error",
        }
		return CreateList([]ListItem{{Unit: errorUnit}}) // <-- Use the updated ListItem struct
	}

	items := make([]ListItem, 0, len(units)) // <-- Use the updated ListItem struct
	for _, unit := range units {
		// Create a ListItem for each fetched Unit
		items = append(items, ListItem{Unit: unit}) // <-- Embed the system.Unit
	}

	return CreateList(items) // Returns a list.Model
}

// NewLists initializes all the necessary lists for the application.
// This function can be called from the main model setup.
// It will now fetch and return lists based on the enhanced Unit data.
func NewLists() []list.Model {
	// Fetching units can take time. Consider showing a loading indicator later.
	unitsList := InitUnitsList()

	return []list.Model{
		InitOptionsList(),    // Index 0, matches constants.TabOptions
		InitCommandsList(),   // Index 1, matches constants.TabCommands
		unitsList,            // Index 2, matches constants.TabUnits
	}
}

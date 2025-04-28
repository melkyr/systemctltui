// package listui
package listui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	"systemctltui/internal/system"
)

// ListItem implements list.Item. <-- Renamed from listItem
type ListItem struct { // <-- Renamed from listItem
	title, desc string
}

// Update methods to use ListItem
func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }

// CreateList creates a new bubbletea list model from a slice of ListItems.
func CreateList(items []ListItem) list.Model {
	// You might want to make dimensions configurable or dynamic later
	const width, height = 60, 20 // Example fixed size
	delegate := list.NewDefaultDelegate()

	l := list.New(convert(items), delegate, width, height)

	// --- Add this line to hide the list's default title ---
	l.SetShowTitle(false)
	// ----------------------------------------------------

	l.SetShowPagination(true)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)
	// l.SetFilteringEnabled(true) // This line is redundant, it's already true above

	return l
}

// convert converts a slice of ListItem to a slice of list.Item.
// Update function signature
func convert(source []ListItem) []list.Item { // <-- Use ListItem
	out := make([]list.Item, len(source))
	for i, it := range source {
		out[i] = it
	}
	return out
}

// InitOptionsList creates the list for global options.
// Update return type if needed, but it returns list.Model so it's fine
func InitOptionsList() list.Model {
	items := []ListItem{ // <-- Use ListItem
		{title: "-h, --help", desc: "Show help text"},
		{title: "--version", desc: "Show version"},
	}
	return CreateList(items)
}

// InitCommandsList creates the list for systemctl commands.
// Update return type if needed
func InitCommandsList() list.Model {
	items := []ListItem{ // <-- Use ListItem
		{title: "status", desc: "Show unit status and logs"},
		{title: "start", desc: "Start one or more units"},
		{title: "stop", desc: "Stop one or more units"},
		{title: "restart", desc: "Restart one or more units"},
		{title: "enable", desc: "Enable one or more units"},
		{title: "disable", desc: "Disable one or more units"},
	}
	return CreateList(items)
}

// InitUnitsList fetches units from the system and creates the list.
// Update return type if needed
func InitUnitsList() list.Model {
	units, err := system.FetchUnits()
	if err != nil {
		log.Printf("Error fetching units: %v", err)
		// Use ListItem here as well
		return CreateList([]ListItem{{title: "Error", desc: fmt.Sprintf("Failed to fetch units: %v", err)}}) // <-- Use ListItem
	}

	items := make([]ListItem, 0, len(units)) // <-- Use ListItem
	for _, unit := range units {
		items = append(items, ListItem{title: unit.Name, desc: ""}) // <-- Use ListItem
	}

	return CreateList(items)
}

// NewLists initializes all the necessary lists for the application.
func NewLists() []list.Model {
	unitsList := InitUnitsList()

	return []list.Model{
		InitOptionsList(),
		InitCommandsList(),
		unitsList,
	}
}

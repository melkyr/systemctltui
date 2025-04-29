// package tui
package tui

import (
	"sort"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"systemctltui/internal/listui"
	"systemctltui/internal/constants"
	"systemctltui/internal/system"
)

// AppState defines the different states of the application view.
type AppState int

const (
	StateBrowse AppState = iota // Browse tabs and lists
	StatePreview                  // Showing command preview
	StateOutput                   // Showing command output/error
	StateFiltering                // Showing unit filter dialog
)

// model represents the main state of the TUI application.
type model struct {
	activeTab int
	tabs      []string
	lists     []list.Model
	showHelp  bool
	width, height int

	// State for command preview and execution
	state           AppState
	selectedCommand string
	selectedUnit    string
	previewCommand  string
	commandOutput   string

	// State for unit filtering
	FullUnitList      []system.Unit
	filterList        list.Model
	currentUnitFilter string
}

// NewModel initializes the main application model.
func NewModel() model {
	tabs := []string{"Global Options", "Commands", "Units"}
	lists := listui.NewLists()

	// Extract the full list of units
	var fullUnitList []system.Unit
	unitsListModel := lists[constants.TabUnits]
	for _, item := range unitsListModel.Items() {
		if listItem, ok := item.(listui.ListItem); ok {
			fullUnitList = append(fullUnitList, listItem.Unit)
		}
	}

	// Determine unique unit types
	uniqueTypes := make(map[string]bool)
	uniqueTypes["All"] = true
	for _, unit := range fullUnitList {
		if unit.Type != "" {
			uniqueTypes[unit.Type] = true
		}
	}

	// Build filter items
	var filterItems []listui.SimpleListItem
	for unitType := range uniqueTypes {
		filterItems = append(filterItems, listui.SimpleListItem{
			TitleValue: unitType,
			DescValue:  "",
		})
	}

	// Sort filter items alphabetically
	sort.Slice(filterItems, func(i, j int) bool {
		return filterItems[i].TitleValue < filterItems[j].TitleValue
	})

	// Create the filter list model
	filterList := listui.CreateSimpleList(filterItems)

	return model{
		activeTab:        constants.TabOptions,
		tabs:             tabs,
		lists:            lists,
		showHelp:         false,
		width:            0,
		height:           0,
		state:            StateBrowse,

		FullUnitList:      fullUnitList,
		filterList:        filterList,
		currentUnitFilter: "All",
	}
}

// Init performs initial setup for the bubbletea model.
func (m model) Init() tea.Cmd {
	return nil
}

// Update and View methods are defined in update.go and view.go


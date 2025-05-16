package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// MenuOption represents a menu option with a key and description
type MenuOption struct {
	Key         string
	Description string
}

// CreateMenuBar creates a menu bar with the given options
func CreateMenuBar(options []MenuOption) *tview.TextView {
	menuText := tview.NewTextView()
	var menuString strings.Builder

	for i, option := range options {
		if i > 0 {
			menuString.WriteString("   ")
		}
		var menuItem = formatMenuOption(option)
		menuString.WriteString(menuItem)
	}

	menuText.SetText(menuString.String())
	return menuText
}

// formatMenuOption formats a single menu option for display in the menu bar.
func formatMenuOption(option MenuOption) string {
	menuItem := fmt.Sprintf("[white]%s[d:] %s", option.Key, option.Description)
	return menuItem
}

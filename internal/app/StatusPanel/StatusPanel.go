package StatusPanel

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Create creates a panel that displays the game statusbar
func Create(status string, borderColor tcell.Color, backgroundColor tcell.Color) *tview.Flex {
	statusPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create statusbar text view
	statusTextView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(status)

	// Add the text view to the panel
	statusPanel.AddItem(statusTextView, 1, 0, false)

	// Set the border and background
	statusPanel.SetBorder(true)
	statusPanel.SetBorderColor(borderColor)
	statusPanel.SetBackgroundColor(backgroundColor)

	return statusPanel
}

// UpdateWithGameTime updates the status panel to include the total game time
func UpdateWithGameTime(panel *tview.Flex, status string, totalGameTime time.Duration) {
	statusTextView := panel.GetItem(0).(*tview.TextView)
	statusTextView.SetText(fmt.Sprintf("%s | Total Game Time: %v", status, totalGameTime))
}

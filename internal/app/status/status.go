package status

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Status creates a panel that displays the game status
func Status(status string, borderColor tcell.Color, backgroundColor tcell.Color) *tview.Flex {
	statusPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create status text view
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

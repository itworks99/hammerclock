package hammerclock

import (
	"github.com/rivo/tview"
)

// CreateEndGameConfirmationModal creates a modal dialog asking for confirmation to end the game
func CreateEndGameConfirmationModal(view *View) *tview.Modal {
	modal := tview.NewModal().
		SetText("Would you like to end the current game?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 { // "Yes" is the first button (index 0)
				view.MessageChan <- &endGameConfirmMsg{Confirmed: true}
			} else {
				view.MessageChan <- &endGameConfirmMsg{Confirmed: false}
			}
		})

	// Style the modal
	modal.SetBorder(true)
	modal.SetTitle(" Confirm End Game ")

	return modal
}

// ShowConfirmationModal displays a confirmation modal in the application
func (v *View) ShowConfirmationModal(modal *tview.Modal) {
	// Center the modal in a flex container
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(modal, 10, 1, true).
				AddItem(nil, 0, 1, false),
			60, 1, true,
		).
		AddItem(nil, 0, 1, false)

	// Create a new pages object to layer the modal over the main UI
	pages := tview.NewPages().
		AddPage("background", v.MainFlex, true, true).
		AddPage("modal", flex, true, true)

	// Set the pages as the application's root
	v.App.SetRoot(pages, true)
}

// RestoreMainUI restores the main UI after modal is closed
func (v *View) RestoreMainUI() {
	v.App.SetRoot(v.MainFlex, true)
}

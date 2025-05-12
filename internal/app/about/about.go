package about

import (
	"io"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mdp/qrterminal/v3"
	"github.com/rivo/tview"
	"hammerclock/config"
)

func About(mainColor tcell.Color) *tview.Flex {
	aboutPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create content with about information
	contentBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(mainColor).
		SetDynamicColors(true)

	// Use a string.Builder to capture the QR code output
	renderStr := new(strings.Builder)
	config := qrterminal.Config{
		Level:      qrterminal.M,
		Writer:     io.Writer(renderStr),
		HalfBlocks: true,
	}
	qrterminal.GenerateWithConfig(hammerclockConfig.GitHubUrl, config)

	// Build about content
	var content strings.Builder

	content.WriteString(renderStr.String() + "\n")
	content.WriteString("[d:]v." + hammerclockConfig.Version + "\n\n")
	content.WriteString("A terminal-based timer and phase tracker for tabletop games\n\n")
	content.WriteString(hammerclockConfig.GitHubUrl + "\n\n\n\n")
	content.WriteString("Press [white]A[d:] to return to the main screen")

	contentBox.SetText(content.String())

	// Add the boxes to the panel
	aboutPanel.AddItem(tview.NewBox(), 1, 0, false). // Spacer
								AddItem(contentBox, 0, 1, false)

	// Set the border and background
	aboutPanel.SetBorder(true)
	aboutPanel.SetTitle(" About ")
	aboutPanel.SetBorderColor(tcell.ColorYellow)
	aboutPanel.SetBackgroundColor(tcell.ColorBlack)

	return aboutPanel
}

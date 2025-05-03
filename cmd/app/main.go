package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"hammerclock/internal/app"
	"hammerclock/pkg/settings"
)

// CLI usage information
var cliUsage = `
Hammerclock - A terminal-based timer and phase tracker for tabletop games

Usage:
  hammerclock [options]

Options:
  -s <file>    Specify a custom settings file (default: defaultRules.json)
  -h, --help   Show this help message

Examples:
  hammerclock                     # Run with default settings
  hammerclock -s mySettings.json # Run with custom settings
`

func main() {
	// Parse command line flags
	settingsFileFlag := flag.String("s", "defaultRules.json", "Path to the settings file")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, cliUsage)
	}
	flag.Parse()

	// Load settings from file
	settings := settings.LoadSettings(*settingsFileFlag)

	// Create the model
	model := app.NewModel()
	model.Settings = settings
	model.Phases = settings.Phases
	model.CurrentColorPalette = app.GetColorPaletteByName(settings.ColorPalette)

	// Create players based on settings
	players := make([]*app.Player, settings.PlayerCount)
	for i := 0; i < settings.PlayerCount; i++ {
		playerName := fmt.Sprintf("Player %d", i+1) // Default name as fallback
		if i < len(settings.PlayerNames) {
			playerName = settings.PlayerNames[i]
		}
		players[i] = &app.Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
		}
	}
	model.Players = players

	// Create the view
	view := app.NewView(&model)

	// Set up message channel for communication between components
	msgChan := make(chan app.Message)
	done := make(chan struct{})

	// Set up input capture to send key press messages
	app.SetupInputCapture(view.App, msgChan)

	// Start the ticker to send tick messages
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Update the clock display
				view.App.QueueUpdateDraw(func() {
					view.UpdateClock()
				})

				// Send a tick message
				msgChan <- &app.TickMsg{}
			case <-done:
				return
			}
		}
	}()

	// Start the main update loop
	go func() {
		for {
			select {
			case msg := <-msgChan:
				// Process the message and update the model
				updatedModel, cmd := app.Update(msg, model)
				model = updatedModel

				// Render the updated model
				view.App.QueueUpdateDraw(func() {
					view.Render(&model)
				})

				// Execute the command if there is one
				if cmd != nil {
					go func() {
						// Execute the command and send any resulting message
						if resultMsg := cmd(); resultMsg != nil {
							msgChan <- resultMsg
						}
					}()
				}
			case <-done:
				return
			}
		}
	}()

	// Start the application
	if err := view.App.SetRoot(view.MainFlex, true).EnableMouse(true).Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
	}

	// Clean up when the application exits
	close(done)
}

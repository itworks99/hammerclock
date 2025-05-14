package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"hammerclock/components"
	"hammerclock/components/hammerclock/Palette"
	hammerclockConfig "hammerclock/config"
	"hammerclock/internal/app"
	"hammerclock/internal/app/LogPanel"
)

// CLI usage information
var cliUsage = `
Hammerclock - A terminal-based timer and phase tracker for tabletop games

Usage:
  hammerclock [options]

options:
  -o <file>    Specify a custom options file (default: defaultRules.json)
  -h, --help   Show this help message

Examples:
  hammerclock                     # Run with default options
  hammerclock -o myOptions.json   # Run with custom options
`

// main function initializes the application, loads options, and starts the main loop
func main() {
	// Parse command line flags
	optionsFileFlag := flag.String("o", hammerclockConfig.DefaultOptionsFilename, "Path to the options file")
	flag.Usage = func() {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, cliUsage)
	}
	flag.Parse()

	// Load options from file
	options := components.LoadOptions(*optionsFileFlag)

	// Create the model
	model := app.NewModel()
	model.Options = options
	model.Phases = options.Rules[options.Default].Phases
	model.CurrentColorPalette = Palette.GetColorPaletteByName(options.ColorPalette)

	// Create players based on options
	players := make([]*app.Player, options.PlayerCount)
	for i := 0; i < options.PlayerCount; i++ {
		playerName := fmt.Sprintf("Player %d", i+1) // Default name as fallback
		if i < len(options.PlayerNames) {
			playerName = options.PlayerNames[i]
		}
		players[i] = &app.Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
			TurnCount:    0,
			ActionLog:    []LogPanel.LogEntry{}, // Initialize empty action log
		}

		// Add initial player log message
		if i == 0 {
			app.AddLogEntry(players[i], "Player initialized as active player (Turn %d)", players[i].TurnCount)
		} else {
			app.AddLogEntry(players[i], "Player initialized (Turn %d)", players[i].TurnCount)
		}
	}
	model.Players = players

	// Set up message channel for communication between components
	msgChan := make(chan app.Message)
	done := make(chan struct{})

	// Create the view
	view := app.NewView(&model, msgChan)

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
					view.UpdateClock(&model)
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

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"hammerclock/internal/hammerclock"
	"hammerclock/internal/hammerclock/common"
	"hammerclock/internal/hammerclock/config"
	"hammerclock/internal/hammerclock/logging"
	"hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/palette"
)

// CLI usage information
var cliUsage = `
Hammerclock ` + hammerclockConfig.Version + `
Terminal-based timer and phase tracker for tabletop games

Usage:
  hammerclock [options]

options:
  -o <file>    Specify a custom options file (default: default.json)
  -h, --help   Show this help message

Examples:
  hammerclock                     # Run with default options
  hammerclock -o myOptions.json   # Run with custom options
`

func main() {
	// Initialize logging
	logging.Initialise()
	fmt.Println("Hammerclock", hammerclockConfig.Version, "starting up...")
	fmt.Println("Logs will be written to logs.csv in the current directory")

	// Parse command line flags
	optionsFileFlag := flag.String("o", hammerclockConfig.DefaultOptionsFilename, "Path to the loadedOptions file")
	flag.Usage = func() {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, cliUsage)
	}
	flag.Parse()

	// Load loadedOptions from file
	loadedOptions := options.LoadOptions(*optionsFileFlag)

	// CreateAboutPanel the model
	model := hammerclock.NewModel()
	model.Options = loadedOptions
	model.Phases = loadedOptions.Rules[loadedOptions.Default].Phases
	model.CurrentColorPalette = palette.ColorPaletteByName(loadedOptions.ColorPalette)

	// CreateAboutPanel players based on loadedOptions
	players := make([]*common.Player, loadedOptions.PlayerCount)
	for i := 0; i < loadedOptions.PlayerCount; i++ {
		playerName := fmt.Sprintf("Player %d", i+1)
		if i < len(loadedOptions.PlayerNames) {
			playerName = loadedOptions.PlayerNames[i]
		}
		players[i] = &common.Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
			TurnCount:    0,
			ActionLog:    []common.LogEntry{},
		}
	}
	model.Players = players

	// Set up message channel for communication between components
	msgChan := make(chan common.Message)
	done := make(chan struct{})

	// CreateAboutPanel the view
	view := hammerclock.NewView(&model, msgChan)

	// Set up input capture to send key press messages
	hammerclock.SetupInputCapture(view.App, msgChan)

	// Start the ticker to send tick messages
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Only update if the game is started
				if model.GameStarted {
					view.App.QueueUpdateDraw(func() {
						view.UpdateClock(&model)
					})
				}
				msgChan <- &common.TickMsg{}
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
				updatedModel, cmd := hammerclock.Update(msg, model)
				model = updatedModel

				_ = options.SaveOptions(model.Options, "", true)

				view.App.QueueUpdateDraw(func() {
					view.Render(&model)
				})

				if cmd != nil {
					go func() {
						if resultMsg := cmd(); resultMsg != nil {
							// Special handling for ShowModalMsg
							if showModal, ok := resultMsg.(*common.ShowModalMsg); ok {
								view.App.QueueUpdateDraw(func() {
									switch showModal.Type {
									case "EndGameConfirm":
										modal := hammerclock.CreateEndGameConfirmationModal(view)
										view.ShowConfirmationModal(modal)
									}
								})
							} else if _, ok := resultMsg.(*common.RestoreMainUIMsg); ok {
								view.App.QueueUpdateDraw(func() {
									view.RestoreMainUI()
								})
							} else {
								msgChan <- resultMsg
							}
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
	logging.Cleanup()
}

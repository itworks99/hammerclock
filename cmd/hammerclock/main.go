package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"hammerclock/internal/hammerclock"
	"hammerclock/internal/hammerclock/config"
	"hammerclock/internal/hammerclock/logging"
	options2 "hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/palette"
	"hammerclock/internal/hammerclock/ui"
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
	optionsFileFlag := flag.String("o", hammerclockConfig.DefaultOptionsFilename, "Path to the options file")
	flag.Usage = func() {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, cliUsage)
	}
	flag.Parse()

	// Load options from file
	options := options2.LoadOptions(*optionsFileFlag)

	// CreateAboutPanel the model
	model := hammerclock.NewModel()
	model.Options = options
	model.Phases = options.Rules[options.Default].Phases
	model.CurrentColorPalette = palette.ColorPaletteByName(options.ColorPalette)

	// CreateAboutPanel players based on options
	players := make([]*hammerclock.Player, options.PlayerCount)
	for i := 0; i < options.PlayerCount; i++ {
		playerName := fmt.Sprintf("Player %d", i+1)
		if i < len(options.PlayerNames) {
			playerName = options.PlayerNames[i]
		}
		players[i] = &hammerclock.Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
			TurnCount:    0,
			ActionLog:    []ui.LogEntry{},
		}

		// Add initial player log message
		if i == 0 {
			createInitialLog(players[i], &model, "Initialized - active player")
		} else {
			createInitialLog(players[i], &model, "Initialized")
		}
	}
	model.Players = players

	// Set up message channel for communication between components
	msgChan := make(chan hammerclock.Message)
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
				msgChan <- &hammerclock.TickMsg{}
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

				_ = options2.SaveOptions(model.Options, "", true)

				view.App.QueueUpdateDraw(func() {
					view.Render(&model)
				})

				if cmd != nil {
					go func() {
						if resultMsg := cmd(); resultMsg != nil {
							// Special handling for ShowModalMsg
							if showModal, ok := resultMsg.(*hammerclock.ShowModalMsg); ok {
								view.App.QueueUpdateDraw(func() {
									switch showModal.Type {
									case "EndGameConfirm":
										modal := hammerclock.CreateEndGameConfirmationModal(view)
										view.ShowConfirmationModal(modal)
									}
								})
							} else if _, ok := resultMsg.(*hammerclock.RestoreMainUIMsg); ok {
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

// createInitialLog adds initial log entries for players at startup
func createInitialLog(player *hammerclock.Player, model *hammerclock.Model, format string, args ...any) {
	currentPhase := ""
	if player.CurrentPhase < len(model.Options.Rules[model.Options.Default].Phases) && player.CurrentPhase >= 0 {
		currentPhase = model.Options.Rules[model.Options.Default].Phases[player.CurrentPhase]
	}

	logEntry := ui.LogEntry{
		DateTime:   time.Now().Local().Format(hammerclockConfig.DefaultLogDateTimeFormat),
		PlayerName: player.Name,
		Turn:       player.TurnCount,
		Phase:      currentPhase,
		Message:    fmt.Sprintf(format, args...),
	}

	// Add to in-memory player action log for UI
	player.ActionLog = append(player.ActionLog, logEntry)

	// Always log initialization messages to CSV regardless of LoggingEnabled setting
	// This ensures the log file is created at startup even if logging is disabled
	logging.WriteLogEntry(logEntry, true)
}

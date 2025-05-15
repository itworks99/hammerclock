package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"hammerclock/components/hammerclock/Palette"
	"hammerclock/components/hammerclock/fileio"
	hammerclockConfig "hammerclock/config"
	"hammerclock/internal/app"
	logpanel "hammerclock/internal/app/LogPanel"
)

// Buffered channel for log entries
var logChannel = make(chan logpanel.LogEntry, 100)

// CLI usage information
var cliUsage = `
Hammerclock ` + hammerclockConfig.Version + `
Terminal-based timer and phase tracker for tabletop games

Usage:
  hammerclock [options]

options:
  -o <file>    Specify a custom options file (default: defaultRules.json)
  -h, --help   Show this help message

Examples:
  hammerclock                     # Run with default options
  hammerclock -o myOptions.json   # Run with custom options
`

func main() {
	var logWg sync.WaitGroup
	logWg.Add(1)
	// Start background log writer
	go func() {
		defer logWg.Done()
		for entry := range logChannel {
			logpanel.WriteLogEntry(entry)
		}
	}()

	// Parse command line flags
	optionsFileFlag := flag.String("o", hammerclockConfig.DefaultOptionsFilename, "Path to the options file")
	flag.Usage = func() {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, cliUsage)
	}
	flag.Parse()

	// Load options from file
	options := fileio.LoadOptions(*optionsFileFlag)

	// Create the model
	model := app.NewModel()
	model.Options = options
	model.Phases = options.Rules[options.Default].Phases
	model.CurrentColorPalette = Palette.ColorPaletteByName(options.ColorPalette)

	// Create players based on options
	players := make([]*app.Player, options.PlayerCount)
	for i := 0; i < options.PlayerCount; i++ {
		playerName := fmt.Sprintf("Player %d", i+1)
		if i < len(options.PlayerNames) {
			playerName = options.PlayerNames[i]
		}
		players[i] = &app.Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
			TurnCount:    0,
			ActionLog:    []logpanel.LogEntry{},
		}

		// Add initial player log message
		if i == 0 {
			AddLogEntryBuffered(players[i], &model, "Initialized - active player")
		} else {
			AddLogEntryBuffered(players[i], &model, "Initialized")
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
				// Only update if the game is started
				if model.GameStarted {
					view.App.QueueUpdateDraw(func() {
						view.UpdateClock(&model)
					})
				}
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
				updatedModel, cmd := app.Update(msg, model)
				model = updatedModel

				_ = fileio.SaveOptions(model.Options, "", true)

				view.App.QueueUpdateDraw(func() {
					view.Render(&model)
				})

				if cmd != nil {
					go func() {
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
	close(logChannel)
	logWg.Wait()
}

// AddLogEntryBuffered adds a log entry and writes to the log channel if enabled
func AddLogEntryBuffered(player *app.Player, model *app.Model, format string, args ...any) {
	currentPhase := ""
	if player.CurrentPhase < len(model.Options.Rules[model.Options.Default].Phases) && player.CurrentPhase >= 0 {
		currentPhase = model.Options.Rules[model.Options.Default].Phases[player.CurrentPhase]
	}

	logEntry := logpanel.LogEntry{
		DateTime:   time.Now().Local().Format(hammerclockConfig.DefaultLogDateTimeFormat),
		PlayerName: player.Name,
		Turn:       player.TurnCount,
		Phase:      currentPhase,
		Message:    fmt.Sprintf(format, args...),
	}

	player.ActionLog = append(player.ActionLog, logEntry)
	if model.Options.EnableCSVLog {
		select {
		case logChannel <- logEntry:
			// sent successfully
		default:
			// channel full, drop log entry to avoid UI lag
		}
	}
}

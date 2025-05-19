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
	logging.Initialise()
	fmt.Println("Hammerclock", hammerclockConfig.Version, "starting up...")
	fmt.Println("Logs will be written to logs.csv in the current directory")

	optionsFileFlag := flag.String("o", hammerclockConfig.DefaultOptionsFilename, "Path to the loadedOptions file")
	flag.Usage = func() {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, cliUsage)
	}
	flag.Parse()

	loadedOptions := options.LoadOptions(*optionsFileFlag)

	model := hammerclock.NewModel()
	model.Options = loadedOptions
	model.Phases = loadedOptions.Rules[loadedOptions.Default].Phases
	model.CurrentColorPalette = palette.ColorPaletteByName(loadedOptions.ColorPalette)

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

	msgChan := make(chan common.Message)
	done := make(chan struct{})

	view := hammerclock.NewView(&model, msgChan)
	hammerclock.SetupInputCapture(view.App, msgChan)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
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
							if showModal, ok := resultMsg.(*common.ShowModalMsg); ok {
								view.App.QueueUpdateDraw(func() {
									switch showModal.Type {
									case "EndGameConfirm":
										modal := hammerclock.CreateEndGameConfirmationModal(view)
										hammerclock.ShowConfirmationModal(view, modal)
									}
								})
							} else if _, ok := resultMsg.(*common.RestoreMainUIMsg); ok {
								view.App.QueueUpdateDraw(func() {
									view.RestoreMainView()
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

	if err := view.App.SetRoot(view.MainView, true).EnableMouse(true).Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
	}

	close(done)
	logging.Cleanup()
}

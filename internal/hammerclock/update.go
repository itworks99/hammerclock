package hammerclock

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"hammerclock/internal/hammerclock/common"
	"hammerclock/internal/hammerclock/logging"
	"hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/palette"
	"hammerclock/internal/hammerclock/rules"
)

// Command represents a Command that can be executed after an update
type Command func() common.Message

// noCommand is a Command that does nothing
func noCommand() common.Message {
	return nil
}

// Update processes a message and returns an updated model and a command to execute
func Update(msg common.Message, model common.Model) (common.Model, Command) {
	switch msg := msg.(type) {
	case *common.StartGameMsg:
		return handleStartGame(model)
	case *common.EndGameMsg:
		return handleEndGame(model)
	case *common.EndGameConfirmMsg:
		return handleEndGameConfirm(msg, model)
	case *common.ShowEndGameConfirmMsg:
		return handleShowEndGameConfirm(model)
	case *common.SwitchTurnsMsg:
		return handleSwitchTurns(model)
	case *common.NextPhaseMsg:
		return handleNextPhase(model)
	case *common.PrevPhaseMsg:
		return handlePrevPhase(model)
	case *common.ShowOptionsMsg:
		return handleShowOptions(model)
	case *common.ShowAboutMsg:
		return handleShowAbout(model)
	case *common.ShowMainScreenMsg:
		return handleShowMainScreen(model)
	case *common.RestoreMainUIMsg:
		return model, noCommand
	case *common.TickMsg:
		return handleTick(model)
	case *common.KeyPressMsg:
		return handleKeyPress(msg, model)
	// Handle option update messages
	case *common.SetRulesetMsg:
		return handleSetRuleset(msg, model)
	case *common.SetPlayerCountMsg:
		return handleSetPlayerCount(msg, model)
	case *common.SetPlayerNameMsg:
		return handleSetPlayerName(msg, model)
	case *common.SetColorPaletteMsg:
		return handleSetColorPalette(msg, model)
	case *common.SetTimeFormatMsg:
		return handleSetTimeFormat(msg, model)
	case *common.SetOneTurnForAllPlayersMsg:
		return handleSetOneTurnForAllPlayers(msg, model)
	case *common.SetEnableLogMsg:
		newModel := model
		newModel.Options.LoggingEnabled = msg.Value
		// Persist options to disk
		_ = options.SaveOptions(newModel.Options, "", true)
		return newModel, noCommand
	default:
		return model, noCommand
	}
}

// handleStartGame handles the startGameMsg
func handleStartGame(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model

	// Toggle between start and pause
	if model.GameStatus == gamePaused {
		// Resume the game
		newModel.GameStatus = gameInProgress

		// Log action for active player(s)
		for i, player := range model.Players {
			if player.IsTurn {
				logging.AddLogEntry(newModel.Players[i], &newModel, "Game resumed")
			}
		}
	} else if model.GameStatus == gameInProgress {
		// Pause the game
		newModel.GameStatus = gamePaused

		// Log action for active player(s)
		for i, player := range model.Players {
			if player.IsTurn {
				logging.AddLogEntry(newModel.Players[i], &newModel, "Game paused")
			}
		}
	} else {
		// Start the game if not already started
		newModel.GameStatus = gameInProgress
		newModel.GameStarted = true

		// Log action for active player(s)
		for i, player := range model.Players {
			if player.IsTurn {
				logging.AddLogEntry(newModel.Players[i], &newModel, "Game started")
			}
		}
	}

	return newModel, noCommand
}

// handleEndGame handles the endGameMsg
func handleEndGame(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model

	// Only handle if the game was started
	if model.GameStarted {
		// Reset game state
		newModel.GameStatus = gameNotStarted
		newModel.GameStarted = false
		newModel.TotalGameTime = 0

		// Log action for players
		for i := range model.Players {
			// Reset player state
			newModel.Players[i].TimeElapsed = 0
			newModel.Players[i].TurnCount = 0
			newModel.Players[i].CurrentPhase = 0

			// Clear the action log
			newModel.Players[i].ActionLog = []common.LogEntry{}

			// Keep turn state of player 1
			if i == 0 {
				newModel.Players[i].IsTurn = true
				logging.AddLogEntry(newModel.Players[i], &newModel, "Game ended - reset to initial state")
			} else {
				newModel.Players[i].IsTurn = false
				logging.AddLogEntry(newModel.Players[i], &newModel, "Game ended")
			}
		}
	}

	return newModel, noCommand
}

// handleEndGameConfirm handles the endGameConfirmMsg
func handleEndGameConfirm(msg *common.EndGameConfirmMsg, model common.Model) (common.Model, Command) {
	// CreateAboutPanel a command that will restore the main UI after handling the confirmation
	restoreUICmd := func() common.Message {
		return &common.ShowMainScreenMsg{}
	}

	// If user confirmed ending the game, proceed with the game ending logic
	if msg.Confirmed {
		// Get the updated model after ending the game
		newModel, _ := handleEndGame(model)
		return newModel, restoreUICmd
	}

	// If user canceled, just restore the UI
	return model, restoreUICmd
}

// handleShowEndGameConfirm handles the showEndGameConfirmMsg
func handleShowEndGameConfirm(model common.Model) (common.Model, Command) {
	// Return the model unchanged and a command that will show the confirmation dialog
	return model, func() common.Message {
		// This will be handled by the main.go to show the dialog
		return &common.ShowModalMsg{Type: "EndGameConfirm"}
	}
}

// handleSwitchTurns handles the switchTurnsMsg
func handleSwitchTurns(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model
	newPlayers := make([]*common.Player, len(model.Players))

	// Log for currently active players that their turn is ending
	for i, player := range model.Players {
		// CreateAboutPanel a copy of each player to avoid modifying the original
		newPlayer := *player
		newPlayers[i] = &newPlayer

		if player.IsTurn {
			logging.AddLogEntry(newPlayers[i], &newModel, "Turn %d ended", player.TurnCount)
		}

		// Switch turns
		newPlayers[i].IsTurn = !player.IsTurn

		if newPlayers[i].IsTurn {
			// Increment turn count when a player's turn begins
			newPlayers[i].TurnCount++
			newPlayers[i].CurrentPhase = 0
			// Log for newly active players that their turn is starting
			logging.AddLogEntry(newPlayers[i], &newModel, "Turn %d started", newPlayers[i].TurnCount)
			if len(model.Phases) > 0 {
				logging.AddLogEntry(newPlayers[i], &newModel, "Turn %d - Entered phase: %s", newPlayers[i].TurnCount, model.Phases[0])
			}
		}
	}

	// Update the model with the new players
	newModel.Players = newPlayers

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		newModel.CurrentScreen = "main"
	}

	return newModel, noCommand
}

// handleNextPhase handles the nextPhaseMsg
func handleNextPhase(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model
	newPlayers := make([]*common.Player, len(model.Players))

	// Move forward in the phase
	for i, player := range model.Players {
		// CreateAboutPanel a copy of each player
		newPlayer := *player
		newPlayers[i] = &newPlayer

		if player.IsTurn && player.CurrentPhase < len(model.Phases)-1 {
			newPlayers[i].CurrentPhase = player.CurrentPhase + 1

			// Log the phase change
			logging.AddLogEntry(newPlayers[i], &newModel, "Started phase: %s",
				model.Phases[newPlayers[i].CurrentPhase])
		}
	}

	// Update the model with the new players
	newModel.Players = newPlayers

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		newModel.CurrentScreen = "main"
	}

	return newModel, noCommand
}

// handlePrevPhase handles the prevPhaseMsg
func handlePrevPhase(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model
	newPlayers := make([]*common.Player, len(model.Players))

	// Move backward in the phase
	for i, player := range model.Players {
		// CreateAboutPanel a copy of each player
		newPlayer := *player
		newPlayers[i] = &newPlayer

		if player.IsTurn && player.CurrentPhase > 0 {
			newPlayers[i].CurrentPhase = player.CurrentPhase - 1

			// Log the phase change
			logging.AddLogEntry(newPlayers[i], &newModel, "Started phase: %s",
				model.Phases[newPlayers[i].CurrentPhase])
		}
	}

	// Update the model with the new players
	newModel.Players = newPlayers

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		newModel.CurrentScreen = "main"
	}

	return newModel, noCommand
}

// handleShowOptions handles the showOptionsMsg
func handleShowOptions(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model

	// Toggle between main screen and options screen
	if model.CurrentScreen == "options" {
		newModel.CurrentScreen = "main"
	} else {
		newModel.CurrentScreen = "options"
	}

	return newModel, noCommand
}

// handleShowAbout handles the showAboutMsg
func handleShowAbout(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model

	// Toggle between main screen and about screen
	if model.CurrentScreen == "about" {
		newModel.CurrentScreen = "main"
	} else {
		newModel.CurrentScreen = "about"
	}

	return newModel, noCommand
}

// handleShowMainScreen handles the showMainScreenMsg
func handleShowMainScreen(model common.Model) (common.Model, Command) {
	// CreateAboutPanel a copy of the model to avoid modifying the original
	newModel := model

	// Return to the main screen
	newModel.CurrentScreen = "main"

	// Return a command that will restore the main UI from any modal
	return newModel, func() common.Message {
		return &common.RestoreMainUIMsg{}
	}
}

// handleTick handles the TickMsg
func handleTick(model common.Model) (common.Model, Command) {
	// Only increment time if the game is in progress (not paused)
	if model.GameStarted && model.GameStatus == gameInProgress {
		// CreateAboutPanel a copy of the model to avoid modifying the original
		newModel := model
		newPlayers := make([]*common.Player, len(model.Players))

		// Increment total game time
		newModel.TotalGameTime += 1 * time.Second

		for i, player := range model.Players {
			// CreateAboutPanel a copy of each player
			newPlayer := *player
			newPlayers[i] = &newPlayer

			if player.IsTurn {
				newPlayers[i].TimeElapsed += 1 * time.Second
			}
		}

		// Update the model with the new players
		newModel.Players = newPlayers
		return newModel, noCommand
	}

	// Don't return a TickCommand here as we already have a ticker in main.go
	return model, noCommand
}

// handleKeyPress handles the keyPressMsg
func handleKeyPress(msg *common.KeyPressMsg, model common.Model) (common.Model, Command) {
	switch msg.Key {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		// Quit the application
		// This will be handled in the main function
		return model, noCommand
	case tcell.KeyRune:
		switch string(msg.Rune) {
		case "o", "O":
			return handleShowOptions(model)
		case "a", "A":
			// Toggle about screen
			return handleShowAbout(model)
		case "s", "S":
			// Start/pause/resume game
			return handleStartGame(model)
		case "e", "E":
			// End game (only if game has started)
			if model.GameStarted {
				// Show confirmation dialog instead of directly ending the game
				return model, func() common.Message {
					// Return a command that will show the confirmation dialog
					return &common.ShowEndGameConfirmMsg{}
				}
			}
		case "p", "P":
			// Next phase
			return handleNextPhase(model)
		case "b", "B":
			// Previous phase
			return handlePrevPhase(model)
		case "q", "Q":
			// Quit the application
			// This will be handled in the main function
			return model, noCommand
		case " ":
			// Switch turns
			return handleSwitchTurns(model)
		}
	default:
		// Handle other keys if needed
	}

	return model, noCommand
}

// SetupInputCapture sets up the input capture for the tview application
func SetupInputCapture(app *tview.Application, msgChan chan<- common.Message) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Send a KeyPressMsg to the message channel
		msgChan <- &common.KeyPressMsg{Key: event.Key(), Rune: event.Rune()}

		// Handle specific keys and prevent them from propagating
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'o', 'O', 'a', 'A', 's', 'S', 'e', 'E', 'p', 'P', 'b', 'B', 'q', 'Q', ' ':
				return nil
			}
		default:
			// Handle other keys if needed
		}
		return event
	})
}

// Option update handlers
// handleSetRuleset handles changes to the selected ruleset
func handleSetRuleset(msg *common.SetRulesetMsg, model common.Model) (common.Model, Command) {
	newModel := model
	newModel.Options.Default = msg.Index
	newModel.Phases = model.Options.Rules[msg.Index].Phases
	return newModel, noCommand
}

// handleSetPlayerCount handles changes to the player count
func handleSetPlayerCount(msg *common.SetPlayerCountMsg, model common.Model) (common.Model, Command) {
	if msg.Count <= 0 {
		return model, noCommand
	}

	newModel := model
	newModel.Options.PlayerCount = msg.Count

	// Ensure player names slice has the right length
	if len(newModel.Options.PlayerNames) < msg.Count {
		newModel.Options.PlayerNames = append(
			append([]string{}, newModel.Options.PlayerNames...),
			make([]string, msg.Count-len(newModel.Options.PlayerNames))...)
	}
	return newModel, noCommand
}

// handleSetPlayerName handles changes to a player's name
func handleSetPlayerName(msg *common.SetPlayerNameMsg, model common.Model) (common.Model, Command) {
	if msg.Index < 0 || msg.Index >= len(model.Options.PlayerNames) {
		return model, noCommand
	}

	newModel := model
	newNames := append([]string{}, newModel.Options.PlayerNames...)
	newNames[msg.Index] = msg.Name
	newModel.Options.PlayerNames = newNames
	return newModel, noCommand
}

// handleSetColorPalette handles changes to the color palette
func handleSetColorPalette(msg *common.SetColorPaletteMsg, model common.Model) (common.Model, Command) {
	newModel := model
	newModel.Options.ColorPalette = msg.Name
	newModel.CurrentColorPalette = palette.ColorPaletteByName(msg.Name)
	return newModel, noCommand
}

// handleSetTimeFormat handles changes to the time format
func handleSetTimeFormat(msg *common.SetTimeFormatMsg, model common.Model) (common.Model, Command) {
	newModel := model
	newModel.Options.TimeFormat = msg.Format
	return newModel, noCommand
}

// handleSetOneTurnForAllPlayers handles changes to the "One Turn For All Players" option
func handleSetOneTurnForAllPlayers(msg *common.SetOneTurnForAllPlayersMsg, model common.Model) (common.Model, Command) {
	newModel := model
	newRules := append([]rules.Rules{}, newModel.Options.Rules...)
	newRule := newRules[newModel.Options.Default]
	newRule.OneTurnForAllPlayers = msg.Value
	newRules[newModel.Options.Default] = newRule
	newModel.Options.Rules = newRules
	return newModel, noCommand
}

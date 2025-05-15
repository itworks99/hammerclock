package app

import (
	"time"

	"hammerclock/components/hammerclock/Palette"
	"hammerclock/components/hammerclock/Rules"
	"hammerclock/components/hammerclock/fileio"
	logpanel "hammerclock/internal/app/LogPanel"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Message represents a message that can be sent to the Update function
type Message interface {
	// This is a marker interface
}

// StartGameMsg is sent when the user wants to start/pause/resume the game
type StartGameMsg struct{}

// SwitchTurnsMsg is sent when the user wants to switch turns
type SwitchTurnsMsg struct{}

// NextPhaseMsg is sent when the user wants to move to the next phase
type NextPhaseMsg struct{}

// PrevPhaseMsg is sent when the user wants to move to the previous phase
type PrevPhaseMsg struct{}

// ShowOptionsMsg is sent when the user wants to show the options screen
type ShowOptionsMsg struct{}

// ShowAboutMsg is sent when the user wants to show the about screen
type ShowAboutMsg struct{}

// ShowMainScreenMsg is sent when the user wants to return to the main screen
type ShowMainScreenMsg struct{}

// TickMsg is sent every second to update the clock and player times
type TickMsg struct{}

// KeyPressMsg is sent when a key is pressed
type KeyPressMsg struct {
	Key  tcell.Key
	Rune rune
}

// EndGameMsg is sent when the user wants to end the current game
type EndGameMsg struct{}

// EndGameConfirmMsg is sent when the user confirms or cancels ending the game
type EndGameConfirmMsg struct {
	Confirmed bool
}

// ShowEndGameConfirmMsg is sent to show the end game confirmation dialog
type ShowEndGameConfirmMsg struct{}

// ShowModalMsg is sent to show a modal dialog
type ShowModalMsg struct {
	Type string
}

// RestoreMainUIMsg is sent to restore the main UI after a modal dialog
type RestoreMainUIMsg struct{}

// SetRulesetMsg is sent when the user selects a different ruleset
type SetRulesetMsg struct {
	Index int
}

// SetPlayerCountMsg is sent when the user changes the player count
type SetPlayerCountMsg struct {
	Count int
}

// SetPlayerNameMsg is sent when a player name is changed
type SetPlayerNameMsg struct {
	Index int
	Name  string
}

// SetColorPaletteMsg is sent when the color palette is changed
type SetColorPaletteMsg struct {
	Name string
}

// SetTimeFormatMsg is sent when the time format is changed
type SetTimeFormatMsg struct {
	Format string
}

// SetOneTurnForAllPlayersMsg is sent when the "One Turn For All Players" option is toggled
type SetOneTurnForAllPlayersMsg struct {
	Value bool
}

// SetEnableCSVLogMsg is sent when the user toggles CSV logging
type SetEnableCSVLogMsg struct {
	Value bool
}

// Command represents a command that can be executed after an update
type Command func() Message

// NoCommand is a command that does nothing
func NoCommand() Message {
	return nil
}

// Update processes a message and returns an updated model and a command to execute
func Update(msg Message, model Model) (Model, Command) {
	switch msg := msg.(type) {
	case *StartGameMsg:
		return handleStartGame(model)
	case *EndGameMsg:
		return handleEndGame(model)
	case *EndGameConfirmMsg:
		return handleEndGameConfirm(msg, model)
	case *ShowEndGameConfirmMsg:
		return handleShowEndGameConfirm(model)
	case *SwitchTurnsMsg:
		return handleSwitchTurns(model)
	case *NextPhaseMsg:
		return handleNextPhase(model)
	case *PrevPhaseMsg:
		return handlePrevPhase(model)
	case *ShowOptionsMsg:
		return handleShowOptions(model)
	case *ShowAboutMsg:
		return handleShowAbout(model)
	case *ShowMainScreenMsg:
		return handleShowMainScreen(model)
	case *RestoreMainUIMsg:
		return model, NoCommand
	case *TickMsg:
		return handleTick(model)
	case *KeyPressMsg:
		return handleKeyPress(msg, model)
	// Handle option update messages
	case *SetRulesetMsg:
		return handleSetRuleset(msg, model)
	case *SetPlayerCountMsg:
		return handleSetPlayerCount(msg, model)
	case *SetPlayerNameMsg:
		return handleSetPlayerName(msg, model)
	case *SetColorPaletteMsg:
		return handleSetColorPalette(msg, model)
	case *SetTimeFormatMsg:
		return handleSetTimeFormat(msg, model)
	case *SetOneTurnForAllPlayersMsg:
		return handleSetOneTurnForAllPlayers(msg, model)
	case *SetEnableCSVLogMsg:
		newModel := model
		newModel.Options.EnableCSVLog = msg.Value
		// Persist options to disk
		_ = fileio.SaveOptions(newModel.Options, "", true)
		return newModel, NoCommand
	default:
		return model, NoCommand
	}
}

// handleStartGame handles the StartGameMsg
func handleStartGame(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model

	// Toggle between start and pause
	if model.GameStatus == GamePaused {
		// Resume the game
		newModel.GameStatus = GameInProgress

		// Log action for active player(s)
		for i, player := range model.Players {
			if player.IsTurn {
				AddLogEntry(newModel.Players[i], &newModel, "Game resumed")
			}
		}
	} else if model.GameStatus == GameInProgress {
		// Pause the game
		newModel.GameStatus = GamePaused

		// Log action for active player(s)
		for i, player := range model.Players {
			if player.IsTurn {
				AddLogEntry(newModel.Players[i], &newModel, "Game paused")
			}
		}
	} else {
		// Start the game if not already started
		newModel.GameStatus = GameInProgress
		newModel.GameStarted = true

		// Log action for active player(s)
		for i, player := range model.Players {
			if player.IsTurn {
				AddLogEntry(newModel.Players[i], &newModel, "Game started")
			}
		}
	}

	return newModel, NoCommand
}

// handleEndGame handles the EndGameMsg
func handleEndGame(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model

	// Only handle if the game was started
	if model.GameStarted {
		// Reset game state
		newModel.GameStatus = GameNotStarted
		newModel.GameStarted = false
		newModel.TotalGameTime = 0

		// Log action for players
		for i, _ := range model.Players {
			// Reset player state
			newModel.Players[i].TimeElapsed = 0
			newModel.Players[i].TurnCount = 0
			newModel.Players[i].CurrentPhase = 0
			
			// Clear the action log
			newModel.Players[i].ActionLog = []logpanel.LogEntry{}
			
			// Keep turn state of player 1
			if i == 0 {
				newModel.Players[i].IsTurn = true
				AddLogEntry(newModel.Players[i], &newModel, "Game ended - reset to initial state")
			} else {
				newModel.Players[i].IsTurn = false
				AddLogEntry(newModel.Players[i], &newModel, "Game ended")
			}
		}
	}

	return newModel, NoCommand
}

// handleEndGameConfirm handles the EndGameConfirmMsg
func handleEndGameConfirm(msg *EndGameConfirmMsg, model Model) (Model, Command) {
	// Create a command that will restore the main UI after handling the confirmation
	restoreUICmd := func() Message {
		return &ShowMainScreenMsg{}
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

// handleShowEndGameConfirm handles the ShowEndGameConfirmMsg
func handleShowEndGameConfirm(model Model) (Model, Command) {
	// Return the model unchanged and a command that will show the confirmation dialog
	return model, func() Message {
		// This will be handled by the main.go to show the dialog
		return &ShowModalMsg{Type: "EndGameConfirm"}
	}
}

// handleSwitchTurns handles the SwitchTurnsMsg
func handleSwitchTurns(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model
	newPlayers := make([]*Player, len(model.Players))

	// Log for currently active players that their turn is ending
	for i, player := range model.Players {
		// Create a copy of each player to avoid modifying the original
		newPlayer := *player
		newPlayers[i] = &newPlayer

		if player.IsTurn {
			AddLogEntry(newPlayers[i], &newModel, "Turn %d ended", player.TurnCount)
		}

		// Switch turns
		newPlayers[i].IsTurn = !player.IsTurn

		if newPlayers[i].IsTurn {
			// Increment turn count when a player's turn begins
			newPlayers[i].TurnCount++
			newPlayers[i].CurrentPhase = 0
			// Log for newly active players that their turn is starting
			AddLogEntry(newPlayers[i], &newModel, "Turn %d started", newPlayers[i].TurnCount)
			if len(model.Phases) > 0 {
				AddLogEntry(newPlayers[i], &newModel, "Turn %d - Entered phase: %s", newPlayers[i].TurnCount, model.Phases[0])
			}
		}
	}

	// Update the model with the new players
	newModel.Players = newPlayers

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		newModel.CurrentScreen = "main"
	}

	return newModel, NoCommand
}

// handleNextPhase handles the NextPhaseMsg
func handleNextPhase(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model
	newPlayers := make([]*Player, len(model.Players))

	// Move forward in the phase
	for i, player := range model.Players {
		// Create a copy of each player
		newPlayer := *player
		newPlayers[i] = &newPlayer

		if player.IsTurn && player.CurrentPhase < len(model.Phases)-1 {
			newPlayers[i].CurrentPhase = player.CurrentPhase + 1

			// Log the phase change
			AddLogEntry(newPlayers[i], &newModel, "Started phase: %s",
				model.Phases[newPlayers[i].CurrentPhase])
		}
	}

	// Update the model with the new players
	newModel.Players = newPlayers

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		newModel.CurrentScreen = "main"
	}

	return newModel, NoCommand
}

// handlePrevPhase handles the PrevPhaseMsg
func handlePrevPhase(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model
	newPlayers := make([]*Player, len(model.Players))

	// Move backward in the phase
	for i, player := range model.Players {
		// Create a copy of each player
		newPlayer := *player
		newPlayers[i] = &newPlayer

		if player.IsTurn && player.CurrentPhase > 0 {
			newPlayers[i].CurrentPhase = player.CurrentPhase - 1

			// Log the phase change
			AddLogEntry(newPlayers[i], &newModel, "Started phase: %s",
				model.Phases[newPlayers[i].CurrentPhase])
		}
	}

	// Update the model with the new players
	newModel.Players = newPlayers

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		newModel.CurrentScreen = "main"
	}

	return newModel, NoCommand
}

// handleShowOptions handles the ShowOptionsMsg
func handleShowOptions(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model

	// Toggle between main screen and options screen
	if model.CurrentScreen == "options" {
		newModel.CurrentScreen = "main"
	} else {
		newModel.CurrentScreen = "options"
	}

	return newModel, NoCommand
}

// handleShowAbout handles the ShowAboutMsg
func handleShowAbout(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model

	// Toggle between main screen and about screen
	if model.CurrentScreen == "about" {
		newModel.CurrentScreen = "main"
	} else {
		newModel.CurrentScreen = "about"
	}

	return newModel, NoCommand
}

// handleShowMainScreen handles the ShowMainScreenMsg
func handleShowMainScreen(model Model) (Model, Command) {
	// Create a copy of the model to avoid modifying the original
	newModel := model

	// Return to the main screen
	newModel.CurrentScreen = "main"

	// Return a command that will restore the main UI from any modal
	return newModel, func() Message {
		return &RestoreMainUIMsg{}
	}
}

// handleTick handles the TickMsg
func handleTick(model Model) (Model, Command) {
	// Only increment time if the game is in progress (not paused)
	if model.GameStarted && model.GameStatus == GameInProgress {
		// Create a copy of the model to avoid modifying the original
		newModel := model
		newPlayers := make([]*Player, len(model.Players))

		// Increment total game time
		newModel.TotalGameTime += 1 * time.Second

		for i, player := range model.Players {
			// Create a copy of each player
			newPlayer := *player
			newPlayers[i] = &newPlayer

			if player.IsTurn {
				newPlayers[i].TimeElapsed += 1 * time.Second
			}
		}

		// Update the model with the new players
		newModel.Players = newPlayers
		return newModel, NoCommand
	}

	// Don't return a TickCommand here as we already have a ticker in main.go
	return model, NoCommand
}

// handleKeyPress handles the KeyPressMsg
func handleKeyPress(msg *KeyPressMsg, model Model) (Model, Command) {
	switch msg.Key {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		// Quit the application
		// This will be handled in the main function
		return model, NoCommand
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
				return model, func() Message {
					// Return a command that will show the confirmation dialog
					return &ShowEndGameConfirmMsg{}
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
			return model, NoCommand
		case " ":
			// Switch turns
			return handleSwitchTurns(model)
		}
	default:
		// Handle other keys if needed
	}

	return model, NoCommand
}

// SetupInputCapture sets up the input capture for the tview application
func SetupInputCapture(app *tview.Application, msgChan chan<- Message) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Send a KeyPressMsg to the message channel
		msgChan <- &KeyPressMsg{Key: event.Key(), Rune: event.Rune()}

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
func handleSetRuleset(msg *SetRulesetMsg, model Model) (Model, Command) {
	newModel := model
	newModel.Options.Default = msg.Index
	newModel.Phases = model.Options.Rules[msg.Index].Phases
	return newModel, NoCommand
}

// handleSetPlayerCount handles changes to the player count
func handleSetPlayerCount(msg *SetPlayerCountMsg, model Model) (Model, Command) {
	if msg.Count <= 0 {
		return model, NoCommand
	}

	newModel := model
	newModel.Options.PlayerCount = msg.Count

	// Ensure player names slice has the right length
	if len(newModel.Options.PlayerNames) < msg.Count {
		newModel.Options.PlayerNames = append(
			append([]string{}, newModel.Options.PlayerNames...),
			make([]string, msg.Count-len(newModel.Options.PlayerNames))...)
	}
	return newModel, NoCommand
}

// handleSetPlayerName handles changes to a player's name
func handleSetPlayerName(msg *SetPlayerNameMsg, model Model) (Model, Command) {
	if msg.Index < 0 || msg.Index >= len(model.Options.PlayerNames) {
		return model, NoCommand
	}

	newModel := model
	newNames := append([]string{}, newModel.Options.PlayerNames...)
	newNames[msg.Index] = msg.Name
	newModel.Options.PlayerNames = newNames
	return newModel, NoCommand
}

// handleSetColorPalette handles changes to the color palette
func handleSetColorPalette(msg *SetColorPaletteMsg, model Model) (Model, Command) {
	newModel := model
	newModel.Options.ColorPalette = msg.Name
	newModel.CurrentColorPalette = Palette.ColorPaletteByName(msg.Name)
	return newModel, NoCommand
}

// handleSetTimeFormat handles changes to the time format
func handleSetTimeFormat(msg *SetTimeFormatMsg, model Model) (Model, Command) {
	newModel := model
	newModel.Options.TimeFormat = msg.Format
	return newModel, NoCommand
}

// handleSetOneTurnForAllPlayers handles changes to the "One Turn For All Players" option
func handleSetOneTurnForAllPlayers(msg *SetOneTurnForAllPlayersMsg, model Model) (Model, Command) {
	newModel := model
	newRules := append([]Rules.Rules{}, newModel.Options.Rules...)
	newRule := newRules[newModel.Options.Default]
	newRule.OneTurnForAllPlayers = msg.Value
	newRules[newModel.Options.Default] = newRule
	newModel.Options.Rules = newRules
	return newModel, NoCommand
}

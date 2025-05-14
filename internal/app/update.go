package app

import (
	"time"

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
	case *TickMsg:
		return handleTick(model)
	case *KeyPressMsg:
		return handleKeyPress(msg, model)
	default:
		return model, NoCommand
	}
}

// handleStartGame handles the StartGameMsg
func handleStartGame(model Model) (Model, Command) {
	// Toggle between start and pause
	if model.GameStatus == GamePaused {
		// Resume the game
		model.GameStatus = GameInProgress
		
		// Log action for active player(s)
		for _, player := range model.Players {
			if player.IsTurn {
				AddLogEntry(player, "Game resumed (Turn %d)", player.TurnCount)
			}
		}
	} else if model.GameStatus == GameInProgress {
		// Pause the game
		model.GameStatus = GamePaused
		
		// Log action for active player(s)
		for _, player := range model.Players {
			if player.IsTurn {
				AddLogEntry(player, "Game paused (Turn %d)", player.TurnCount)
			}
		}
	} else {
		// Start the game if not already started
		model.GameStatus = GameInProgress
		model.GameStarted = true
		
		// Log action for active player(s)
		for _, player := range model.Players {
			if player.IsTurn {
				AddLogEntry(player, "Game started (Turn %d)", player.TurnCount)
			}
		}
	}

	return model, NoCommand
}

// handleSwitchTurns handles the SwitchTurnsMsg
func handleSwitchTurns(model Model) (Model, Command) {
	// Log for currently active players that their turn is ending
	for _, player := range model.Players {
		if player.IsTurn {
			AddLogEntry(player, "Turn %d ended", player.TurnCount)
		}
	}

	// Switch turns
	for _, player := range model.Players {
		player.IsTurn = !player.IsTurn
		if player.IsTurn {
			// Increment turn count when a player's turn begins
			player.TurnCount++
			player.CurrentPhase = 0
			// Log for newly active players that their turn is starting
			AddLogEntry(player, "Turn %d started", player.TurnCount)
			if len(model.Phases) > 0 {
				AddLogEntry(player, "Turn %d - Entered phase: %s", player.TurnCount, model.Phases[0])
			}
		}
	}

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		model.CurrentScreen = "main"
	}

	return model, NoCommand
}

// handleNextPhase handles the NextPhaseMsg
func handleNextPhase(model Model) (Model, Command) {
	// Move forward in the phase
	for _, player := range model.Players {
		if player.IsTurn && player.CurrentPhase < len(model.Phases)-1 {
			oldPhase := player.CurrentPhase
			player.CurrentPhase = player.CurrentPhase + 1
			
			// Log the phase change
			AddLogEntry(player, "Turn %d - Moved from phase: %s to phase: %s", 
				player.TurnCount,
				model.Phases[oldPhase], 
				model.Phases[player.CurrentPhase])
		}
	}

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		model.CurrentScreen = "main"
	}

	return model, NoCommand
}

// handlePrevPhase handles the PrevPhaseMsg
func handlePrevPhase(model Model) (Model, Command) {
	// Move backward in the phase
	for _, player := range model.Players {
		if player.IsTurn && player.CurrentPhase > 0 {
			oldPhase := player.CurrentPhase
			player.CurrentPhase = player.CurrentPhase - 1
			
			// Log the phase change
			AddLogEntry(player, "Turn %d - Moved from phase: %s to phase: %s", 
				player.TurnCount,
				model.Phases[oldPhase], 
				model.Phases[player.CurrentPhase])
		}
	}

	// If we're not on the main screen, this is a good time to return to it
	if model.CurrentScreen != "main" {
		model.CurrentScreen = "main"
	}

	return model, NoCommand
}

// handleShowOptions handles the ShowOptionsMsg
func handleShowOptions(model Model) (Model, Command) {
	// Toggle between main screen and options screen
	if model.CurrentScreen == "options" {
		model.CurrentScreen = "main"
	} else {
		model.CurrentScreen = "options"
	}

	return model, NoCommand
}

// handleShowAbout handles the ShowAboutMsg
func handleShowAbout(model Model) (Model, Command) {
	// Toggle between main screen and about screen
	if model.CurrentScreen == "about" {
		model.CurrentScreen = "main"
	} else {
		model.CurrentScreen = "about"
	}

	return model, NoCommand
}

// handleShowMainScreen handles the ShowMainScreenMsg
func handleShowMainScreen(model Model) (Model, Command) {
	// Return to the main screen
	model.CurrentScreen = "main"

	return model, NoCommand
}

// handleTick handles the TickMsg
func handleTick(model Model) (Model, Command) {
	// Only increment time if the game is in progress (not paused)
	if model.GameStarted && model.GameStatus == GameInProgress {
		for _, player := range model.Players {
			if player.IsTurn {
				oldTimeElapsed := player.TimeElapsed
				player.TimeElapsed += 1 * time.Second
				
				// Log time milestone every minute
				if oldTimeElapsed.Minutes() < player.TimeElapsed.Minutes() {
					minutes := int(player.TimeElapsed.Minutes())
					if minutes > 0 && minutes % 1 == 0 {
						AddLogEntry(player, "Turn %d - Time elapsed: %d minute(s)", player.TurnCount, minutes)
					}
				}
			}
		}
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
			case 'o', 'O', 'a', 'A', 's', 'S', 'p', 'P', 'b', 'B', 'q', 'Q', ' ':
				return nil
			}
		default:
			// Handle other keys if needed
		}
		return event
	})
}

package hammerclock

import (
	"fmt"
	"time"

	"hammerclock/internal/hammerclock/logging"
	"hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/palette"
)

// Model represents the entire application state
type Model struct {
	// Game state
	Players             []*Player
	Phases              []string
	GameStatus          gameStatus
	CurrentScreen       string // Can be "main", "options", or "about"
	GameStarted         bool
	Options             options.Options
	CurrentColorPalette palette.ColorPalette
	TotalGameTime       time.Duration // Total elapsed time for the entire game
}

// Player represents a player in the game
type Player struct {
	Name         string
	TimeElapsed  time.Duration // Time elapsed for the player
	IsTurn       bool          // Indicates if it's this player's turn
	CurrentPhase int           // Current phase of the game for this player
	TurnCount    int           // Counter to track number of turns completed
	ArmyList     []unit
	ActionLog    []logging.LogEntry // Log of player actions during the game
}

// unit represents a unit in a player's army
type unit struct {
	Name   string
	Points int
}

// gameStatus represents the current state of the game
type gameStatus string

const (
	gameNotStarted gameStatus = "Game Not Started"
	gameInProgress gameStatus = "Game In Progress"
	gamePaused     gameStatus = "Game Paused"
)

// NewModel creates a new model with default values
func NewModel() Model {
	// Initialize with default options
	opts := options.DefaultOptions

	// CreateAboutPanel players
	players := make([]*Player, opts.PlayerCount)
	model := Model{
		Players:             players,
		Phases:              opts.Rules[opts.Default].Phases,
		GameStatus:          gameNotStarted,
		CurrentScreen:       "main",
		GameStarted:         false,
		Options:             opts,
		CurrentColorPalette: palette.K9sPalette,
		TotalGameTime:       0,
	}

	for i := 0; i < opts.PlayerCount; i++ {
		playerName := fmt.Sprintf("Player %d", i+1)
		if i < len(opts.PlayerNames) {
			playerName = opts.PlayerNames[i]
		}
		players[i] = &Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
			ActionLog:    []logging.LogEntry{}, // Initialize empty action log
		}

		// Add initial log entry
		if i == 0 {
			addLogEntry(players[i], &model, "Initialized - active player")
		} else {
			addLogEntry(players[i], &model, "Initialized")
		}
	}

	return model
}

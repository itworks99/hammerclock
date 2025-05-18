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
	GameStatus          GameStatus
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
	ArmyList     []Unit
	ActionLog    []logging.LogEntry // Log of player actions during the game
}

// Unit represents a unit in a player's army
type Unit struct {
	Name   string
	Points int
}

// GameStatus represents the current state of the game
type GameStatus string

const (
	GameNotStarted GameStatus = "Game Not Started"
	GameInProgress GameStatus = "Game In Progress"
	GamePaused     GameStatus = "Game Paused"
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
		GameStatus:          GameNotStarted,
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
			AddLogEntry(players[i], &model, "Initialized - active player")
		} else {
			AddLogEntry(players[i], &model, "Initialized")
		}
	}

	return model
}

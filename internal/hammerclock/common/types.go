package common

import (
	"time"

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
	ArmyList     []unit
	ActionLog    []LogEntry // Log of player actions during the game
}

// unit represents a unit in a player's army
type unit struct {
	Name   string
	Points int
}

// GameStatus represents the current state of the game
type GameStatus string

// LogEntry represents a single log entry with details about an action.
type LogEntry struct {
	DateTime   string
	PlayerName string
	Turn       int
	Phase      string
	Message    string
}

// Message represents a message that can be sent to the Update function
type Message interface {
}

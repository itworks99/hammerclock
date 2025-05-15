package app

import (
	"time"

	"hammerclock/components/hammerclock/Palette"
	Rules2 "hammerclock/components/hammerclock/Rules"
	hammerclockConfig "hammerclock/config"
	"hammerclock/internal/app/LogPanel"
)

// Model represents the entire application state
type Model struct {
	// Game state
	Players             []*Player
	Phases              []string
	GameStatus          GameStatus
	CurrentScreen       string // Can be "main", "options", or "about"
	GameStarted         bool
	Options             Options
	CurrentColorPalette Palette.ColorPalette
}

// Player represents a player in the game
type Player struct {
	Name         string
	TimeElapsed  time.Duration // Time elapsed for the player
	IsTurn       bool          // Indicates if it's this player's turn
	CurrentPhase int           // Current phase of the game for this player
	TurnCount    int           // Counter to track number of turns completed
	ArmyList     []Unit
	ActionLog    []LogPanel.LogEntry // Log of player actions during the game
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

// Options defines the configuration for a game, including player details, phases, and display preferences.
type Options struct {
	Default      int            `json:"default"`
	Rules        []Rules2.Rules `json:"rules"`
	PlayerCount  int            `json:"playerCount"`
	PlayerNames  []string       `json:"playerNames"`
	ColorPalette string         `json:"colorPalette"`
	TimeFormat   string         `json:"timeFormat"`   // AMPM or 24h
	EnableCSVLog bool           `json:"enableCSVLog"` // Enable/disable CSV logging
}

// DefaultPlayerNames Generate default player names
func DefaultPlayerNames() []string {
	var playerNames []string
	for i := range hammerclockConfig.DefaultPlayerCount {
		playerNames = append(playerNames, hammerclockConfig.DefaultPlayerPrefix+" "+string(rune(i+1)))
	}
	return playerNames
}

// DefaultOptions Default options
var DefaultOptions = Options{
	Default:      0,
	Rules:        Rules2.AllRules,
	PlayerCount:  hammerclockConfig.DefaultPlayerCount,
	PlayerNames:  DefaultPlayerNames(),
	ColorPalette: hammerclockConfig.DefaultColorPalette,
	TimeFormat:   "AMPM",
	EnableCSVLog: true, // CSV logging enabled by default
}

// NewModel creates a new model with default values
func NewModel() Model {
	// Initialize with default options
	options := DefaultOptions

	// Create players
	players := make([]*Player, options.PlayerCount)
	model := Model{
		Players:             players,
		Phases:              options.Rules[options.Default].Phases,
		GameStatus:          GameNotStarted,
		CurrentScreen:       "main",
		GameStarted:         false,
		Options:             options,
		CurrentColorPalette: Palette.K9sPalette,
	}

	for i := 0; i < options.PlayerCount; i++ {
		playerName := options.PlayerNames[i]
		players[i] = &Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
			ActionLog:    []LogPanel.LogEntry{}, // Initialize empty action log
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

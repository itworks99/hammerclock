package app

import (
	"github.com/gdamore/tcell/v2"
	"time"
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
	CurrentColorPalette ColorPalette
}

// Player represents a player in the game
type Player struct {
	Name         string
	TimeElapsed  time.Duration
	IsTurn       bool
	CurrentPhase int
	ArmyList     []Unit
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
	Name                 string   `json:"name"`
	Default              bool     `json:"default"`
	PlayerCount          int      `json:"playerCount"`
	PlayerNames          []string `json:"playerNames"`
	Phases               []string `json:"phases"`
	OneTurnForAllPlayers bool     `json:"oneTurnForAllPlayers"`
	ColorPalette         string   `json:"colorPalette"`
	TimeFormat           string   `json:"timeFormat"` // AMPM or 24h
}

// ColorPalette contains all the colors used in the application
type ColorPalette struct {
	Blue     tcell.Color
	Cyan     tcell.Color
	White    tcell.Color
	DimWhite tcell.Color
	Yellow   tcell.Color
	Green    tcell.Color
	Red      tcell.Color
	Black    tcell.Color
}

// K9sPalette K9s color palette
var K9sPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(36, 96, 146),   // Dark blue for backgrounds
	Cyan:     tcell.NewRGBColor(0, 183, 235),   // Cyan for highlights
	White:    tcell.NewRGBColor(255, 255, 255), // White for primary text
	DimWhite: tcell.NewRGBColor(180, 180, 180), // Dimmed white for inactive panels
	Yellow:   tcell.NewRGBColor(253, 185, 19),  // Yellow for warnings
	Green:    tcell.NewRGBColor(0, 200, 83),    // Green for success/active states
	Red:      tcell.NewRGBColor(255, 0, 0),     // Red for errors/critical states
	Black:    tcell.NewRGBColor(0, 0, 0),       // Black for default backgrounds
}

// DraculaPalette Dracula color palette
var DraculaPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(189, 147, 249), // Purple
	Cyan:     tcell.NewRGBColor(139, 233, 253), // Cyan
	White:    tcell.NewRGBColor(248, 248, 242), // Foreground
	DimWhite: tcell.NewRGBColor(174, 174, 169), // Dimmed foreground
	Yellow:   tcell.NewRGBColor(241, 250, 140), // Yellow
	Green:    tcell.NewRGBColor(80, 250, 123),  // Green
	Red:      tcell.NewRGBColor(255, 85, 85),   // Red
	Black:    tcell.NewRGBColor(40, 42, 54),    // Background
}

// MonokaiPalette Monokai color palette
var MonokaiPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(102, 217, 239), // Blue
	Cyan:     tcell.NewRGBColor(102, 217, 239), // Blue (same as Blue for Monokai)
	White:    tcell.NewRGBColor(248, 248, 242), // Foreground
	DimWhite: tcell.NewRGBColor(174, 174, 169), // Dimmed foreground
	Yellow:   tcell.NewRGBColor(230, 219, 116), // Yellow
	Green:    tcell.NewRGBColor(166, 226, 46),  // Green
	Red:      tcell.NewRGBColor(249, 38, 114),  // Red
	Black:    tcell.NewRGBColor(39, 40, 34),    // Background
}

// DefaultOptions Default options
var DefaultOptions = Options{
	Name:                 "W40K 10th Edition",
	Default:              true,
	PlayerCount:          2,
	PlayerNames:          []string{"Player 1", "Player 2"},
	Phases:               []string{"Command Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Fight Phase", "End Phase"},
	OneTurnForAllPlayers: false,
	ColorPalette:         "k9s",
	TimeFormat:           "AMPM",
}

// NewModel creates a new model with default values
func NewModel() Model {
	// Initialize with default options
	options := DefaultOptions

	// Create players
	players := make([]*Player, options.PlayerCount)
	for i := 0; i < options.PlayerCount; i++ {
		playerName := options.PlayerNames[i]
		players[i] = &Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
		}
	}

	return Model{
		Players:             players,
		Phases:              options.Phases,
		GameStatus:          GameNotStarted,
		CurrentScreen:       "main",
		GameStarted:         false,
		Options:             options,
		CurrentColorPalette: K9sPalette,
	}
}

// GetColorPaletteByName returns the color palette for the given name
func GetColorPaletteByName(name string) ColorPalette {
	switch name {
	case "dracula":
		return DraculaPalette
	case "monokai":
		return MonokaiPalette
	default: // "k9s" or any other value defaults to k9s
		return K9sPalette
	}
}

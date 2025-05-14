package app

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"hammerclock/config"
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
	TurnCount    int      // Counter to track number of turns completed
	ArmyList     []Unit
	ActionLog    []string // Log of player actions during the game
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
	Default      int      `json:"default"`
	Rules        []Rules  `json:"rules"`
	PlayerCount  int      `json:"playerCount"`
	PlayerNames  []string `json:"playerNames"`
	ColorPalette string   `json:"colorPalette"`
	TimeFormat   string   `json:"timeFormat"` // AMPM or 24h
}

// Rules defines the rules for a specific game, including the name, phases, and whether players are only taking
// one turn (in that case, phases are being ignored).
type Rules struct {
	Name                 string   `json:"name"`
	Phases               []string `json:"phases"`
	OneTurnForAllPlayers bool     `json:"oneTurnForAllPlayers"`
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

// WarhammerPalette represents the color theme for Warhammer 40K
var WarhammerPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(38, 57, 132),   // Ultramarine Blue
	Cyan:     tcell.NewRGBColor(23, 155, 215),  // Tyranid Blue
	White:    tcell.NewRGBColor(255, 250, 240), // Imperial White
	DimWhite: tcell.NewRGBColor(180, 170, 150), // Bone Color
	Yellow:   tcell.NewRGBColor(245, 180, 26),  // Imperial Gold
	Green:    tcell.NewRGBColor(0, 120, 50),    // Dark Angels Green
	Red:      tcell.NewRGBColor(190, 0, 0),     // Blood Angels Red
	Black:    tcell.NewRGBColor(10, 10, 10),    // Abaddon Black
}

// KillTeamPalette represents the color theme for Kill Team
var KillTeamPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(63, 81, 153),   // Night Lords Blue
	Cyan:     tcell.NewRGBColor(0, 169, 157),   // Tactical Turquoise
	White:    tcell.NewRGBColor(230, 230, 230), // Tactical White
	DimWhite: tcell.NewRGBColor(150, 150, 150), // Urban Gray
	Yellow:   tcell.NewRGBColor(255, 193, 0),   // Warning Yellow
	Green:    tcell.NewRGBColor(76, 99, 25),    // Camo Green
	Red:      tcell.NewRGBColor(200, 40, 40),   // Target Red
	Black:    tcell.NewRGBColor(5, 5, 5),       // Shadow Black
}

// DefaultPlayerNames Generate default player names
func DefaultPlayerNames() []string {
	var playerNames []string
	for i := 0; i < hammerclockConfig.DefaultPlayerCount; i++ {
		playerNames = append(playerNames, hammerclockConfig.DefaultPlayerPrefix+" "+string(rune(i+1)))
	}
	return playerNames
}

// DefaultOptions Default options
var DefaultOptions = Options{
	Default:      0,
	Rules:        allRules,
	PlayerCount:  hammerclockConfig.DefaultPlayerCount,
	PlayerNames:  DefaultPlayerNames(),
	ColorPalette: hammerclockConfig.DefaultColorPalette,
	TimeFormat:   "AMPM",
}

// allRules contains all the rules available in the application
var allRules = []Rules{
	WarhammerRules,
	KillTeamRules,
	NecromundaRules,
	AgeOfSigmarRules,
	WarcryRules,
	BloodBowlRules,
	BunnyKingdomRules,
	ChessRules,
}

// WarhammerRules Warhammer rules
var WarhammerRules = Rules{
	Name:                 "Warhammer 40K (10th Edition)",
	Phases:               []string{"Command Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Fight Phase", "End Phase"},
	OneTurnForAllPlayers: false,
}

// KillTeamRules Kill Team rules
var KillTeamRules = Rules{
	Name:                 "Kill Team (2021)",
	Phases:               []string{"Initiative Phase", "Movement Phase", "Shooting Phase", "Fight Phase", "Morale Phase"},
	OneTurnForAllPlayers: false,
}

// NecromundaRules Necromunda rules
var NecromundaRules = Rules{
	Name:                 "Necromunda (2022 edition)",
	Phases:               []string{"Recovery Phase", "Action Phase", "End Phase"},
	OneTurnForAllPlayers: false,
}

// AgeOfSigmarRules Age of Sigmar rules
var AgeOfSigmarRules = Rules{
	Name:                 "Age of Sigmar (4th Edition)",
	Phases:               []string{"Start of Turn Phase", "Hero Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Combat Phase", "End of Turn Phase"},
	OneTurnForAllPlayers: false,
}

// WarcryRules Warcry rules
var WarcryRules = Rules{
	Name:                 "Warcry (3rd edition)",
	Phases:               []string{"Set Up Phase", "Players' Phase (activating models alternately)", "End Phase"},
	OneTurnForAllPlayers: false,
}

// BloodBowlRules Blood Bowl rules
var BloodBowlRules = Rules{
	Name:                 "Blood Bowl (2020 edition)",
	Phases:               []string{"Pre-Match Phase", "Kick-Off Phase", "Team Turn (both teams alternate)", "End of Turn Phase", "Post-Match Phase"},
	OneTurnForAllPlayers: false,
}

// BunnyKingdomRules Bunny Kingdom rules
var BunnyKingdomRules = Rules{
	Name: "Bunny Kingdom",
	Phases: []string{"Draft Phase (players select cards)",
		"Build Phase (place cards on the board)",
		"Scoring Phase (calculate points based on card placement)"},
	OneTurnForAllPlayers: false,
}

// ChessRules Chess rules
var ChessRules = Rules{
	Name:                 "Chess",
	Phases:               []string{},
	OneTurnForAllPlayers: true,
}

// AddLogEntry adds a log entry to a player's action log
func AddLogEntry(player *Player, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, fmt.Sprintf(format, args...))
	player.ActionLog = append(player.ActionLog, logEntry)
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
			ActionLog:    []string{}, // Initialize empty action log
		}
		
		// Add initial log entry
		if i == 0 {
			AddLogEntry(players[i], "Player initialized as active player")
		} else {
			AddLogEntry(players[i], "Player initialized")
		}
	}

	return Model{
		Players:             players,
		Phases:              options.Rules[options.Default].Phases,
		GameStatus:          GameNotStarted,
		CurrentScreen:       "main",
		GameStarted:         false,
		Options:             options,
		CurrentColorPalette: K9sPalette,
	}
}

// GetColorPalettes returns a list of available color palettes
func GetColorPalettes() []string {
	return []string{
		"k9s",
		"dracula",
		"monokai",
		"warhammer",
		"killteam",
	}
}

// GetColorPaletteByName returns the color palette for the given name
func GetColorPaletteByName(name string) ColorPalette {
	switch name {
	case "dracula":
		return DraculaPalette
	case "monokai":
		return MonokaiPalette
	case "warhammer":
		return WarhammerPalette
	case "killteam":
		return KillTeamPalette
	default: // "k9s" or any other value defaults to k9s
		return K9sPalette
	}
}

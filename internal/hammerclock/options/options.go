package options

import (
	"hammerclock/internal/hammerclock/config"
	"hammerclock/internal/hammerclock/rules"
)

// Options defines the configuration for a game, including player details, phases, and display preferences.
type Options struct {
	Default      int           `json:"default"`
	Rules        []rules.Rules `json:"rules"`
	PlayerCount  int           `json:"playerCount"`
	PlayerNames  []string      `json:"playerNames"`
	ColorPalette string        `json:"colorPalette"`
	TimeFormat   string        `json:"timeFormat"`   // AMPM or 24h
	EnableCSVLog bool          `json:"enableCSVLog"` // Enable/disable CSV logging
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
	Rules:        rules.AllRules,
	PlayerCount:  hammerclockConfig.DefaultPlayerCount,
	PlayerNames:  DefaultPlayerNames(),
	ColorPalette: hammerclockConfig.DefaultColorPalette,
	TimeFormat:   "AMPM",
	EnableCSVLog: true, // CSV logging enabled by default
}

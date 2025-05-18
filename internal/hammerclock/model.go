package hammerclock

import (
	"fmt"

	"hammerclock/internal/hammerclock/common"
	"hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/palette"
)

const (
	gameNotStarted common.GameStatus = "Game Not Started"
	gameInProgress common.GameStatus = "Game In Progress"
	gamePaused     common.GameStatus = "Game Paused"
)

// NewModel creates a new model with default values
func NewModel() common.Model {
	// Initialize with default options
	opts := options.DefaultOptions

	// CreateAboutPanel players
	players := make([]*common.Player, opts.PlayerCount)
	model := common.Model{
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
		players[i] = &common.Player{
			Name:         playerName,
			TimeElapsed:  0,
			IsTurn:       i == 0,
			CurrentPhase: 0,
			ActionLog:    []common.LogEntry{}, // Initialize empty action log
		}
	}

	return model
}

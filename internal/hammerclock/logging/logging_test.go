package logging

import (
	"fmt"
	"testing"

	"hammerclock/internal/hammerclock/common"
	"hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/rules"
)

var testModel = &common.Model{
	Players: []*common.Player{
		{Name: "Player 1", IsTurn: true},
		{Name: "Player 2"},
	},
	Phases: []string{"Setup", "Movement", "Shooting", "Melee", "End"},
	//GameStatus:    gameNotStarted,
	CurrentScreen: "main",
	Options: options.Options{
		TimeFormat: "24h",
		Rules: []rules.Rules{
			{
				Name:                 "Default Rules",
				Phases:               []string{"Setup", "Movement", "Shooting", "Melee", "End"},
				OneTurnForAllPlayers: true,
			},
		},
		Default: 0,
	},
	TotalGameTime: 0,
}

func TestInitialiseSetsUpLoggingCorrectly(t *testing.T) {
	Initialise()
	if logChannel == nil {
		t.Error("Expected logChannel to be initialized")
	}
	if !logInitialized {
		t.Error("Expected logInitialized to be true")
	}
	Cleanup()
}

func TestCleanupClosesLogChannelAndResetsState(t *testing.T) {
	Initialise()
	Cleanup()
	if logInitialized {
		t.Error("Expected logInitialized to be false after cleanup")
	}
	if _, ok := <-logChannel; ok {
		t.Error("Expected logChannel to be closed after cleanup")
	}
}

func TestSendLogEntryDropsEntryWhenChannelIsFull(t *testing.T) {
	Initialise()
	defer Cleanup()

	for i := 0; i < cap(logChannel); i++ {
		sendLogEntry(common.LogEntry{Message: fmt.Sprintf("Log %d", i)})
	}

	// Attempt to send one more log entry
	sendLogEntry(common.LogEntry{Message: "Dropped log"})
	if len(logChannel) != cap(logChannel) {
		t.Errorf("Expected logChannel to remain full, got %d entries", len(logChannel))
	}
}

func TestAddLogEntryAppendsToPlayerActionLog(t *testing.T) {
	player := &common.Player{Name: "Player 1"}
	model := testModel

	AddLogEntry(player, model, "Test message")
	if len(player.ActionLog) != 1 {
		t.Errorf("Expected player.ActionLog to have 1 entry, got %d", len(player.ActionLog))
	}
	if player.ActionLog[0].Message != "Test message" {
		t.Errorf("Expected log message to be 'Test message', got '%s'", player.ActionLog[0].Message)
	}
}

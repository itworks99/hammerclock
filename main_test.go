package main

import (
	"testing"
	"time"
)

// TestReadSettingsFile tests that the settings file can be read correctly
func TestReadSettingsFile(t *testing.T) {
	settings := readSettingsFile()

	// Check that the settings are loaded correctly
	if settings.Name != "W40K 10th Edition" {
		t.Errorf("Expected game name to be 'W40K 10th Edition', got '%s'", settings.Name)
	}

	if settings.PlayerCount != 2 {
		t.Errorf("Expected player count to be 2, got %d", settings.PlayerCount)
	}

	if len(settings.Phases) != 6 {
		t.Errorf("Expected 6 phases, got %d", len(settings.Phases))
	}

	// Check the first and last phase
	if settings.Phases[0] != "Command Phase" {
		t.Errorf("Expected first phase to be 'Command Phase', got '%s'", settings.Phases[0])
	}

	if settings.Phases[5] != "End Phase" {
		t.Errorf("Expected last phase to be 'End Phase', got '%s'", settings.Phases[5])
	}
}

// TestPlayerFunctionality tests the Player struct functionality
func TestPlayerFunctionality(t *testing.T) {
	// Initialize phases for testing
	phases = []string{"Command Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Fight Phase", "End Phase"}

	// Create a new player
	player := &Player{
		name:         "Test Player",
		timeElapsed:  0,
		isTurn:       true,
		currentPhase: 0,
		armyList:     []Unit{},
	}

	// Test initial state
	if player.name != "Test Player" {
		t.Errorf("Expected player name to be 'Test Player', got '%s'", player.name)
	}

	if player.timeElapsed != 0 {
		t.Errorf("Expected time elapsed to be 0, got %v", player.timeElapsed)
	}

	if !player.isTurn {
		t.Errorf("Expected player to be in turn")
	}

	if player.currentPhase != 0 {
		t.Errorf("Expected current phase to be 0, got %d", player.currentPhase)
	}

	// Test phase advancement
	player.currentPhase = (player.currentPhase + 1) % len(phases)
	if player.currentPhase != 1 {
		t.Errorf("Expected current phase to be 1 after advancement, got %d", player.currentPhase)
	}

	// Test time tracking
	player.timeElapsed += 5 * time.Second
	if player.timeElapsed != 5*time.Second {
		t.Errorf("Expected time elapsed to be 5 seconds, got %v", player.timeElapsed)
	}

	// Test loading army list
	loadArmyListFromFile(player)
	if len(player.armyList) != 3 {
		t.Errorf("Expected 3 units in army list, got %d", len(player.armyList))
	}

	if player.armyList[0].Name != "Space Marine" {
		t.Errorf("Expected first unit to be 'Space Marine', got '%s'", player.armyList[0].Name)
	}
}

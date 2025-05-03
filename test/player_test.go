package test

import (
	"testing"
	"time"

	"hammerclock/internal/app"
)

// TestPlayerFunctionality tests the Player struct functionality
func TestPlayerFunctionality(t *testing.T) {
	// Initialize phases for testing
	phases := []string{"Command Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Fight Phase", "End Phase"}

	// Create a new player
	player := &app.Player{
		Name:         "Test Player",
		TimeElapsed:  0,
		IsTurn:       true,
		CurrentPhase: 0,
		ArmyList:     []app.Unit{},
	}

	// Test initial state
	if player.Name != "Test Player" {
		t.Errorf("Expected player name to be 'Test Player', got '%s'", player.Name)
	}

	if player.TimeElapsed != 0 {
		t.Errorf("Expected time elapsed to be 0, got %v", player.TimeElapsed)
	}

	if !player.IsTurn {
		t.Errorf("Expected player to be in turn")
	}

	if player.CurrentPhase != 0 {
		t.Errorf("Expected current phase to be 0, got %d", player.CurrentPhase)
	}

	// Test phase advancement
	player.CurrentPhase = (player.CurrentPhase + 1) % len(phases)
	if player.CurrentPhase != 1 {
		t.Errorf("Expected current phase to be 1 after advancement, got %d", player.CurrentPhase)
	}

	// Test time tracking
	player.TimeElapsed += 5 * time.Second
	if player.TimeElapsed != 5*time.Second {
		t.Errorf("Expected time elapsed to be 5 seconds, got %v", player.TimeElapsed)
	}

	// Add a unit to the army list
	player.ArmyList = append(player.ArmyList, app.Unit{
		Name:   "Space Marine",
		Points: 100,
	})

	// Test army list
	if len(player.ArmyList) != 1 {
		t.Errorf("Expected army list to have 1 unit, got %d", len(player.ArmyList))
	}

	if player.ArmyList[0].Name != "Space Marine" {
		t.Errorf("Expected first unit to be 'Space Marine', got '%s'", player.ArmyList[0].Name)
	}

	if player.ArmyList[0].Points != 100 {
		t.Errorf("Expected first unit to have 100 points, got %d", player.ArmyList[0].Points)
	}
}

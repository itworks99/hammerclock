package main

import (
	"os"
	"testing"
	"time"
)

// TestReadSettingsFile tests that the settings file can be read correctly
// and created if it doesn't exist
func TestReadSettingsFile(t *testing.T) {
	// First, save the current state of the file
	fileExists := true
	_, err := os.Stat(defaultSettingsFilename)
	if os.IsNotExist(err) {
		fileExists = false
	}

	// If the file exists, make a backup
	var backupData []byte
	if fileExists {
		backupData, err = os.ReadFile(defaultSettingsFilename)
		if err != nil {
			t.Fatalf("Failed to read existing settings file: %v", err)
		}
	}

	// Remove the settings file if it exists (for testing)
	err = os.Remove(defaultSettingsFilename)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove settings file: %v", err)
	}

	// Verify the file was removed
	_, err = os.Stat(defaultSettingsFilename)
	if !os.IsNotExist(err) {
		t.Fatalf("Settings file still exists after removal")
	}

	// Call readSettingsFile to test if it creates the file
	settings := readSettingsFile(defaultSettingsFilename)

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

	// Verify that the file was created
	_, err = os.Stat(defaultSettingsFilename)
	if os.IsNotExist(err) {
		t.Errorf("Settings file was not created")
	}

	// Restore the original file if it existed
	if fileExists {
		err = os.WriteFile(defaultSettingsFilename, backupData, 0644)
		if err != nil {
			t.Fatalf("Failed to restore settings file: %v", err)
		}
	}
}

// TestLoadCustomSettingsFile tests that a custom settings file can be loaded
func TestLoadCustomSettingsFile(t *testing.T) {
	// Path to the custom settings file
	customSettingsFile := "rules/chess.json"

	// Check if the file exists
	_, err := os.Stat(customSettingsFile)
	if os.IsNotExist(err) {
		t.Fatalf("Custom settings file %s does not exist", customSettingsFile)
	}

	// Load the custom settings
	settings := readSettingsFile(customSettingsFile)

	// Check that the settings are loaded correctly
	if settings.Name != "Chess" {
		t.Errorf("Expected game name to be 'Chess', got '%s'", settings.Name)
	}

	if settings.PlayerCount != 2 {
		t.Errorf("Expected player count to be 2, got %d", settings.PlayerCount)
	}

	if len(settings.Phases) != 3 {
		t.Errorf("Expected 3 phases, got %d", len(settings.Phases))
	}

	// Check the phases
	expectedPhases := []string{"Opening", "Middle Game", "End Game"}
	for i, phase := range settings.Phases {
		if phase != expectedPhases[i] {
			t.Errorf("Expected phase %d to be '%s', got '%s'", i, expectedPhases[i], phase)
		}
	}

	// Check player names
	expectedPlayerNames := []string{"White", "Black"}
	for i, name := range settings.PlayerNames {
		if name != expectedPlayerNames[i] {
			t.Errorf("Expected player name %d to be '%s', got '%s'", i, expectedPlayerNames[i], name)
		}
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

	if player.armyList[0].Name != "Space Marine" {
		t.Errorf("Expected first unit to be 'Space Marine', got '%s'", player.armyList[0].Name)
	}
}

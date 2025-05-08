package test

import (
	"hammerclock/components"
	"hammerclock/config"
	"os"
	"testing"
)

// TestLoadOptionsFile tests that the settings file can be read correctly
// and created if it doesn't exist
func TestLoadOptionsFile(t *testing.T) {
	// First, save the current state of the file
	fileExists := true
	_, err := os.Stat(hammerclockConfig.DefaultOptionsFilename)
	if os.IsNotExist(err) {
		fileExists = false
	}

	// If the file exists, make a backup
	var backupData []byte
	if fileExists {
		backupData, err = os.ReadFile(hammerclockConfig.DefaultOptionsFilename)
		if err != nil {
			t.Fatalf("Failed to read existing settings file: %v", err)
		}
	}

	// Remove the settings file if it exists (for testing)
	err = os.Remove(hammerclockConfig.DefaultOptionsFilename)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove settings file: %v", err)
	}

	// Verify the file was removed
	_, err = os.Stat(hammerclockConfig.DefaultOptionsFilename)
	if !os.IsNotExist(err) {
		t.Fatalf("options file still exists after removal")
	}

	// Call LoadOptions to test if it creates the file
	loadedOptions := components.LoadOptions(hammerclockConfig.DefaultOptionsFilename)

	// Check that the settings are loaded correctly
	if loadedOptions.Name != "W40K 10th Edition" {
		t.Errorf("Expected game name to be 'W40K 10th Edition', got '%s'", loadedOptions.Name)
	}

	if loadedOptions.PlayerCount != 2 {
		t.Errorf("Expected player count to be 2, got %d", loadedOptions.PlayerCount)
	}

	if len(loadedOptions.Phases) != 6 {
		t.Errorf("Expected 6 phases, got %d", len(loadedOptions.Phases))
	}

	// Check the first and last phase
	if loadedOptions.Phases[0] != "Command Phase" {
		t.Errorf("Expected first phase to be 'Command Phase', got '%s'", loadedOptions.Phases[0])
	}

	if loadedOptions.Phases[5] != "End Phase" {
		t.Errorf("Expected last phase to be 'End Phase', got '%s'", loadedOptions.Phases[5])
	}

	// Verify that the file was created
	_, err = os.Stat(hammerclockConfig.DefaultOptionsFilename)
	if os.IsNotExist(err) {
		t.Errorf("options file was not created")
	}

	// Restore the original file if it existed
	if fileExists {
		err = os.WriteFile(hammerclockConfig.DefaultOptionsFilename, backupData, 0644)
		if err != nil {
			t.Fatalf("Failed to restore settings file: %v", err)
		}
	}
}

// TestLoadCustomOptionsFile tests that a custom settings file can be loaded
func TestLoadCustomOptionsFile(t *testing.T) {
	// Path to the custom settings file
	customOptionsFile := "../rules/chess.json"

	// Check if the file exists
	_, err := os.Stat(customOptionsFile)
	if os.IsNotExist(err) {
		t.Fatalf("Custom settings file %s does not exist", customOptionsFile)
	}

	// Load the custom settings
	loadedOptions := components.LoadOptions(customOptionsFile)

	// Check that the settings are loaded correctly
	if loadedOptions.Name != "Chess" {
		t.Errorf("Expected game name to be 'Chess', got '%s'", loadedOptions.Name)
	}

	if loadedOptions.PlayerCount != 2 {
		t.Errorf("Expected player count to be 2, got %d", loadedOptions.PlayerCount)
	}

	if len(loadedOptions.Phases) != 3 {
		t.Errorf("Expected 3 phases, got %d", len(loadedOptions.Phases))
	}

	// Check the phases
	expectedPhases := []string{"Opening", "Middle Game", "End Game"}
	for i, phase := range loadedOptions.Phases {
		if phase != expectedPhases[i] {
			t.Errorf("Expected phase %d to be '%s', got '%s'", i, expectedPhases[i], phase)
		}
	}

	// Check player names
	expectedPlayerNames := []string{"White", "Black"}
	for i, name := range loadedOptions.PlayerNames {
		if name != expectedPlayerNames[i] {
			t.Errorf("Expected player name %d to be '%s', got '%s'", i, expectedPlayerNames[i], name)
		}
	}
}

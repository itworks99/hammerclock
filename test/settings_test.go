package test

import (
	"os"
	"testing"

	"hammerclock/pkg/settings"
)

// TestLoadSettingsFile tests that the settings file can be read correctly
// and created if it doesn't exist
func TestLoadSettingsFile(t *testing.T) {
	// First, save the current state of the file
	fileExists := true
	_, err := os.Stat(settings.DefaultSettingsFilename)
	if os.IsNotExist(err) {
		fileExists = false
	}

	// If the file exists, make a backup
	var backupData []byte
	if fileExists {
		backupData, err = os.ReadFile(settings.DefaultSettingsFilename)
		if err != nil {
			t.Fatalf("Failed to read existing settings file: %v", err)
		}
	}

	// Remove the settings file if it exists (for testing)
	err = os.Remove(settings.DefaultSettingsFilename)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove settings file: %v", err)
	}

	// Verify the file was removed
	_, err = os.Stat(settings.DefaultSettingsFilename)
	if !os.IsNotExist(err) {
		t.Fatalf("Settings file still exists after removal")
	}

	// Call LoadSettings to test if it creates the file
	loadedSettings := settings.LoadSettings(settings.DefaultSettingsFilename)

	// Check that the settings are loaded correctly
	if loadedSettings.Name != "W40K 10th Edition" {
		t.Errorf("Expected game name to be 'W40K 10th Edition', got '%s'", loadedSettings.Name)
	}

	if loadedSettings.PlayerCount != 2 {
		t.Errorf("Expected player count to be 2, got %d", loadedSettings.PlayerCount)
	}

	if len(loadedSettings.Phases) != 6 {
		t.Errorf("Expected 6 phases, got %d", len(loadedSettings.Phases))
	}

	// Check the first and last phase
	if loadedSettings.Phases[0] != "Command Phase" {
		t.Errorf("Expected first phase to be 'Command Phase', got '%s'", loadedSettings.Phases[0])
	}

	if loadedSettings.Phases[5] != "End Phase" {
		t.Errorf("Expected last phase to be 'End Phase', got '%s'", loadedSettings.Phases[5])
	}

	// Verify that the file was created
	_, err = os.Stat(settings.DefaultSettingsFilename)
	if os.IsNotExist(err) {
		t.Errorf("Settings file was not created")
	}

	// Restore the original file if it existed
	if fileExists {
		err = os.WriteFile(settings.DefaultSettingsFilename, backupData, 0644)
		if err != nil {
			t.Fatalf("Failed to restore settings file: %v", err)
		}
	}
}

// TestLoadCustomSettingsFile tests that a custom settings file can be loaded
func TestLoadCustomSettingsFile(t *testing.T) {
	// Path to the custom settings file
	customSettingsFile := "../rules/chess.json"

	// Check if the file exists
	_, err := os.Stat(customSettingsFile)
	if os.IsNotExist(err) {
		t.Fatalf("Custom settings file %s does not exist", customSettingsFile)
	}

	// Load the custom settings
	loadedSettings := settings.LoadSettings(customSettingsFile)

	// Check that the settings are loaded correctly
	if loadedSettings.Name != "Chess" {
		t.Errorf("Expected game name to be 'Chess', got '%s'", loadedSettings.Name)
	}

	if loadedSettings.PlayerCount != 2 {
		t.Errorf("Expected player count to be 2, got %d", loadedSettings.PlayerCount)
	}

	if len(loadedSettings.Phases) != 3 {
		t.Errorf("Expected 3 phases, got %d", len(loadedSettings.Phases))
	}

	// Check the phases
	expectedPhases := []string{"Opening", "Middle Game", "End Game"}
	for i, phase := range loadedSettings.Phases {
		if phase != expectedPhases[i] {
			t.Errorf("Expected phase %d to be '%s', got '%s'", i, expectedPhases[i], phase)
		}
	}

	// Check player names
	expectedPlayerNames := []string{"White", "Black"}
	for i, name := range loadedSettings.PlayerNames {
		if name != expectedPlayerNames[i] {
			t.Errorf("Expected player name %d to be '%s', got '%s'", i, expectedPlayerNames[i], name)
		}
	}
}

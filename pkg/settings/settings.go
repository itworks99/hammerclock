package settings

import (
	"encoding/json"
	"fmt"
	"os"

	"hammerclock/internal/app"
)

// DefaultSettingsFilename is the default filename for the settings file
const DefaultSettingsFilename = "defaultRules.json"

// LoadSettings loads the settings from a file
func LoadSettings(filename string) app.Settings {
	var settings app.Settings

	// Check if the settings file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// If the requested file is not the default one, inform the user and use default
		if filename != DefaultSettingsFilename {
			fmt.Printf("Settings file '%s' not found, using default settings file\n", filename)
			return LoadSettings(DefaultSettingsFilename)
		}

		// Default file doesn't exist, create it from defaultSettings
		fmt.Println("Default settings file not found, creating it")

		// Convert defaultSettings to JSON
		jsonData, err := json.MarshalIndent(app.DefaultSettings, "", "  ")
		if err != nil {
			fmt.Println("Error creating default settings file:", err)
			return app.DefaultSettings
		}

		// Write the JSON data to the file
		err = os.WriteFile(DefaultSettingsFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing default settings file:", err)
			return app.DefaultSettings
		}

		return app.DefaultSettings
	} else if err != nil {
		// Some other error occurred
		fmt.Printf("Error checking settings file '%s': %v\n", filename, err)
		if filename != DefaultSettingsFilename {
			fmt.Println("Falling back to default settings")
			return LoadSettings(DefaultSettingsFilename)
		}
		return app.DefaultSettings
	}

	// File exists, read it
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading settings file '%s': %v\n", filename, err)
		if filename != DefaultSettingsFilename {
			fmt.Println("Falling back to default settings")
			return LoadSettings(DefaultSettingsFilename)
		}
		return app.DefaultSettings
	}

	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		fmt.Printf("Error processing settings file '%s': %v\n", filename, err)
		if filename != DefaultSettingsFilename {
			fmt.Println("Falling back to default settings")
			return LoadSettings(DefaultSettingsFilename)
		}
		return app.DefaultSettings
	}

	return settings
}

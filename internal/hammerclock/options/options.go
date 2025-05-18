package options

import (
	"encoding/json"
	"fmt"
	"os"

	"hammerclock/internal/hammerclock/config"
	"hammerclock/internal/hammerclock/rules"
)

// Options defines the configuration for a game, including player details, phases, and display preferences.
type Options struct {
	Default        int           `json:"default"`
	Rules          []rules.Rules `json:"rules"`
	PlayerCount    int           `json:"playerCount"`
	PlayerNames    []string      `json:"playerNames"`
	ColorPalette   string        `json:"colorPalette"`
	TimeFormat     string        `json:"timeFormat"`     // AMPM or 24h
	LoggingEnabled bool          `json:"loggingEnabled"` // Enable/disable CSV logging
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
	Default:        0,
	Rules:          rules.AllRules,
	PlayerCount:    hammerclockConfig.DefaultPlayerCount,
	PlayerNames:    DefaultPlayerNames(),
	ColorPalette:   hammerclockConfig.DefaultColorPalette,
	TimeFormat:     "AMPM",
	LoggingEnabled: true, // CSV logging enabled by default
}

// LoadOptions loads the options from a file
func LoadOptions(filename string) Options {
	var opts Options

	// Check if the options file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// If the requested file is not the default one, inform the user and use default
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Printf("options file '%s' not found, using default options file\n", filename)
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}

		// Default file doesn't exist, create it from defaultOptions
		fmt.Println("Default options file not found, creating it")

		// Convert default options to JSON
		jsonData, err := json.MarshalIndent(DefaultOptions, "", "  ")
		if err != nil {
			fmt.Println("Error creating default options file:", err)
			return DefaultOptions
		}

		// Write the JSON data to the file
		err = os.WriteFile(hammerclockConfig.DefaultOptionsFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing default options file:", err)
			return DefaultOptions
		}

		return DefaultOptions
	} else if err != nil {
		// Some other error occurred
		fmt.Printf("Error checking options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return DefaultOptions
	}

	// File exists, read it
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return DefaultOptions
	}

	// Unmarshal the JSON data into the options struct
	err = json.Unmarshal(byteValue, &opts)
	if err != nil {
		fmt.Printf("Error parsing options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return DefaultOptions
	}

	return opts
}

// SaveOptions saves the options to a file
func SaveOptions(opts Options, filename string, silent bool) error {
	// If no filename is specified, use the default
	if filename == "" {
		filename = hammerclockConfig.DefaultOptionsFilename
	}

	// Convert options to JSON
	jsonData, err := json.MarshalIndent(opts, "", "  ")
	if err != nil {
		if !silent {
			fmt.Printf("Error marshalling options: %v\n", err)
		}
		return err
	}

	// Write the JSON data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil && !silent {
		fmt.Printf("Error writing options file '%s': %v\n", filename, err)
	}

	return err
}

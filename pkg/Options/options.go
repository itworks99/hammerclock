package Options

import (
	"encoding/json"
	"fmt"
	"os"

	"hammerclock/internal/app"
)

// DefaultOptionsFilename is the default filename for the options file
const DefaultOptionsFilename = "defaultRules.json"

// LoadOptions loads the options from a file
func LoadOptions(filename string) app.Options {
	var options app.Options

	// Check if the options file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// If the requested file is not the default one, inform the user and use default
		if filename != DefaultOptionsFilename {
			fmt.Printf("Options file '%s' not found, using default options file\n", filename)
			return LoadOptions(DefaultOptionsFilename)
		}

		// Default file doesn't exist, create it from defaultOptions
		fmt.Println("Default options file not found, creating it")

		// Convert defaultOptions to JSON
		jsonData, err := json.MarshalIndent(app.DefaultOptions, "", "  ")
		if err != nil {
			fmt.Println("Error creating default options file:", err)
			return app.DefaultOptions
		}

		// Write the JSON data to the file
		err = os.WriteFile(DefaultOptionsFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing default options file:", err)
			return app.DefaultOptions
		}

		return app.DefaultOptions
	} else if err != nil {
		// Some other error occurred
		fmt.Printf("Error checking options file '%s': %v\n", filename, err)
		if filename != DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(DefaultOptionsFilename)
		}
		return app.DefaultOptions
	}

	// File exists, read it
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading options file '%s': %v\n", filename, err)
		if filename != DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(DefaultOptionsFilename)
		}
		return app.DefaultOptions
	}

	err = json.Unmarshal(byteValue, &options)
	if err != nil {
		fmt.Printf("Error processing options file '%s': %v\n", filename, err)
		if filename != DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(DefaultOptionsFilename)
		}
		return app.DefaultOptions
	}

	return options
}

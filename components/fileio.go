package components

import (
	"encoding/json"
	"fmt"
	"os"

	"hammerclock/config"
	"hammerclock/internal/app"
)

// LoadOptions loads the options from a file
func LoadOptions(filename string) app.Options {
	var options app.Options

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
		jsonData, err := json.MarshalIndent(app.DefaultOptions, "", "  ")
		if err != nil {
			fmt.Println("Error creating default options file:", err)
			return app.DefaultOptions
		}

		// Write the JSON data to the file
		err = os.WriteFile(hammerclockConfig.DefaultOptionsFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing default options file:", err)
			return app.DefaultOptions
		}

		return app.DefaultOptions
	} else if err != nil {
		// Some other error occurred
		fmt.Printf("Error checking options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return app.DefaultOptions
	}

	// File exists, read it
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return app.DefaultOptions
	}

	err = json.Unmarshal(byteValue, &options)
	if err != nil {
		fmt.Printf("Error processing options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return app.DefaultOptions
	}

	return options
}

// SaveOptions saves the options to a file.
// If filename is empty, it uses DefaultOptionsFilename.
// If overwrite is false and the file exists, it returns an error.
// It returns an error if the operation fails.
func SaveOptions(options app.Options, filename string, overwrite bool) error {
	// If filename is empty, use default
	if filename == "" {
		filename = hammerclockConfig.DefaultOptionsFilename
	}

	// Check if the file exists and we're not allowed to overwrite
	if !overwrite {
		_, err := os.Stat(filename)
		if err == nil {
			// File exists
			return fmt.Errorf("file '%s' already exists and overwrite is set to false", filename)
		} else if !os.IsNotExist(err) {
			// Some other error occurred when checking file statusbar
			return fmt.Errorf("error checking if file '%s' exists: %w", filename, err)
		}
	}

	// Convert options to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling options to JSON: %w", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing options to file '%s': %w", filename, err)
	}

	return nil
}

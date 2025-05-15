package fileio

import (
	"encoding/json"
	"fmt"
	"os"

	hammerclockConfig "hammerclock/config"
	"hammerclock/internal/app/options"
)

// LoadOptions loads the options from a file
func LoadOptions(filename string) options.Options {
	var opts options.Options

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
		jsonData, err := json.MarshalIndent(options.DefaultOptions, "", "  ")
		if err != nil {
			fmt.Println("Error creating default options file:", err)
			return options.DefaultOptions
		}

		// Write the JSON data to the file
		err = os.WriteFile(hammerclockConfig.DefaultOptionsFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing default options file:", err)
			return options.DefaultOptions
		}

		return options.DefaultOptions
	} else if err != nil {
		// Some other error occurred
		fmt.Printf("Error checking options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return options.DefaultOptions
	}

	// File exists, read it
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return options.DefaultOptions
	}

	// Unmarshal the JSON data into the options struct
	err = json.Unmarshal(byteValue, &opts)
	if err != nil {
		fmt.Printf("Error parsing options file '%s': %v\n", filename, err)
		if filename != hammerclockConfig.DefaultOptionsFilename {
			fmt.Println("Falling back to default options")
			return LoadOptions(hammerclockConfig.DefaultOptionsFilename)
		}
		return options.DefaultOptions
	}

	return opts
}

// SaveOptions saves the options to a file
func SaveOptions(opts options.Options, filename string, silent bool) error {
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

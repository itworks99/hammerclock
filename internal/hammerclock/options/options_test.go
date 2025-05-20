package options

import (
	"encoding/json"
	"os"
	"testing"

	"hammerclock/internal/hammerclock/config"
)

func TestLoadOptionsFromNonExistentFileUsesDefaultOptions(t *testing.T) {
	filename := "nonexistent.json"
	opts := LoadOptions(filename)

	if opts.Default != DefaultOptions.Default {
		t.Errorf("Expected default options, got %+v", opts)
	}
}

func TestLoadOptionsFromInvalidFileFallsBackToDefault(t *testing.T) {
	filename := "invalid.json"
	err := os.WriteFile(filename, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}
	defer os.Remove(filename)

	opts := LoadOptions(filename)
	if opts.Default != DefaultOptions.Default {
		t.Errorf("Expected default options, got %+v", opts)
	}
}

func TestSaveOptionsCreatesFileWithCorrectContent(t *testing.T) {
	filename := "test_options.json"
	defer os.Remove(filename)

	err := SaveOptions(DefaultOptions, filename, false)
	if err != nil {
		t.Fatalf("Failed to save options: %v", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	var loadedOpts Options
	err = json.Unmarshal(data, &loadedOpts)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved options: %v", err)
	}

	if loadedOpts.Default != DefaultOptions.Default {
		t.Errorf("Expected saved options to match default options, got %+v", loadedOpts)
	}
}

func TestSaveOptionsHandlesEmptyFilenameGracefully(t *testing.T) {
	err := SaveOptions(DefaultOptions, "", false)
	if err != nil {
		t.Errorf("Expected no error when saving with empty filename, got %v", err)
	}
}

func TestLoadOptionsHandlesCorruptedDefaultFileGracefully(t *testing.T) {
	filename := hammerclockConfig.DefaultOptionsFilename
	err := os.WriteFile(filename, []byte("corrupted json"), 0644)
	if err != nil {
		t.Fatalf("Failed to create corrupted default file: %v", err)
	}
	defer os.Remove(filename)

	opts := LoadOptions(filename)
	if opts.Default != DefaultOptions.Default {
		t.Errorf("Expected fallback to default options, got %+v", opts)
	}
}

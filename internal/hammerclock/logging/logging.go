// Package logging provides buffered logging capability to write log entries to CSV files
package logging

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Buffered channel for log entries
var logChannel chan LogEntry
var logInitialized bool
var logWg sync.WaitGroup
var logMutex sync.Mutex

// Initialise sets up the background log writer
func Initialise() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logInitialized {
		return
	}

	logChannel = make(chan LogEntry, 100)
	logWg.Add(1)
	// Start background log writer
	go func() {
		defer logWg.Done()
		defer func() {
			// Recover from any panics in the background goroutine
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic in log writer: %v\n", r)
			}
		}()

		for entry := range logChannel {
			WriteLogEntry(entry)
		}
	}()
	logInitialized = true
}

// Cleanup closes the log channel and waits for the background writer to finish
func Cleanup() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if !logInitialized {
		return
	}

	close(logChannel)
	logWg.Wait()
	logInitialized = false
}

// AddLogEntry sends a log entry to the buffered channel if enableLogging is true
func AddLogEntry(entry LogEntry, enableLogging bool) {
	if !enableLogging {
		// Skip logging but print debug info if in debug mode
		return
	}

	// Make sure logging is initialized
	if !logInitialized {
		Initialise()
	}

	select {
	case logChannel <- entry:
		// sent successfully
	default:
		// channel full, drop log entry to avoid UI lag
		return
	}
}

// WriteLogEntry appends a LogEntry to logs.csv in CSV format.
func WriteLogEntry(entry LogEntry) {
	// Use default log directory (current working directory)
	logDir := "."
	fileName := "logs.csv"
	filePath := fileName

	// If an environment variable for log directory is set, use it
	if envLogDir := os.Getenv("HAMMERCLOCK_LOG_DIR"); envLogDir != "" {
		logDir = envLogDir
		filePath = filepath.Join(logDir, fileName)

		// Ensure the directory exists
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Error creating log directory %s: %v\n", logDir, err)
			// Fall back to current directory
			filePath = fileName
		}
	}

	fileExists := false

	// Check if file exists before opening
	if _, err := os.Stat(filePath); err == nil {
		fileExists = true
	}

	// Open file with appropriate flags
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing log file: %v\n", err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header if it's a new file
	if !fileExists {
		if err := writer.Write([]string{"DateTime", "PlayerName", "Turn", "Phase", "Message"}); err != nil {
			fmt.Printf("Error writing CSV header: %v\n", err)
			return
		}
	}

	// Write the log entry data
	if err := writer.Write([]string{
		entry.DateTime,
		entry.PlayerName,
		fmt.Sprintf("%d", entry.Turn),
		entry.Phase,
		entry.Message,
	}); err != nil {
		fmt.Printf("Error writing CSV entry: %v\n", err)
	}
}

// LogEntry represents a single log entry with details about an action.
type LogEntry struct {
	DateTime   string
	PlayerName string
	Turn       int
	Phase      string
	Message    string
}

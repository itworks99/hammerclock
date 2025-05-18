// Package logging provides buffered logging capability to write log entries to CSV files
package logging

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"hammerclock/internal/hammerclock/common"
	"hammerclock/internal/hammerclock/config"
)

// Buffered channel for log entries
var logChannel chan common.LogEntry
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

	logChannel = make(chan common.LogEntry, 100)
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
			writeLogEntry(entry)
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

// SendLogEntry sends a log entry to the buffered channel if enableLogging is true
func SendLogEntry(entry common.LogEntry) {
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

// writeLogEntry appends a LogEntry to logs.csv in CSV format.
func writeLogEntry(entry common.LogEntry) {
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

// AddLogEntry adds a log entry to a player's action log
func AddLogEntry(player *common.Player, model *common.Model, format string, args ...any) {
	currentPhase := ""
	if player.CurrentPhase < len(model.Options.Rules[model.Options.Default].Phases) && player.CurrentPhase >= 0 {
		currentPhase = model.Options.Rules[model.Options.Default].Phases[player.CurrentPhase]
	}

	logEntry := common.LogEntry{
		DateTime:   time.Now().Local().Format(hammerclockConfig.DefaultLogDateTimeFormat),
		PlayerName: player.Name,
		Turn:       player.TurnCount,
		Phase:      currentPhase,
		Message:    fmt.Sprintf(format, args...),
	}

	// Add to in-memory player action log for UI
	player.ActionLog = append(player.ActionLog, logEntry)

	// Send log entry to the logging channel
	SendLogEntry(logEntry)
}

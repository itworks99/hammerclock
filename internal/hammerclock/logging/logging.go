// Package logging provides buffered logging capability to write log entries to CSV files
package logging

import (
	"fmt"
	"sync"

	"hammerclock/internal/hammerclock/ui"
)

// Buffered channel for log entries
var logChannel chan ui.LogEntry
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

	logChannel = make(chan ui.LogEntry, 100)
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
			ui.WriteLogEntry(entry)
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

// WriteLogEntry sends a log entry to the buffered channel if enableLogging is true
func WriteLogEntry(entry ui.LogEntry, enableLogging bool) {
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

package ui

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateLogView initializes a scrollable, word-wrapped text view for logs.
func CreateLogView() *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetWrap(true)
}

// CreateLogContainer wraps the log view in a container with auto-scroll and input handling.
func CreateLogContainer(logView *tview.TextView) *tview.Flex {
	logView.SetChangedFunc(
		func() {
			logView.ScrollToEnd()
		},
	) // Auto-scroll on content change.
	SetupLogViewInputHandling(logView) // Enable mouse scrolling.

	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(logView, 0, 1, true)
}

// SetupLogViewInputHandling enables mouse scrolling for the log view.
func SetupLogViewInputHandling(logView *tview.TextView) {
	logView.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		row, _ := logView.GetScrollOffset()
		switch action {
		case tview.MouseScrollUp:
			if row > 0 {
				logView.ScrollTo(row-1, 0)
			}
		case tview.MouseScrollDown:
			if row < logView.GetWrappedLineCount()-1 {
				logView.ScrollTo(row+1, 0)
			}
		default:
			// Handle other mouse actions as needed
		}
		return action, nil
	})
}

// SetLogContent updates the log view with the provided log entries.
func SetLogContent(logView *tview.TextView, logEntries interface{}) {
	if logView == nil {
		return
	}

	var logText strings.Builder
	switch entries := logEntries.(type) {
	case []interface{}:
		for _, entry := range entries {
			logText.WriteString(formatLogEntry(entry))
		}
	case []string:
		for _, entry := range entries {
			logText.WriteString(entry + "\n")
		}
	default:
		if slice, ok := tryGetSlice(logEntries); ok {
			for _, entry := range slice {
				logText.WriteString(formatLogEntry(entry))
			}
		} else {
			logText.WriteString(formatLogEntry(logEntries))
		}
	}

	if logText.String() != logView.GetText(false) {
		logView.SetText(logText.String())
	}
}

// formatLogEntry converts a log entry to a string.
func formatLogEntry(entry interface{}) string {
	if logEntry, ok := entry.(LogEntry); ok {
		return logEntry.DisplayString() + "\n"
	}
	if stringer, ok := entry.(fmt.Stringer); ok {
		return stringer.String() + "\n"
	}
	return fmt.Sprintf("%v\n", entry)
}

// tryGetSlice converts an interface to a slice of interfaces if possible.
func tryGetSlice(obj interface{}) ([]interface{}, bool) {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Slice {
		return nil, false
	}

	result := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = val.Index(i).Interface()
	}
	return result, true
}

// CreateString returns a detailed string representation of the log entry.
func (le LogEntry) CreateString() string {
	return fmt.Sprintf("[%s] %s | Turn %d | %s | %s",
		le.DateTime, le.PlayerName, le.Turn, le.Phase, le.Message)
}

// DisplayString returns a simplified string representation for UI display.
func (le LogEntry) DisplayString() string {
	if spaceIdx := strings.Index(le.DateTime, " "); spaceIdx != -1 {
		return fmt.Sprintf("[%s] %s", le.DateTime[spaceIdx+1:], le.Message)
	}
	return fmt.Sprintf("[%s] %s", le.DateTime, le.Message)
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

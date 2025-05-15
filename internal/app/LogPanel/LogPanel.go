package LogPanel

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateLogView creates a text view for the log with scrolling enabled.
func CreateLogView() *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetScrollable(true). // Enable scrolling
		SetWordWrap(true). // Enable word wrapping
		SetWrap(true) // Enable text wrapping
}

// CreateLogContainer creates a container with a log view.
func CreateLogContainer(logView *tview.TextView) *tview.Flex {
	// Create a layout for the log view
	logContainer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(logView, 0, 1, true) // Log view takes all the space

	// Set auto-scroll behavior when content changes
	logView.SetChangedFunc(func() {
		logView.ScrollToEnd()
	})

	// Setup mouse handling for scrolling
	SetupLogViewInputHandling(logView)

	return logContainer
}

// SetupLogViewInputHandling configures keyboard and mouse input for scrolling a log view.
func SetupLogViewInputHandling(logView *tview.TextView) {
	// Set up mouse handling for scrolling
	logView.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseScrollUp:
			// Scroll up
			row, _ := logView.GetScrollOffset()
			if row > 0 {
				logView.ScrollTo(row-1, 0)
			}
			return action, nil // Consume the event
		case tview.MouseScrollDown:
			// Scroll down
			row, _ := logView.GetScrollOffset()
			totalLines := logView.GetWrappedLineCount()
			if row < totalLines-1 {
				logView.ScrollTo(row+1, 0)
			}
			return action, nil // Consume the event
		default:
			// Handle other mouse events
		}

		// Pass other mouse events through
		return action, event
	})
}

// SetLogContent updates the content of a log view with the given log entries.
// logEntries must have a String() method.
func SetLogContent(logView *tview.TextView, logEntries interface{}) {
	if logView == nil {
		return
	}

	var logText strings.Builder

	// Handle different types of log entries
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

	// Update content only if it has changed
	if logText.String() != logView.GetText(false) {
		logView.SetText(logText.String())
	}
}

// formatLogEntry formats a single log entry as a string.
func formatLogEntry(entry interface{}) string {
	// Check if it's a LogEntry type to use DisplayString
	if logEntry, ok := entry.(LogEntry); ok {
		return logEntry.DisplayString() + "\n"
	}
	// Fall back to regular String method for other types
	if stringer, ok := entry.(fmt.Stringer); ok {
		return stringer.String() + "\n"
	}
	return fmt.Sprintf("%v\n", entry)
}

// tryGetSlice attempts to convert an interface to a slice of interfaces
func tryGetSlice(obj interface{}) ([]interface{}, bool) {
	if obj == nil {
		return nil, false
	}

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Slice {
		return nil, false
	}

	length := val.Len()
	result := make([]interface{}, length)
	for i := 0; i < length; i++ {
		result[i] = val.Index(i).Interface()
	}

	return result, true
}

// LogEntry represents a single entry in a player's action log
type LogEntry struct {
	DateTime   string
	PlayerName string
	Turn       int
	Phase      string
	Message    string
}

// String returns a formatted string representation of the log entry
// This is used when the entry needs to be displayed with all details
func (le LogEntry) String() string {
	return fmt.Sprintf("[%s] %s | Turn %d | %s | %s", 
		le.DateTime,
		le.PlayerName,
		le.Turn,
		le.Phase,
		le.Message)
}

// DisplayString returns a simplified string representation for UI display
// This omits player name, turn, and phase since they're visible elsewhere in the UI
// Also extracts only the time part from the DateTime field
func (le LogEntry) DisplayString() string {
	// Extract just the time part (the last 8 characters in the format "15:04:05")
	timeOnly := le.DateTime
	
	// If the DateTime follows the format "2006-01-02 15:04:05"
	// Extract just the time part (everything after the space)
	if spaceIdx := strings.Index(le.DateTime, " "); spaceIdx != -1 && spaceIdx+1 < len(le.DateTime) {
		timeOnly = le.DateTime[spaceIdx+1:]
	}
	
	return fmt.Sprintf("[%s] %s", 
		timeOnly,
		le.Message)
}

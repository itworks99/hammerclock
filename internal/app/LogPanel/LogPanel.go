package LogPanel

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CreateLogView creates a text view for the log with scrolling enabled.
func CreateLogView() *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetScrollable(true). // Enable scrolling
		SetWordWrap(true).   // Enable word wrapping
		SetWrap(true)        // Enable text wrapping
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

	// Setup keyboard and mouse handling for scrolling
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
		if len(entries) == 0 {
			return
		}
		for _, entry := range entries {
			if stringer, ok := entry.(fmt.Stringer); ok {
				logText.WriteString(stringer.String() + "\n")
			} else {
				logText.WriteString(fmt.Sprintf("%v\n", entry))
			}
		}
	case []string:
		if len(entries) == 0 {
			return
		}
		for _, entry := range entries {
			logText.WriteString(entry + "\n")
		}
	default:
		// For any other type, try to get the slice and iterate over it
		if slice, ok := tryGetSlice(logEntries); ok && len(slice) > 0 {
			for _, entry := range slice {
				if stringer, ok := entry.(fmt.Stringer); ok {
					logText.WriteString(stringer.String() + "\n")
				} else {
					logText.WriteString(fmt.Sprintf("%v\n", entry))
				}
			}
		} else {
			// Single item
			if stringer, ok := logEntries.(fmt.Stringer); ok {
				logText.WriteString(stringer.String() + "\n")
			} else {
				logText.WriteString(fmt.Sprintf("%v\n", logEntries))
			}
		}
	}

	// Only update if content has changed
	if logText.String() != logView.GetText(false) {
		logView.SetText(logText.String())
	}
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
	for i := range length {
		result[i] = val.Index(i).Interface()
	}

	return result, true
}

// LogEntry represents a single entry in a player's action log
type LogEntry struct {
	DateTime   time.Time
	PlayerName string
	Turn       int
	Phase      string
	Message    string
}

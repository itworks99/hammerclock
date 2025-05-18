package ui

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"hammerclock/internal/hammerclock/logging"
)

// CreateLogView initializes a scrollable, word-wrapped text view for logs.
func CreateLogView() *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)
}

// CreateLogContainer wraps the log view in a container with auto-scroll and input handling.
func CreateLogContainer(logView *tview.TextView) *tview.Flex {
	logView.SetChangedFunc(func() { logView.ScrollToEnd() }) // Auto-scroll on content change.
	SetupLogViewInputHandling(logView)                       // Enable mouse scrolling.

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
			logView.ScrollTo(row-1, 0)
		case tview.MouseScrollDown:
			logView.ScrollTo(row+1, 0)
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
	switch e := entry.(type) {
	case logging.LogEntry:
		return DisplayLogEntry(e) + "\n"
	case fmt.Stringer:
		return e.String() + "\n"
	default:
		return fmt.Sprintf("%v\n", entry)
	}
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

// DisplayLogEntry returns a simplified string representation for UI display.
func DisplayLogEntry(le logging.LogEntry) string {
	if spaceIdx := strings.Index(le.DateTime, " "); spaceIdx != -1 {
		return fmt.Sprintf("[%s] %s", le.DateTime[spaceIdx+1:], le.Message)
	}
	return fmt.Sprintf("[%s] %s", le.DateTime, le.Message)
}

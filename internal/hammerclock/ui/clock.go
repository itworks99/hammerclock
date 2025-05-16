package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TimeFormat determines the clock format string based on the model's time format setting (AMPM or 24-hour).
func TimeFormat(option string) string {
	// Determine the clock layout based on the options
	if option == "AMPM" {
		return "03:04:05 PM"
	}
	return "15:04:05"
}

// Display displays the current time in the specified format.
func Display(format string, color tcell.Color) *tview.TextView {
	hClock := tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true).
		SetTextColor(color)

	// Set the clock format based on the model's options
	var clockFormat = TimeFormat(format)

	hClock.SetText(time.Now().Format(clockFormat))
	return hClock
}

// TimeFormatToIndex converts the time format string to an index
func TimeFormatToIndex(format string) int {
	if format == "AMPM" {
		return 0
	}
	return 1 // Default to 24-hour format
}

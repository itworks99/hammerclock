package common

import "github.com/gdamore/tcell/v2"

// PrevPhaseMsg is sent when the user wants to move to the previous phase
type PrevPhaseMsg struct{}

// ShowOptionsMsg is sent when the user wants to show the options screen
type ShowOptionsMsg struct{}

// ShowAboutMsg is sent when the user wants to show the about screen
type ShowAboutMsg struct{}

// ShowMainScreenMsg is sent when the user wants to return to the main screen
type ShowMainScreenMsg struct{}

// TickMsg is sent every second to update the clock and player times
type TickMsg struct{}

// KeyPressMsg is sent when a key is pressed
type KeyPressMsg struct {
	Key  tcell.Key
	Rune rune
}

// EndGameMsg is sent when the user wants to end the current game
type EndGameMsg struct{}

// EndGameConfirmMsg is sent when the user confirms or cancels ending the game
type EndGameConfirmMsg struct {
	Confirmed bool
}

// ShowEndGameConfirmMsg is sent to show the end game confirmation dialog
type ShowEndGameConfirmMsg struct{}

// ShowExitConfirmMsg is sent to show the exit confirmation dialog
type ShowExitConfirmMsg struct{}

// ExitConfirmMsg is sent when the user confirms or cancels exiting the application
type ExitConfirmMsg struct {
	Confirmed bool
}

// ShowModalMsg is sent to show a modal dialog
type ShowModalMsg struct {
	Type string
}

// RestoreMainUIMsg is sent to restore the main UI after a modal dialog
type RestoreMainUIMsg struct{}

// SetRulesetMsg is sent when the user selects a different ruleset
type SetRulesetMsg struct {
	Index int
}

// SetPlayerCountMsg is sent when the user changes the player count
type SetPlayerCountMsg struct {
	Count int
}

// SetPlayerNameMsg is sent when a player name is changed
type SetPlayerNameMsg struct {
	Index int
	Name  string
}

// SetColorPaletteMsg is sent when the color palette is changed
type SetColorPaletteMsg struct {
	Name string
}

// SetTimeFormatMsg is sent when the time format is changed
type SetTimeFormatMsg struct {
	Format string
}

// SetOneTurnForAllPlayersMsg is sent when the "One Turn For All Players" option is toggled
type SetOneTurnForAllPlayersMsg struct {
	Value bool
}

// SetEnableLogMsg is sent when the user toggles CSV logging
type SetEnableLogMsg struct {
	Value bool
}

// StartGameMsg is sent when the user wants to start/pause/resume the game
type StartGameMsg struct{}

// SwitchTurnsMsg is sent when the user wants to switch turns
type SwitchTurnsMsg struct{}

// NextPhaseMsg is sent when the user wants to move to the next phase
type NextPhaseMsg struct{}

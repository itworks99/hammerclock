package hammerclock

import "github.com/gdamore/tcell/v2"

// prevPhaseMsg is sent when the user wants to move to the previous phase
type prevPhaseMsg struct{}

// showOptionsMsg is sent when the user wants to show the options screen
type showOptionsMsg struct{}

// showAboutMsg is sent when the user wants to show the about screen
type showAboutMsg struct{}

// showMainScreenMsg is sent when the user wants to return to the main screen
type showMainScreenMsg struct{}

// TickMsg is sent every second to update the clock and player times
type TickMsg struct{}

// keyPressMsg is sent when a key is pressed
type keyPressMsg struct {
	Key  tcell.Key
	Rune rune
}

// endGameMsg is sent when the user wants to end the current game
type endGameMsg struct{}

// endGameConfirmMsg is sent when the user confirms or cancels ending the game
type endGameConfirmMsg struct {
	Confirmed bool
}

// showEndGameConfirmMsg is sent to show the end game confirmation dialog
type showEndGameConfirmMsg struct{}

// ShowModalMsg is sent to show a modal dialog
type ShowModalMsg struct {
	Type string
}

// RestoreMainUIMsg is sent to restore the main UI after a modal dialog
type RestoreMainUIMsg struct{}

// setRulesetMsg is sent when the user selects a different ruleset
type setRulesetMsg struct {
	Index int
}

// setPlayerCountMsg is sent when the user changes the player count
type setPlayerCountMsg struct {
	Count int
}

// setPlayerNameMsg is sent when a player name is changed
type setPlayerNameMsg struct {
	Index int
	Name  string
}

// setColorPaletteMsg is sent when the color palette is changed
type setColorPaletteMsg struct {
	Name string
}

// setTimeFormatMsg is sent when the time format is changed
type setTimeFormatMsg struct {
	Format string
}

// setOneTurnForAllPlayersMsg is sent when the "One Turn For All Players" option is toggled
type setOneTurnForAllPlayersMsg struct {
	Value bool
}

// setEnableLogMsg is sent when the user toggles CSV logging
type setEnableLogMsg struct {
	Value bool
}

// startGameMsg is sent when the user wants to start/pause/resume the game
type startGameMsg struct{}

// switchTurnsMsg is sent when the user wants to switch turns
type switchTurnsMsg struct{}

// nextPhaseMsg is sent when the user wants to move to the next phase
type nextPhaseMsg struct{}

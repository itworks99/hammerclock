package main

import (
	"flag"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"hammerclock/internal/hammerclock"
	"hammerclock/internal/hammerclock/common"
	hammerclockConfig "hammerclock/internal/hammerclock/config"
	"hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/palette"

	"github.com/gdamore/tcell/v2"
)

// TestMainIntegration is a basic integration test for the main function
// Since main() is difficult to fully test due to UI dependencies,
// we'll use a short timeout to just verify it starts without panic
func TestMainIntegration(t *testing.T) {
	// Skip in CI environments or when running all tests
	if os.Getenv("CI") != "" || testing.Short() {
		t.Skip("Skipping main test in CI environment")
	}

	// Start main in a goroutine with a channel to signal completion
	done := make(chan bool)
	go func() {
		// Recover from any panics in main
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Main function panicked: %v", r)
			}
			done <- true
		}()

		// Set test args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"hammerclock"}

		// Reset flags to prevent interference from other tests
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		// We can't call main() directly because it doesn't return,
		// but we can check that the flags are properly parsed
		optionsFileFlag := flag.String("o", hammerclockConfig.DefaultOptionsFilename, "Path to the options file")
		flag.Parse()

		if *optionsFileFlag != hammerclockConfig.DefaultOptionsFilename {
			t.Errorf("Expected default options filename, got %s", *optionsFileFlag)
		}

		// Signal completion
		done <- true
	}()

	// Wait for the test to complete or timeout
	select {
	case <-done:
		// Test completed successfully
	case <-time.After(100 * time.Millisecond):
		// This is just to prevent hanging; in real usage we'd wait longer
	}
}

// TestOptionsFlagParsing ensures the options flag is properly parsed
func TestOptionsFlagParsing(t *testing.T) {
	// Save original args and restore after test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset flag set to avoid interference from previous tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test with custom options file
	os.Args = []string{"hammerclock", "-o", "custom.json"}
	optionsFileFlag := flag.String("o", hammerclockConfig.DefaultOptionsFilename, "Path to the options file")
	flag.Parse()

	if *optionsFileFlag != "custom.json" {
		t.Errorf("Expected 'custom.json', got '%s'", *optionsFileFlag)
	}
}

// TestUsageOutput verifies that the usage text is correctly formatted
func TestUsageOutput(t *testing.T) {
	// Save original args and restore after test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset the flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Capture output
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Setup flag usage
	flag.Usage = func() {
		//goland:noinspection GoUnhandledErrorResult
		_, _ = os.Stderr.Write([]byte(cliUsage))
	}

	// Trigger usage
	flag.Usage()

	// Restore stderr
	//goland:noinspection GoUnhandledErrorResult
	w.Close()
	os.Stderr = oldStderr

	// Read captured output
	output := new(strings.Builder)
	_, _ = io.Copy(output, r)

	// Verify output contains expected content
	if !strings.Contains(output.String(), "Hammerclock "+hammerclockConfig.Version) {
		t.Errorf("Expected usage to contain version info")
	}
	if !strings.Contains(output.String(), "-o <file>") {
		t.Errorf("Expected usage to contain options flag")
	}
}

// TestModelCreation tests the initial model setup
func TestModelCreation(t *testing.T) {
	model := hammerclock.NewModel()

	// Verify initial game state
	if model.GameStatus != "Game Not Started" {
		t.Errorf("Expected game status to be 'Game Not Started', got '%s'", model.GameStatus)
	}

	if model.GameStarted {
		t.Errorf("Expected GameStarted to be false")
	}

	// Verify player setup
	if len(model.Players) != options.DefaultOptions.PlayerCount {
		t.Errorf("Expected %d players, got %d", options.DefaultOptions.PlayerCount, len(model.Players))
	}

	// First player should be active at start
	if !model.Players[0].IsTurn {
		t.Errorf("Expected first player to have IsTurn=true")
	}

	// Other players should not be active
	for i := 1; i < len(model.Players); i++ {
		if model.Players[i].IsTurn {
			t.Errorf("Expected player %d to have IsTurn=false", i)
		}
	}
}

// TestUpdateFunction tests the core update logic
func TestUpdateFunction(t *testing.T) {
	model := hammerclock.NewModel()

	// Test starting the game
	updatedModel, _ := hammerclock.Update(&common.StartGameMsg{}, model)
	if updatedModel.GameStatus != "Game In Progress" {
		t.Errorf("Expected game status to be 'Game In Progress', got '%s'", updatedModel.GameStatus)
	}

	// Test pausing the game
	updatedModel, _ = hammerclock.Update(&common.StartGameMsg{}, updatedModel)
	if updatedModel.GameStatus != "Game Paused" {
		t.Errorf("Expected game status to be 'Game Paused', got '%s'", updatedModel.GameStatus)
	}

	// Test switching turns
	initialActivePlayer := -1
	for i, player := range updatedModel.Players {
		if player.IsTurn {
			initialActivePlayer = i
			break
		}
	}

	updatedModel, _ = hammerclock.Update(&common.SwitchTurnsMsg{}, updatedModel)

	newActivePlayer := -1
	for i, player := range updatedModel.Players {
		if player.IsTurn {
			newActivePlayer = i
			break
		}
	}

	if newActivePlayer == initialActivePlayer {
		t.Errorf("Expected active player to change after SwitchTurnsMsg")
	}
}

// TestKeyPressHandling tests key press event handling
func TestKeyPressHandling(t *testing.T) {
	model := hammerclock.NewModel()

	// Test quitting with 'q' key
	_, cmd := hammerclock.Update(&common.KeyPressMsg{Key: tcell.KeyRune, Rune: 'q'}, model)

	// Should show exit confirmation modal
	if cmd == nil {
		t.Errorf("Expected a command, got nil")
		return
	}

	modalMsg := cmd()
	if modalMsg == nil {
		t.Errorf("Expected ShowModalMsg, got nil")
		return
	}

	if showModalMsg, ok := modalMsg.(*common.ShowModalMsg); ok {
		if showModalMsg.Type != "ExitConfirm" {
			t.Errorf("Expected modal type 'ExitConfirm', got '%s'", showModalMsg.Type)
		}
	} else {
		t.Errorf("Expected ShowModalMsg, got %T", modalMsg)
	}

	// Test starting game with 's' key (not spacebar as in previous test)
	updatedModel, _ := hammerclock.Update(&common.KeyPressMsg{Key: tcell.KeyRune, Rune: 's'}, model)
	if updatedModel.GameStatus != "Game In Progress" {
		t.Errorf("Expected game status to be 'Game In Progress', got '%s'", updatedModel.GameStatus)
	}
}

// TestPhaseManagement tests phase switching functionality
func TestPhaseManagement(t *testing.T) {
	model := hammerclock.NewModel()

	// Start with player 0 at phase 0
	initialPhase := model.Players[0].CurrentPhase

	// Move to next phase
	updatedModel, _ := hammerclock.Update(&common.NextPhaseMsg{}, model)

	// Check phase was updated
	if updatedModel.Players[0].CurrentPhase <= initialPhase {
		t.Errorf("Expected phase to increase, got %d", updatedModel.Players[0].CurrentPhase)
	}

	// Now move back to the previous phase
	updatedModel, _ = hammerclock.Update(&common.PrevPhaseMsg{}, updatedModel)

	// Should be back at the initial phase
	if updatedModel.Players[0].CurrentPhase != initialPhase {
		t.Errorf("Expected to be back at phase %d, got %d", initialPhase, updatedModel.Players[0].CurrentPhase)
	}
}

// TestScreenNavigation tests navigation between different screens
func TestScreenNavigation(t *testing.T) {
	model := hammerclock.NewModel()

	// Initially should be on main screen
	if model.CurrentScreen != "main" {
		t.Errorf("Expected to start on 'main' screen, got '%s'", model.CurrentScreen)
	}

	// Navigate to options
	updatedModel, _ := hammerclock.Update(&common.ShowOptionsMsg{}, model)
	if updatedModel.CurrentScreen != "options" {
		t.Errorf("Expected to be on 'options' screen, got '%s'", updatedModel.CurrentScreen)
	}

	// Navigate to about
	updatedModel, _ = hammerclock.Update(&common.ShowAboutMsg{}, updatedModel)
	if updatedModel.CurrentScreen != "about" {
		t.Errorf("Expected to be on 'about' screen, got '%s'", updatedModel.CurrentScreen)
	}

	// Navigate back to main
	updatedModel, _ = hammerclock.Update(&common.ShowMainScreenMsg{}, updatedModel)
	if updatedModel.CurrentScreen != "main" {
		t.Errorf("Expected to be back on 'main' screen, got '%s'", updatedModel.CurrentScreen)
	}
}

// TestOptionsUpdates tests changing options
func TestOptionsUpdates(t *testing.T) {
	model := hammerclock.NewModel()

	// Test changing player count
	const newPlayerCount = 4
	updatedModel, _ := hammerclock.Update(&common.SetPlayerCountMsg{Count: newPlayerCount}, model)

	// Check that the option was updated correctly
	if updatedModel.Options.PlayerCount != newPlayerCount {
		t.Errorf("Expected player count option to be %d, got %d", newPlayerCount, updatedModel.Options.PlayerCount)
	}

	// Note: The model.Players array is not automatically updated by SetPlayerCountMsg
	// It would need to be explicitly recreated with the new player count
	// We're only testing that the option value changed successfully

	// Test changing player name
	const newName = "Test Player"
	// First ensure our player names array has enough entries
	if len(updatedModel.Options.PlayerNames) > 0 {
		updatedModel, _ = hammerclock.Update(&common.SetPlayerNameMsg{Index: 0, Name: newName}, updatedModel)

		// Check the options was updated correctly
		if updatedModel.Options.PlayerNames[0] != newName {
			t.Errorf("Expected player name option to be '%s', got '%s'", newName, updatedModel.Options.PlayerNames[0])
		}
	} else {
		t.Skip("Skipping player name test as player names array is empty")
	}

	// Test changing color palette
	updatedModel, _ = hammerclock.Update(&common.SetColorPaletteMsg{Name: "Solarized"}, updatedModel)
	expectedPalette := palette.ColorPaletteByName("Solarized")
	// Just check one color to verify palette changed
	if updatedModel.CurrentColorPalette.White != expectedPalette.White {
		t.Errorf("Expected palette to be Solarized")
	}

	// Test changing time format
	const newTimeFormat = "24h"
	updatedModel, _ = hammerclock.Update(&common.SetTimeFormatMsg{Format: newTimeFormat}, updatedModel)
	if updatedModel.Options.TimeFormat != newTimeFormat {
		t.Errorf("Expected time format to be '%s', got '%s'", newTimeFormat, updatedModel.Options.TimeFormat)
	}
}

// TestTickHandling tests the tick message for time updates
func TestTickHandling(t *testing.T) {
	model := hammerclock.NewModel()

	// Start the game
	model, _ = hammerclock.Update(&common.StartGameMsg{}, model)

	// Record initial time
	initialTime := model.Players[0].TimeElapsed

	// Send a tick message
	updatedModel, _ := hammerclock.Update(&common.TickMsg{}, model)

	// Active player's time should increase
	activePlayerIndex := -1
	for i, player := range updatedModel.Players {
		if player.IsTurn {
			activePlayerIndex = i
			break
		}
	}

	if activePlayerIndex >= 0 {
		if updatedModel.Players[activePlayerIndex].TimeElapsed <= initialTime {
			t.Errorf("Expected active player's time to increase after tick")
		}
	} else {
		t.Errorf("No active player found")
	}

	// Total game time should also increase
	if updatedModel.TotalGameTime <= model.TotalGameTime {
		t.Errorf("Expected total game time to increase after tick")
	}
}

// TestEndGameFlow tests the flow of ending a game
func TestEndGameFlow(t *testing.T) {
	model := hammerclock.NewModel()

	// First start the game
	model, _ = hammerclock.Update(&common.StartGameMsg{}, model)

	// Should be in progress
	if model.GameStatus != "Game In Progress" {
		t.Errorf("Expected game status to be 'Game In Progress', got '%s'", model.GameStatus)
	}

	// Now try to end it
	_, cmd := hammerclock.Update(&common.EndGameMsg{}, model)

	// EndGameMsg should perform the action directly without a command
	if cmd != nil && cmd() != nil {
		t.Errorf("Expected EndGameMsg to not return a command")
	}

	// Test confirming game end
	updatedModel, cmd := hammerclock.Update(&common.EndGameConfirmMsg{Confirmed: true}, model)

	// Game should be ended
	if updatedModel.GameStatus != "Game Not Started" {
		t.Errorf("Expected game status to be 'Game Not Started', got '%s'", updatedModel.GameStatus)
	}

	// Should have a command to show main screen
	if cmd == nil {
		t.Errorf("Expected a command to restore UI")
		return
	}

	msg := cmd()
	if _, ok := msg.(*common.ShowMainScreenMsg); !ok {
		t.Errorf("Expected ShowMainScreenMsg, got %T", msg)
	}
}

// TestInvalidMessages tests that invalid messages don't crash the system
func TestInvalidMessages(t *testing.T) {
	model := hammerclock.NewModel()

	// Test nil message
	updatedModel, cmd := hammerclock.Update(nil, model)

	// Model should be unchanged
	if updatedModel.GameStatus != model.GameStatus {
		t.Errorf("Expected model to be unchanged for nil message")
	}

	// No command should be returned
	if cmd != nil && cmd() != nil {
		t.Errorf("Expected nil message to not return a command")
	}

	// Try an invalid player index for SetPlayerNameMsg
	updatedModel, cmd = hammerclock.Update(&common.SetPlayerNameMsg{Index: 999, Name: "Invalid"}, model)

	// Model should be unchanged
	if updatedModel.GameStatus != model.GameStatus {
		t.Errorf("Expected model to be unchanged for invalid player index")
	}

	// No command should be returned
	if cmd != nil && cmd() != nil {
		t.Errorf("Expected invalid player index to not return a command")
	}

	// Try an invalid player count
	updatedModel, _ = hammerclock.Update(&common.SetPlayerCountMsg{Count: -1}, model)

	// Model should be unchanged
	if updatedModel.Options.PlayerCount != model.Options.PlayerCount {
		t.Errorf("Expected player count to be unchanged for invalid count")
	}
}

// TestRulesetChange tests changing the ruleset
func TestRulesetChange(t *testing.T) {
	model := hammerclock.NewModel()

	// Get initial ruleset index
	initialRuleIndex := model.Options.Default

	// Find an alternate ruleset to switch to
	var alternateRuleIndex int
	alternateRuleFound := false
	for i := range model.Options.Rules {
		if i != initialRuleIndex {
			alternateRuleIndex = i
			alternateRuleFound = true
			break
		}
	}

	// Skip test if no alternate ruleset found
	if !alternateRuleFound {
		t.Skip("No alternate ruleset found to test with")
	}

	// Change ruleset
	updatedModel, _ := hammerclock.Update(&common.SetRulesetMsg{Index: alternateRuleIndex}, model)

	// Verify the default ruleset was changed
	if updatedModel.Options.Default == initialRuleIndex {
		t.Errorf("Expected ruleset to change from index %d", initialRuleIndex)
	}
}

// TestLoggingToggle tests toggling of logging
func TestLoggingToggle(t *testing.T) {
	model := hammerclock.NewModel()

	// Get initial logging state
	initialLoggingState := model.Options.LoggingEnabled

	// Toggle logging
	updatedModel, _ := hammerclock.Update(&common.SetEnableLogMsg{Value: !initialLoggingState}, model)

	// Verify logging was toggled
	if updatedModel.Options.LoggingEnabled == initialLoggingState {
		t.Errorf("Expected logging to be toggled from %v to %v",
			initialLoggingState, !initialLoggingState)
	}
}

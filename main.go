package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"strings"
	"time"
)

// Color palettes
// K9s color palette
var (
	k9sBlue   = tcell.NewRGBColor(36, 96, 146)   // Dark blue for backgrounds
	k9sCyan   = tcell.NewRGBColor(0, 183, 235)   // Cyan for highlights
	k9sWhite  = tcell.NewRGBColor(255, 255, 255) // White for primary text
	k9sYellow = tcell.NewRGBColor(253, 185, 19)  // Yellow for warnings
	k9sGreen  = tcell.NewRGBColor(0, 200, 83)    // Green for success/active states
	k9sRed    = tcell.NewRGBColor(255, 0, 0)     // Red for errors/critical states
	k9sBlack  = tcell.NewRGBColor(0, 0, 0)       // Black for default backgrounds
)

// Dracula color palette
var (
	draculaBackground = tcell.NewRGBColor(40, 42, 54)    // Background
	draculaForeground = tcell.NewRGBColor(248, 248, 242) // Foreground
	draculaComment    = tcell.NewRGBColor(98, 114, 164)  // Comment
	draculaPink       = tcell.NewRGBColor(255, 121, 198) // Pink
	draculaPurple     = tcell.NewRGBColor(189, 147, 249) // Purple
	draculaGreen      = tcell.NewRGBColor(80, 250, 123)  // Green
	draculaYellow     = tcell.NewRGBColor(241, 250, 140) // Yellow
	draculaRed        = tcell.NewRGBColor(255, 85, 85)   // Red
	draculaCyan       = tcell.NewRGBColor(139, 233, 253) // Cyan
	draculaOrange     = tcell.NewRGBColor(255, 184, 108) // Orange
)

// Monokai color palette
var (
	monokaiBackground = tcell.NewRGBColor(39, 40, 34)    // Background
	monokaiForeground = tcell.NewRGBColor(248, 248, 242) // Foreground
	monokaiComment    = tcell.NewRGBColor(117, 113, 94)  // Comment
	monokaiRed        = tcell.NewRGBColor(249, 38, 114)  // Red
	monokaiOrange     = tcell.NewRGBColor(253, 151, 31)  // Orange
	monokaiYellow     = tcell.NewRGBColor(230, 219, 116) // Yellow
	monokaiGreen      = tcell.NewRGBColor(166, 226, 46)  // Green
	monokaiBlue       = tcell.NewRGBColor(102, 217, 239) // Blue
	monokaiPurple     = tcell.NewRGBColor(174, 129, 255) // Purple
	monokaiWhite      = tcell.NewRGBColor(255, 255, 255) // White
)

// Current color palette (default to K9s)
var (
	currentBlue     = k9sBlue
	currentCyan     = k9sCyan
	currentWhite    = k9sWhite
	currentDimWhite = tcell.NewRGBColor(180, 180, 180) // Dimmed white for inactive panels
	currentYellow   = k9sYellow
	currentGreen    = k9sGreen
	currentRed      = k9sRed
	currentBlack    = k9sBlack
)

type Player struct {
	name         string
	timeElapsed  time.Duration
	isTurn       bool
	currentPhase int
	armyList     []Unit
}

type Unit struct {
	Name   string
	Points int
}

// GameStatus represents the current state of the game
type GameStatus string

const (
	GameNotStarted GameStatus = "Game Not Started"
	GameInProgress GameStatus = "Game In Progress"
	GamePaused     GameStatus = "Game Paused"
)

type Settings struct {
	Name                 string   `json:"name"`
	Default              bool     `json:"default"`
	PlayerCount          int      `json:"playerCount"`
	PlayerNames          []string `json:"playerNames"`
	Phases               []string `json:"phases"`
	OneTurnForAllPlayers bool     `json:"oneTurnForAllPlayers"`
	ColorPalette         string   `json:"colorPalette"`
}

var defaultSettingsFilename = "defaultRules.json"

var defaultSettings = Settings{
	Name:                 "W40K 10th Edition",
	Default:              true,
	PlayerCount:          2,
	PlayerNames:          []string{"Player 1", "Player 2"},
	Phases:               []string{"Command Phase", "Movement Phase", "Shooting Phase", "Charge Phase", "Fight Phase", "End Phase"},
	OneTurnForAllPlayers: false,
	ColorPalette:         "k9s",
}

var phases []string

// MenuOption represents a menu option with a key and description
type MenuOption struct {
	Key         string
	Description string
	Action      func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool
}

// Track which screen is currently displayed
var currentScreen = "main" // Can be "main", "settings", or "about"

// Track whether the game has started
var gameStarted = false

// Current game status
var gameStatus = GameNotStarted

var menuOptions = []MenuOption{
	{Key: "2", Description: "Settings", Action: func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		// This will be implemented in the main function
		return false
	}},
	{Key: "A", Description: "About", Action: func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		// This will be implemented in the main function
		return false
	}},
}

var instructions = []MenuOption{
	{Key: "S", Description: "Start Game", Action: func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		// Toggle between start and pause
		if gameStatus == GamePaused {
			// Resume the game
			gameStatus = GameInProgress
			statusBar.SetText(string(gameStatus))
			return true // Return true to indicate the game should be started/resumed
		} else if gameStatus == GameInProgress {
			// Pause the game
			gameStatus = GamePaused
			statusBar.SetText(string(gameStatus))
			return false // Return false to indicate the game should remain in started state but paused
		} else {
			// Start the game if not already started
			gameStatus = GameInProgress
			statusBar.SetText(string(gameStatus))
			return true // Return true to indicate the game should be started
		}
	}},
	{Key: "SPACE", Description: "Switch Turns", Action: func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		// Switch turns
		for _, player := range players {
			player.isTurn = !player.isTurn
			if player.isTurn {
				player.currentPhase = 0
			}
		}
		updatePlayerPanels(players, playerPanels)

		// If we're not on the main screen, this is a good time to return to it
		if currentScreen != "main" {
			// Set the flag to return to main screen
			// The actual switching will be handled in the main function
			currentScreen = "main"
		}

		return false
	}},
	{Key: "P", Description: "Next Phase", Action: func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		// Move forward in the phase
		for _, player := range players {
			if player.isTurn {
				player.currentPhase = (player.currentPhase + 1) % len(phases)
			}
		}
		updatePlayerPanels(players, playerPanels)

		// If we're not on the main screen, this is a good time to return to it
		if currentScreen != "main" {
			// Set the flag to return to main screen
			// The actual switching will be handled in the main function
			currentScreen = "main"
		}

		return false
	}},
	{Key: "B", Description: "Previous Phase", Action: func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		// Move backward in the phase
		for _, player := range players {
			if player.isTurn {
				player.currentPhase = (player.currentPhase - 1 + len(phases)) % len(phases)
			}
		}
		updatePlayerPanels(players, playerPanels)

		// If we're not on the main screen, this is a good time to return to it
		if currentScreen != "main" {
			// Set the flag to return to main screen
			// The actual switching will be handled in the main function
			currentScreen = "main"
		}

		return false
	}},
	{Key: "Q", Description: "Quit", Action: func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		// This function will be updated with the done channel in the main function
		return false
	}},
}

// CLI usage information
var cliUsage = `
Hammerclock - A terminal-based timer and phase tracker for tabletop games

Usage:
  hammerclock [options]

Options:
  -s <file>    Specify a custom settings file (default: ` + defaultSettingsFilename + `)
  -h, --help   Show this help message

Examples:
  hammerclock                     # Run with default settings
  hammerclock -s rules/chess.json # Run with custom chess settings
`

func readSettingsFile(filename string) Settings {
	var settings Settings

	// Check if the settings file exists
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// If the requested file is not the default one, inform the user and use default
		if filename != defaultSettingsFilename {
			fmt.Printf("Settings file '%s' not found, using default settings file\n", filename)
			return readSettingsFile(defaultSettingsFilename)
		}

		// Default file doesn't exist, create it from defaultSettings
		fmt.Println("Default settings file not found, creating it")

		// Convert defaultSettings to JSON
		jsonData, err := json.MarshalIndent(defaultSettings, "", "  ")
		if err != nil {
			fmt.Println("Error creating default settings file:", err)
			return defaultSettings
		}

		// Write the JSON data to the file
		err = os.WriteFile(defaultSettingsFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing default settings file:", err)
			return defaultSettings
		}

		return defaultSettings
	} else if err != nil {
		// Some other error occurred
		fmt.Printf("Error checking settings file '%s': %v\n", filename, err)
		if filename != defaultSettingsFilename {
			fmt.Println("Falling back to default settings")
			return readSettingsFile(defaultSettingsFilename)
		}
		return defaultSettings
	}

	// File exists, read it
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading settings file '%s': %v\n", filename, err)
		if filename != defaultSettingsFilename {
			fmt.Println("Falling back to default settings")
			return readSettingsFile(defaultSettingsFilename)
		}
		return defaultSettings
	}

	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		fmt.Printf("Error processing settings file '%s': %v\n", filename, err)
		if filename != defaultSettingsFilename {
			fmt.Println("Falling back to default settings")
			return readSettingsFile(defaultSettingsFilename)
		}
		return defaultSettings
	}

	return settings
}

// applyColorPalette sets the current color variables based on the selected palette
func applyColorPalette(palette string) {
	switch palette {
	case "dracula":
		currentBlue = draculaPurple
		currentCyan = draculaCyan
		currentWhite = draculaForeground
		// Create a dimmed version of draculaForeground (reduce brightness by ~30%)
		currentDimWhite = tcell.NewRGBColor(174, 174, 169)
		currentYellow = draculaYellow
		currentGreen = draculaGreen
		currentRed = draculaRed
		currentBlack = draculaBackground
	case "monokai":
		currentBlue = monokaiBlue
		currentCyan = monokaiBlue
		currentWhite = monokaiForeground
		// Create a dimmed version of monokaiForeground (reduce brightness by ~30%)
		currentDimWhite = tcell.NewRGBColor(174, 174, 169)
		currentYellow = monokaiYellow
		currentGreen = monokaiGreen
		currentRed = monokaiRed
		currentBlack = monokaiBackground
	default: // "k9s" or any other value defaults to k9s
		currentBlue = k9sBlue
		currentCyan = k9sCyan
		currentWhite = k9sWhite
		// Create a dimmed version of k9sWhite (reduce brightness by ~30%)
		currentDimWhite = tcell.NewRGBColor(180, 180, 180)
		currentYellow = k9sYellow
		currentGreen = k9sGreen
		currentRed = k9sRed
		currentBlack = k9sBlack
	}
}

// createPlayerPanel creates a panel for a player using tview
func createPlayerPanel(player *Player, color string) *tview.Flex {
	panel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a box for the player name
	nameBox := tview.NewTextView().
		SetText("Player: " + player.name).
		SetTextColor(currentWhite)

	// Create a box for the time elapsed
	timeBox := tview.NewTextView().
		SetText(fmt.Sprintf("Time Elapsed: %v", player.timeElapsed)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(currentWhite)

	// Create a box for the current phase
	phaseBox := tview.NewTextView().
		SetText(fmt.Sprintf("Phase: %s", phases[player.currentPhase])).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(currentWhite)

	// Add the boxes to the panel
	panel.AddItem(nameBox, 1, 1, false).
		AddItem(tview.NewBox(), 1, 1, false). // Spacer
		AddItem(timeBox, 1, 1, false).
		AddItem(tview.NewBox(), 1, 1, false). // Spacer
		AddItem(phaseBox, 1, 1, false).
		AddItem(tview.NewBox(), 0, 1, false) // Flexible spacer at the bottom

	// Set the border and background
	panel.SetBorder(true)
	panel.SetBackgroundColor(currentBlack)

	// Title will be set by updatePlayerPanels based on turn and game state

	// Set the border color based on the player
	switch color {
	case "blue":
		panel.SetBorderColor(currentBlue)
	case "yellow":
		panel.SetBorderColor(currentYellow)
	case "green":
		panel.SetBorderColor(currentGreen)
	case "red":
		panel.SetBorderColor(currentRed)
	default:
		panel.SetBorderColor(currentBlack)
	}

	return panel
}

// createMenuBar creates a menu bar with the given options
func createMenuBar(options []MenuOption) *tview.TextView {
	menuText := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	var menuString strings.Builder
	for i, option := range options {
		key := option.Key
		item := option.Description

		if i > 0 {
			menuString.WriteString("  ")
		}
		menuString.WriteString("[white:black]" + key + "[:-] " + item)
	}

	menuText.SetText(menuString.String())
	return menuText
}

// updateMenuText updates the menu text based on the game state
func updateMenuText(menu *tview.TextView, status GameStatus) {
	var updatedInstructions []MenuOption

	// Copy the instructions and update the description for the "S" key
	for _, instruction := range instructions {
		if instruction.Key == "S" {
			// Update description based on game status
			if status == GameInProgress {
				instruction.Description = "Pause Game"
			} else if status == GamePaused {
				instruction.Description = "Resume Game"
			} else {
				instruction.Description = "Start Game"
			}
		}
		updatedInstructions = append(updatedInstructions, instruction)
	}

	// Rebuild the menu text
	var menuString strings.Builder
	for i, option := range updatedInstructions {
		key := option.Key
		item := option.Description

		if i > 0 {
			menuString.WriteString("  ")
		}
		menuString.WriteString("[white:black]" + key + "[:-] " + item)
	}

	menu.SetText(menuString.String())
}

// updatePlayerPanels updates the player panels with current information
func updatePlayerPanels(players []*Player, playerPanels []*tview.Flex) {
	for i, player := range players {
		// Get the text views
		nameBox := playerPanels[i].GetItem(0).(*tview.TextView)
		timeBox := playerPanels[i].GetItem(2).(*tview.TextView)
		phaseBox := playerPanels[i].GetItem(4).(*tview.TextView)

		// Update time elapsed
		timeBox.SetText(fmt.Sprintf("Time Elapsed: %v", player.timeElapsed))

		// Update current phase
		phaseBox.SetText(fmt.Sprintf("Phase: %s", phases[player.currentPhase]))

		// Update title and text color based on game state and turn
		if !gameStarted {
			// If game hasn't started, all panels have dimmed text
			playerPanels[i].SetTitle("")
			nameBox.SetTextColor(currentDimWhite)
			timeBox.SetTextColor(currentDimWhite)
			phaseBox.SetTextColor(currentDimWhite)
		} else if player.isTurn {
			// Game started and it's this player's turn
			playerPanels[i].SetTitle(" ACTIVE TURN ")
			// Use normal white for active player
			nameBox.SetTextColor(currentWhite)
			timeBox.SetTextColor(currentWhite)
			phaseBox.SetTextColor(currentWhite)
		} else {
			// Game started but it's not this player's turn
			playerPanels[i].SetTitle("")
			// Use dimmed white for inactive players
			nameBox.SetTextColor(currentDimWhite)
			timeBox.SetTextColor(currentDimWhite)
			phaseBox.SetTextColor(currentDimWhite)
		}
	}
}

// createSettingsScreen creates a screen that displays the current settings
func createSettingsScreen() *tview.Flex {
	settingsPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a title
	titleBox := tview.NewTextView().
		SetText("Settings").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(currentWhite)

	// Create content with current settings information
	contentBox := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(currentWhite).
		SetDynamicColors(true)

	// Build settings content
	var content strings.Builder
	content.WriteString("[::b]Current Game:[:-] " + defaultSettings.Name + "\n\n")
	content.WriteString("[::b]Player Count:[:-] " + fmt.Sprintf("%d", defaultSettings.PlayerCount) + "\n\n")
	content.WriteString("[::b]Players:[:-]\n")
	for i, name := range defaultSettings.PlayerNames {
		content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, name))
	}
	content.WriteString("\n")
	content.WriteString("[::b]Phases:[:-]\n")
	for i, phase := range defaultSettings.Phases {
		content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, phase))
	}
	content.WriteString("\n")
	content.WriteString("[::b]One Turn For All Players:[:-] " + fmt.Sprintf("%t", defaultSettings.OneTurnForAllPlayers) + "\n\n")
	content.WriteString("[::b]Color Palette:[:-] " + defaultSettings.ColorPalette + "\n\n")
	content.WriteString("\nPress [::b]2[:-] to return to the main screen")

	contentBox.SetText(content.String())

	// Add the boxes to the panel
	settingsPanel.AddItem(titleBox, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false). // Spacer
		AddItem(contentBox, 0, 1, false)

	// Set the border and background
	settingsPanel.SetBorder(true)
	settingsPanel.SetTitle(" Settings ")
	settingsPanel.SetBorderColor(currentCyan)
	settingsPanel.SetBackgroundColor(currentBlack)

	return settingsPanel
}

// createAboutScreen creates a screen that displays information about the application
func createAboutScreen() *tview.Flex {
	aboutPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a title
	titleBox := tview.NewTextView().
		SetText("About Hammerclock").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(currentWhite)

	// Create content with about information
	contentBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(currentWhite).
		SetDynamicColors(true)

	// Build about content
	var content strings.Builder
	content.WriteString("[::b]Hammerclock v1.0.0[:-]\n\n")
	content.WriteString("A terminal-based timer and phase tracker for tabletop games\n\n")
	content.WriteString("Created with [::b]Go[:-] and [::b]tview[:-]\n\n")
	content.WriteString("© 2023 Hammerclock Developers\n\n")
	content.WriteString("\nPress [::b]A[:-] to return to the main screen")

	contentBox.SetText(content.String())

	// Add the boxes to the panel
	aboutPanel.AddItem(titleBox, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false). // Spacer
		AddItem(contentBox, 0, 1, false)

	// Set the border and background
	aboutPanel.SetBorder(true)
	aboutPanel.SetTitle(" About ")
	aboutPanel.SetBorderColor(currentYellow)
	aboutPanel.SetBackgroundColor(currentBlack)

	return aboutPanel
}

// createStatusPanel creates a panel that displays the game status
func createStatusPanel(status string) *tview.Flex {
	statusPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create status text view
	statusTextView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(status)

	// Add the text view to the panel
	statusPanel.AddItem(statusTextView, 1, 0, false)

	// Set the border and background
	statusPanel.SetBorder(true)
	statusPanel.SetTitle(" Status ")
	statusPanel.SetBorderColor(currentCyan)
	statusPanel.SetBackgroundColor(currentBlack)

	return statusPanel
}

func main() {
	// Parse command line flags
	settingsFileFlag := flag.String("s", defaultSettingsFilename, "Path to the settings file")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, cliUsage)
	}
	flag.Parse()

	// Use the settings file specified by the flag
	loadedSettings := readSettingsFile(*settingsFileFlag)
	numberOfPlayers := loadedSettings.PlayerCount
	phases = loadedSettings.Phases

	// Apply the selected color palette
	applyColorPalette(loadedSettings.ColorPalette)

	// Apply current color palette to tview styles
	tview.Styles.PrimitiveBackgroundColor = currentBlack
	tview.Styles.ContrastBackgroundColor = currentGreen
	tview.Styles.MoreContrastBackgroundColor = currentCyan
	tview.Styles.BorderColor = currentCyan
	tview.Styles.TitleColor = currentWhite
	tview.Styles.GraphicsColor = currentWhite
	tview.Styles.PrimaryTextColor = currentWhite
	tview.Styles.SecondaryTextColor = currentYellow
	tview.Styles.TertiaryTextColor = currentGreen
	tview.Styles.InverseTextColor = currentRed
	tview.Styles.ContrastSecondaryTextColor = currentYellow

	// Create tview application
	app := tview.NewApplication()

	// Create players
	players := make([]*Player, numberOfPlayers)
	for i := 0; i < numberOfPlayers; i++ {
		playerName := fmt.Sprintf("Player %d", i+1) // Default name as fallback
		if i < len(loadedSettings.PlayerNames) {
			playerName = loadedSettings.PlayerNames[i]
		}
		players[i] = &Player{name: playerName, timeElapsed: 0, isTurn: i == 0, currentPhase: 0}
	}

	// Create the main layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a flex container for the top row
	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Add top menu to the left side
	topMenu := createMenuBar(menuOptions)
	topFlex.AddItem(topMenu, 0, 1, false)

	// Add a spacer for centering
	leftSpacer := tview.NewBox()
	topFlex.AddItem(leftSpacer, 0, 1, false)

	// Add name display in the middle
	nameDisplay := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	nameDisplay.SetText("[white:black]" + loadedSettings.Name + "[:-]")
	topFlex.AddItem(nameDisplay, 0, 1, false)

	// Add a spacer for centering
	rightSpacer := tview.NewBox()
	topFlex.AddItem(rightSpacer, 0, 1, false)

	// Add clock display to the right side
	clockDisplay := tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true)
	clockDisplay.SetText(time.Now().Format("15:04:05"))
	topFlex.AddItem(clockDisplay, 10, 0, false)

	// Add the top flex container to the main layout
	mainFlex.AddItem(topFlex, 1, 0, false)

	// Create player panels container
	playerPanelsContainer := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create player panels
	playerColors := []string{"blue", "yellow", "green", "red"}
	playerPanels := make([]*tview.Flex, numberOfPlayers)

	for i, player := range players {
		playerPanel := createPlayerPanel(player, playerColors[i%len(playerColors)])
		playerPanels[i] = playerPanel
		playerPanelsContainer.AddItem(playerPanel, 0, 1, false)
	}

	// Create settings and about screens (but don't add them yet)
	settingsScreen := createSettingsScreen()
	aboutScreen := createAboutScreen()

	// Add player panels to main layout
	mainFlex.AddItem(playerPanelsContainer, 0, 1, false)

	// Update menuOptions actions to handle screen switching
	menuOptions[0].Action = func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		app.QueueUpdateDraw(func() {
			// Toggle between main screen and settings screen
			if currentScreen == "settings" {
				currentScreen = "main"
				playerPanelsContainer.Clear()
				for _, panel := range playerPanels {
					playerPanelsContainer.AddItem(panel, 0, 1, false)
				}
			} else {
				currentScreen = "settings"
				playerPanelsContainer.Clear()
				playerPanelsContainer.AddItem(settingsScreen, 0, 1, false)
			}
		})
		return false
	}

	menuOptions[1].Action = func(app *tview.Application, players []*Player, playerPanels []*tview.Flex, statusBar *tview.TextView) bool {
		app.QueueUpdateDraw(func() {
			// Toggle between main screen and about screen
			if currentScreen == "about" {
				currentScreen = "main"
				playerPanelsContainer.Clear()
				for _, panel := range playerPanels {
					playerPanelsContainer.AddItem(panel, 0, 1, false)
				}
			} else {
				currentScreen = "about"
				playerPanelsContainer.Clear()
				playerPanelsContainer.AddItem(aboutScreen, 0, 1, false)
			}
		})
		return false
	}

	// Create status panel
	statusPanel := createStatusPanel(string(gameStatus))
	// Extract the status text view for later updates
	statusBar := statusPanel.GetItem(0).(*tview.TextView)
	mainFlex.AddItem(statusPanel, 1, 0, false)

	// Add bottom menu
	bottomMenu := createMenuBar(instructions)
	// Initialize menu text based on initial game status
	updateMenuText(bottomMenu, gameStatus)
	mainFlex.AddItem(bottomMenu, 1, 0, false)

	// Game state
	gameStarted := false
	gameStartedPtr := &gameStarted

	// Set up ticker for time tracking and clock display
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan struct{})
	doneIsClosed := false
	go func() {
		for {
			select {
			case <-ticker.C:
				// Update the clock display with the current time
				app.QueueUpdateDraw(func() {
					clockDisplay.SetText(time.Now().Format("15:04:05"))
				})

				if *gameStartedPtr {
					// Only increment time if the game is in progress (not paused)
					if gameStatus == GameInProgress {
						for _, player := range players {
							if player.isTurn {
								player.timeElapsed += 1 * time.Second
							}
						}
						app.QueueUpdateDraw(func() {
							updatePlayerPanels(players, playerPanels)
						})
					}
				}
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	// Set up key handling using tview's input capture
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			if !doneIsClosed {
				close(done)
				doneIsClosed = true
			}
			app.Stop()
			return nil
		case tcell.KeyRune:
			// Check if the key pressed matches any menu option
			keyRune := string(event.Rune())

			// Check menuOptions
			for _, option := range menuOptions {
				if strings.EqualFold(option.Key, keyRune) && option.Action != nil {
					startGame := option.Action(app, players, playerPanels, statusBar)
					if startGame {
						app.QueueUpdateDraw(func() {
							*gameStartedPtr = true
							gameStarted = true
							updatePlayerPanels(players, playerPanels)
						})
					}
					return nil
				}
			}

			// Check instructions
			for _, instruction := range instructions {
				if strings.EqualFold(instruction.Key, keyRune) && instruction.Action != nil {
					// If this is the "Q" key, close the done channel before executing the action
					if strings.EqualFold(instruction.Key, "Q") {
						if !doneIsClosed {
							close(done)
							doneIsClosed = true
						}
						app.Stop()
						return nil
					}

					// Execute the action
					startGame := instruction.Action(app, players, playerPanels, statusBar)

					// If this is the "S" key, update the menu text based on the new status
					// and update game state if needed, all in one QueueUpdateDraw call
					if strings.EqualFold(instruction.Key, "S") {
						app.QueueUpdateDraw(func() {
							updateMenuText(bottomMenu, gameStatus)
							if startGame {
								*gameStartedPtr = true
								gameStarted = true
								updatePlayerPanels(players, playerPanels)
							}
						})
					} else if startGame {
						// For other keys, update game state if needed
						*gameStartedPtr = true
						gameStarted = true
						app.QueueUpdateDraw(func() {
							updatePlayerPanels(players, playerPanels)
						})
					}

					// Check if we need to switch back to the main screen
					if currentScreen != "main" {
						app.QueueUpdateDraw(func() {
							// Switch back to the main screen
							currentScreen = "main"
							playerPanelsContainer.Clear()
							for _, panel := range playerPanels {
								playerPanelsContainer.AddItem(panel, 0, 1, false)
							}
						})
					}

					return nil
				}
			}

			if event.Rune() == ' ' {
				for _, instruction := range instructions {
					if instruction.Key == "SPACE" && instruction.Action != nil {
						// Execute the action
						startGame := instruction.Action(app, players, playerPanels, statusBar)

						// Update game state if needed
						if startGame {
							app.QueueUpdateDraw(func() {
								*gameStartedPtr = true
								gameStarted = true
								updatePlayerPanels(players, playerPanels)
							})
						}

						// Check if we need to switch back to the main screen
						if currentScreen != "main" {
							app.QueueUpdateDraw(func() {
								// Switch back to the main screen
								currentScreen = "main"
								playerPanelsContainer.Clear()
								for _, panel := range playerPanels {
									playerPanelsContainer.AddItem(panel, 0, 1, false)
								}
							})
						}

						return nil
					}
				}
			}
		default:
			// Just pass the event through
		}
		return event
	})

	// Start the application
	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
	}

	// Ensure the done channel is closed when the application exits
	if !doneIsClosed {
		close(done)
		doneIsClosed = true
	}
}

package app

import (
	"fmt"
	"hammerclock/internal/app/about"
	"hammerclock/internal/app/status"
	"strings"
	"time"

	"github.com/rivo/tview"
)

// View contains all the UI components
type View struct {
	App                   *tview.Application
	MainFlex              *tview.Flex
	PlayerPanelsContainer *tview.Flex
	PlayerPanels          []*tview.Flex
	TopMenu               *tview.TextView
	BottomMenu            *tview.TextView
	StatusPanel           *tview.Flex
	ClockDisplay          *tview.TextView
	OptionsScreen         *tview.Flex
	AboutScreen           *tview.Flex
}

// MenuOption represents a menu option with a key and description
type MenuOption struct {
	Key         string
	Description string
}

// ClockLayout determines the clock format string based on the model's time format setting (AMPM or 24-hour).
func ClockLayout(model *Model) string {
	// Determine the clock layout based on the options
	if model.Options.TimeFormat == "AMPM" {
		return "03:04:05 PM"
	}
	return "15:04:05"
}

// NewView creates a new view with the given model
func NewView(model *Model) *View {
	app := tview.NewApplication()

	// Apply the current color palette to tview styles
	applyTviewStyles(model.CurrentColorPalette)

	// Create the main layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a flex container for the top row
	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create top menu options
	menuOptions := []MenuOption{
		{Key: "O", Description: "options"},
		{Key: "A", Description: "About"},
	}

	// Add a top menu to the left side
	topMenu := createMenuBar(menuOptions).SetDynamicColors(true)
	topFlex.AddItem(topMenu, 0, 1, false)

	// Add a spacer for centering
	leftSpacer := tview.NewBox()
	topFlex.AddItem(leftSpacer, 0, 1, false)

	// Add name display in the middle
	nameDisplay := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	nameDisplay.SetText("[white]" + model.Options.Name + "[-]")
	topFlex.AddItem(nameDisplay, 0, 1, false)

	// Add a spacer for centering
	rightSpacer := tview.NewBox()
	topFlex.AddItem(rightSpacer, 0, 1, false)

	// Add a clock display to the right side
	clockDisplay := tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true)

	// Set the clock format based on the model's options
	var clockLayout = ClockLayout(model)

	clockDisplay.SetText(time.Now().Format(clockLayout))
	topFlex.AddItem(clockDisplay, 10, 0, false)

	// Add the top flex container to the main layout
	mainFlex.AddItem(topFlex, 1, 0, false)

	// Create player panels container
	playerPanelsContainer := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create player panels
	playerColors := []string{"blue", "yellow", "green", "red"}
	playerPanels := make([]*tview.Flex, len(model.Players))

	for i, player := range model.Players {
		playerPanel := createPlayerPanel(player, playerColors[i%len(playerColors)], model)
		playerPanels[i] = playerPanel
		playerPanelsContainer.AddItem(playerPanel, 0, 1, false)
	}

	// Create options and about screens
	optionsScreen := createOptionsScreen(model)
	aboutScreen := about.About(model.CurrentColorPalette.White)

	// Add player panels to the main layout
	mainFlex.AddItem(playerPanelsContainer, 0, 1, false)

	// Create status panel
	statusPanel := status.Status(string(model.GameStatus), model.CurrentColorPalette.Cyan, model.CurrentColorPalette.Black)

	// Add status panels to the main layout
	mainFlex.AddItem(statusPanel, 3, 0, false)

	// Create bottom menu options
	instructions := []MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

	// Add a bottom menu
	bottomMenu := createMenuBar(instructions).SetDynamicColors(true)
	// Initialize menu text based on the initial game status
	updateMenuText(bottomMenu, model.GameStatus)
	mainFlex.AddItem(bottomMenu, 1, 0, false)

	return &View{
		App:                   app,
		MainFlex:              mainFlex,
		PlayerPanelsContainer: playerPanelsContainer,
		PlayerPanels:          playerPanels,
		TopMenu:               topMenu,
		BottomMenu:            bottomMenu,
		StatusPanel:           statusPanel,
		ClockDisplay:          clockDisplay,
		OptionsScreen:         optionsScreen,
		AboutScreen:           aboutScreen,
	}
}

// Render updates the UI based on the current model
func (v *View) Render(model *Model) {
	// Update player panels
	updatePlayerPanels(model.Players, v.PlayerPanels, model)

	// Update status panel
	updateStatusPanel(v.StatusPanel, string(model.GameStatus), model)

	// Update menu text
	updateMenuText(v.BottomMenu, model.GameStatus)

	// Update the screen based on the current screen
	switch model.CurrentScreen {
	case "options":
		v.PlayerPanelsContainer.Clear()
		v.PlayerPanelsContainer.AddItem(v.OptionsScreen, 0, 1, false)
	case "about":
		v.PlayerPanelsContainer.Clear()
		v.PlayerPanelsContainer.AddItem(v.AboutScreen, 0, 1, false)
	default: // "main"
		v.PlayerPanelsContainer.Clear()
		for _, panel := range v.PlayerPanels {
			v.PlayerPanelsContainer.AddItem(panel, 0, 1, false)
		}
	}
}

// UpdateStatusPanel updates the status panel with the current game status
func updateStatusPanel(panel *tview.Flex, s string, model *Model) {
	statusTextView := panel.GetItem(0).(*tview.TextView)
	statusTextView.SetText(s)

	// Set the border color based on the game status
	switch model.GameStatus {
	case GameNotStarted:
		panel.SetBorderColor(model.CurrentColorPalette.Cyan)
	case GameInProgress:
		panel.SetBorderColor(model.CurrentColorPalette.Green)
	case GamePaused:
		panel.SetBorderColor(model.CurrentColorPalette.Yellow)
	}
}

// UpdateClock updates the clock display with the current time
func (v *View) UpdateClock(model *Model) {
	var clockLayout = ClockLayout(model)
	v.ClockDisplay.SetText(time.Now().Format(clockLayout))
}

// applyTviewStyles applies the color palette to tview styles
func applyTviewStyles(palette ColorPalette) {
	tview.Styles.PrimitiveBackgroundColor = palette.Black
	tview.Styles.ContrastBackgroundColor = palette.Green
	tview.Styles.MoreContrastBackgroundColor = palette.Cyan
	tview.Styles.BorderColor = palette.Cyan
	tview.Styles.TitleColor = palette.White
	tview.Styles.GraphicsColor = palette.White
	tview.Styles.PrimaryTextColor = palette.White
	tview.Styles.SecondaryTextColor = palette.Yellow
	tview.Styles.TertiaryTextColor = palette.Green
	tview.Styles.InverseTextColor = palette.Red
	tview.Styles.ContrastSecondaryTextColor = palette.Yellow
}

// createPlayerPanel creates a panel for a player using tview
func createPlayerPanel(player *Player, color string, model *Model) *tview.Flex {
	panel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a box for the player name
	nameBox := tview.NewTextView().
		SetText("Player: " + player.Name).
		SetTextColor(model.CurrentColorPalette.White)

	// Create a box for the time elapsed
	timeBox := tview.NewTextView().
		SetText(fmt.Sprintf("Time Elapsed: %v", player.TimeElapsed)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)

	// Create a box for the current phase
	phaseBox := tview.NewTextView().
		SetText(fmt.Sprintf("Phase: %s", model.Phases[player.CurrentPhase])).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)

	// Add the boxes to the panel
	panel.AddItem(nameBox, 1, 1, false).
		AddItem(tview.NewBox(), 1, 1, false). // Spacer
		AddItem(timeBox, 1, 1, false).
		AddItem(tview.NewBox(), 1, 1, false). // Spacer
		AddItem(phaseBox, 1, 1, false).
		AddItem(tview.NewBox(), 0, 1, false) // Flexible spacer at the bottom

	// Set the border and background
	panel.SetBorder(true)
	panel.SetBackgroundColor(model.CurrentColorPalette.Black)

	// Set the border color based on the player
	switch color {
	case "blue":
		panel.SetBorderColor(model.CurrentColorPalette.Blue)
	case "yellow":
		panel.SetBorderColor(model.CurrentColorPalette.Yellow)
	case "green":
		panel.SetBorderColor(model.CurrentColorPalette.Green)
	case "red":
		panel.SetBorderColor(model.CurrentColorPalette.Red)
	default:
		panel.SetBorderColor(model.CurrentColorPalette.Black)
	}

	return panel
}

// createMenuBar creates a menu bar with the given options
func createMenuBar(options []MenuOption) *tview.TextView {
	menuText := tview.NewTextView()
	var menuString strings.Builder

	for i, option := range options {
		if i > 0 {
			menuString.WriteString("   ")
		}
		var menuItem = formatMenuOption(option)
		menuString.WriteString(menuItem)
	}

	menuText.SetText(menuString.String())
	return menuText
}

// formatMenuOption formats a single menu option for display in the menu bar.
func formatMenuOption(option MenuOption) string {
	menuItem := fmt.Sprintf("[white]%s[d:] %s", option.Key, option.Description)
	return menuItem
}

// updateMenuText updates the menu text based on the game state
func updateMenuText(menu *tview.TextView, status GameStatus) {
	// Define the instructions
	instructions := []MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

	// Update the description for the "S" key based on the game status
	for i := range instructions {
		if instructions[i].Key == "S" {
			switch status {
			case GameInProgress:
				instructions[i].Description = "Pause Game"
			case GamePaused:
				instructions[i].Description = "Resume Game"
			}
		}
	}

	// Build the menu text
	var menuString strings.Builder
	for i, option := range instructions {
		if i > 0 {
			menuString.WriteString("   ")
		}
		menuString.WriteString("[white]" + option.Key + "[d:] " + option.Description)
	}

	menu.SetText(menuString.String())
}

// updatePlayerPanels updates the player panels with current information
func updatePlayerPanels(players []*Player, playerPanels []*tview.Flex, model *Model) {
	for i, player := range players {
		// Get the text views
		nameBox := playerPanels[i].GetItem(0).(*tview.TextView)
		timeBox := playerPanels[i].GetItem(2).(*tview.TextView)
		phaseBox := playerPanels[i].GetItem(4).(*tview.TextView)

		// Update time elapsed
		timeBox.SetText(fmt.Sprintf("Time Elapsed: %v", player.TimeElapsed))

		// Update current phase
		phaseBox.SetText(fmt.Sprintf("Phase: %s", model.Phases[player.CurrentPhase]))

		// Update title and text color based on game state and turn
		if !model.GameStarted {
			// If the game hasn't started, all panels have dimmed text
			playerPanels[i].SetTitle("")
			nameBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			timeBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			phaseBox.SetTextColor(model.CurrentColorPalette.DimWhite)
		} else if player.IsTurn {
			// Game started and it's this player's turn
			playerPanels[i].SetTitle(" ACTIVE TURN ")
			// Use normal white for active player
			nameBox.SetTextColor(model.CurrentColorPalette.White)
			timeBox.SetTextColor(model.CurrentColorPalette.White)
			phaseBox.SetTextColor(model.CurrentColorPalette.White)
		} else {
			// Game started, but it's not this player's turn
			playerPanels[i].SetTitle("")
			// Use dimmed white for inactive players
			nameBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			timeBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			phaseBox.SetTextColor(model.CurrentColorPalette.DimWhite)
		}
	}
}

// createOptionsScreen creates a screen that displays the current options
func createOptionsScreen(model *Model) *tview.Flex {
	optionsPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a title
	titleBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)

	// Create content with current options information
	contentBox := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(model.CurrentColorPalette.White).
		SetDynamicColors(true)

	// Build options content
	var content strings.Builder
	content.WriteString(" [b]Name of the ruleset:[-b] " + model.Options.Name + "\n\n")
	content.WriteString(" [b]Player Count:[-b] " + fmt.Sprintf("%d", model.Options.PlayerCount) + "\n\n")
	content.WriteString(" [b]Players:[-b]\n")
	for i, name := range model.Options.PlayerNames {
		content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, name))
	}
	content.WriteString("\n")
	content.WriteString(" [b]Phases:[-b]\n")
	for i, phase := range model.Options.Phases {
		content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, phase))
	}
	content.WriteString("\n")
	content.WriteString(" [b]One Turn For All Players:[-b] " + fmt.Sprintf("%t", model.Options.OneTurnForAllPlayers) + "\n\n")
	content.WriteString(" [b]Color Palette:[-b] " + model.Options.ColorPalette + "\n\n")
	content.WriteString(" [b]Time Format:[-b] " + model.Options.TimeFormat + "\n\n")
	content.WriteString("\n Press [white]O[:d] to return to the main screen")

	contentBox.SetText(content.String())

	// Add the boxes to the panel
	optionsPanel.AddItem(titleBox, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false). // Spacer
		AddItem(contentBox, 0, 1, false)

	// Set the border and background
	optionsPanel.SetBorder(true)
	optionsPanel.SetTitle(" options ")
	optionsPanel.SetBorderColor(model.CurrentColorPalette.Cyan)
	optionsPanel.SetBackgroundColor(model.CurrentColorPalette.Black)

	return optionsPanel
}

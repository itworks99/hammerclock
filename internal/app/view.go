package app

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mbndr/figlet4go"
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
	StatusBar             *tview.TextView
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
		{Key: "O", Description: "Options"},
		{Key: "A", Description: "About"},
	}

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
	nameDisplay.SetText("[white:black]" + model.Options.Name + "[:-]")
	topFlex.AddItem(nameDisplay, 0, 1, false)

	// Add a spacer for centering
	rightSpacer := tview.NewBox()
	topFlex.AddItem(rightSpacer, 0, 1, false)

	// Add clock display to the right side
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
	aboutScreen := createAboutScreen(model)

	// Add player panels to main layout
	mainFlex.AddItem(playerPanelsContainer, 0, 1, false)

	// Create status panel
	statusPanel := createStatusPanel(string(model.GameStatus), model)
	// Extract the status text view for later updates
	statusBar := statusPanel.GetItem(0).(*tview.TextView)
	mainFlex.AddItem(statusPanel, 1, 0, false)

	// Create bottom menu options
	instructions := []MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

	// Add bottom menu
	bottomMenu := createMenuBar(instructions)
	// Initialize menu text based on initial game status
	updateMenuText(bottomMenu, model.GameStatus)
	mainFlex.AddItem(bottomMenu, 1, 0, false)

	return &View{
		App:                   app,
		MainFlex:              mainFlex,
		PlayerPanelsContainer: playerPanelsContainer,
		PlayerPanels:          playerPanels,
		TopMenu:               topMenu,
		BottomMenu:            bottomMenu,
		StatusBar:             statusBar,
		ClockDisplay:          clockDisplay,
		OptionsScreen:         optionsScreen,
		AboutScreen:           aboutScreen,
	}
}

// Render updates the UI based on the current model
func (v *View) Render(model *Model) {
	// Update player panels
	updatePlayerPanels(model.Players, v.PlayerPanels, model)

	// Update status bar
	v.StatusBar.SetText(string(model.GameStatus))

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
	menuText := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	var menuString strings.Builder
	for i, option := range options {
		key := option.Key
		item := option.Description

		if i > 0 {
			menuString.WriteString("   ")
		}
		menuString.WriteString("[white:black]" + key + "[:-] " + item)
	}

	menuText.SetText(menuString.String())
	return menuText
}

// updateMenuText updates the menu text based on the game state
func updateMenuText(menu *tview.TextView, status GameStatus) {
	var updatedInstructions []MenuOption

	// Define the instructions
	instructions := []MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

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
			menuString.WriteString("   ")
		}
		menuString.WriteString("[white:black]" + key + "[:-] " + item)
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
			// If game hasn't started, all panels have dimmed text
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
			// Game started but it's not this player's turn
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
		SetText("Options").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)

	// Create content with current options information
	contentBox := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(model.CurrentColorPalette.White).
		SetDynamicColors(true)

	// Build options content
	var content strings.Builder
	content.WriteString("[::b]Current Game:[:-] " + model.Options.Name + "\n\n")
	content.WriteString("[::b]Player Count:[:-] " + fmt.Sprintf("%d", model.Options.PlayerCount) + "\n\n")
	content.WriteString("[::b]Players:[:-]\n")
	for i, name := range model.Options.PlayerNames {
		content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, name))
	}
	content.WriteString("\n")
	content.WriteString("[::b]Phases:[:-]\n")
	for i, phase := range model.Options.Phases {
		content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, phase))
	}
	content.WriteString("\n")
	content.WriteString("[::b]One Turn For All Players:[:-] " + fmt.Sprintf("%t", model.Options.OneTurnForAllPlayers) + "\n\n")
	content.WriteString("[::b]Color Palette:[:-] " + model.Options.ColorPalette + "\n\n")
	content.WriteString("\nPress [::b]2[:-] to return to the main screen")

	contentBox.SetText(content.String())

	// Add the boxes to the panel
	optionsPanel.AddItem(titleBox, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false). // Spacer
		AddItem(contentBox, 0, 1, false)

	// Set the border and background
	optionsPanel.SetBorder(true)
	optionsPanel.SetTitle(" Options ")
	optionsPanel.SetBorderColor(model.CurrentColorPalette.Cyan)
	optionsPanel.SetBackgroundColor(model.CurrentColorPalette.Black)

	return optionsPanel
}

// createAboutScreen creates a screen that displays information about the application
func createAboutScreen(model *Model) *tview.Flex {
	aboutPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create content with about information
	contentBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White).
		SetDynamicColors(true)

	// Use figlet4go to create a fancy title
	ascii := figlet4go.NewAsciiRender()
	// Set the font to "slant" for a slanted style
	options := figlet4go.NewRenderOptions()
	options.FontName = "larry3d"
	renderStr, _ := ascii.RenderOpts("Hammerclock", options)

	// Build about content
	var content strings.Builder
	content.WriteString(renderStr)
	content.WriteString("[::b]v0.1[:-]\n\n\n")
	content.WriteString("A terminal-based timer and phase tracker for tabletop games\n\n\n")
	content.WriteString("\nPress [::b]A[:-] to return to the main screen")

	contentBox.SetText(content.String())

	// Add the boxes to the panel
	aboutPanel.AddItem(tview.NewBox(), 1, 0, false). // Spacer
								AddItem(contentBox, 0, 1, false)

	// Set the border and background
	aboutPanel.SetBorder(true)
	aboutPanel.SetTitle(" About ")
	aboutPanel.SetBorderColor(tcell.ColorYellow)
	aboutPanel.SetBackgroundColor(tcell.ColorBlack)

	return aboutPanel
}

// createStatusPanel creates a panel that displays the game status
func createStatusPanel(status string, model *Model) *tview.Flex {
	statusPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create status text view
	statusTextView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(status)

	// Add the text view to the panel
	statusPanel.AddItem(statusTextView, 3, 0, false)

	// Set the border and background
	statusPanel.SetBorder(true)
	statusPanel.SetTitle(" Status ")
	statusPanel.SetBorderColor(model.CurrentColorPalette.Cyan)
	statusPanel.SetBackgroundColor(model.CurrentColorPalette.Black)

	return statusPanel
}

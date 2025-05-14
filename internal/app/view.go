package app

import (
	"strings"
	"time"

	"hammerclock/components/hammerclock/Palette"
	"hammerclock/internal/app/AboutPanel"
	"hammerclock/internal/app/MenuBar"
	"hammerclock/internal/app/StatusPanel"
	"hammerclock/internal/app/clock"

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
	OptionsScreen         *tview.Grid
	AboutScreen           *tview.Flex
	// Message channel for sending messages to the update function
	MessageChan chan<- Message
}

// NewView creates a new view with the given model
func NewView(model *Model, msgChan chan<- Message) *View {
	app := tview.NewApplication()

	// Apply the current color palette to tview styles
	Palette.ApplyColorPalette(model.CurrentColorPalette)

	// Create the main layout
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create a flex container for the top row
	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create top menu options
	menuOptions := []MenuBar.MenuOption{
		{Key: "O", Description: "Options"},
		{Key: "A", Description: "About"},
	}

	// Add a top menu to the left side
	topMenu := MenuBar.CreateMenuBar(menuOptions).SetDynamicColors(true)
	topFlex.AddItem(topMenu, 0, 1, false)

	// Add a spacer for centering
	leftSpacer := tview.NewBox()
	topFlex.AddItem(leftSpacer, 0, 1, false)

	// Add name display in the middle
	nameDisplay := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	nameDisplay.SetText("[white]" + model.Options.Rules[model.Options.Default].Name + "[-]")
	topFlex.AddItem(nameDisplay, 0, 1, false)

	// Add a spacer for centering
	rightSpacer := tview.NewBox()
	topFlex.AddItem(rightSpacer, 0, 1, false)

	// Add a clock display to the right side
	hClock := clock.Display(model.Options.TimeFormat, model.CurrentColorPalette.White)
	topFlex.AddItem(hClock, 10, 0, false)

	// Add the top flex container to the main layout
	mainFlex.AddItem(topFlex, 1, 0, false)

	// Create player panels container
	playerPanelsContainer := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create player panels
	playerColors := []string{"blue", "yellow", "green", "red"}
	playerPanels := make([]*tview.Flex, len(model.Players))

	for i, player := range model.Players {
		playerPanel := CreatePlayerPanel(player, playerColors[i%len(playerColors)], model)
		playerPanels[i] = playerPanel
		playerPanelsContainer.AddItem(playerPanel, 0, 1, false)
	}

	// Create options and about screens
	optionsScreen := createOptionsScreen(model, msgChan)
	aboutScreen := AboutPanel.Create(model.CurrentColorPalette.White)

	// Add player panels to the main layout
	mainFlex.AddItem(playerPanelsContainer, 0, 1, false)

	// Create statusbar panel
	statusPanel := StatusPanel.Create(string(model.GameStatus), model.CurrentColorPalette.Cyan, model.CurrentColorPalette.Black)

	// Add statusbar panels to the main layout
	mainFlex.AddItem(statusPanel, 3, 0, false)

	// Create bottom menu options
	instructions := []MenuBar.MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

	// Add a bottom menu
	bottomMenu := MenuBar.CreateMenuBar(instructions).SetDynamicColors(true)
	// Initialize menu text based on the initial game statusbar
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
		ClockDisplay:          hClock,
		OptionsScreen:         optionsScreen,
		AboutScreen:           aboutScreen,
		MessageChan:           msgChan,
	}
}

// Render updates the UI based on the current model
func (v *View) Render(model *Model) {
	// Update player panels
	updatePlayerPanels(model.Players, v.PlayerPanels, model)

	// Update statusbar panel
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

// updateStatusPanel updates the statusbar panel with the current game statusbar
func updateStatusPanel(panel *tview.Flex, s string, model *Model) {
	statusTextView := panel.GetItem(0).(*tview.TextView)
	statusTextView.SetText(s)

	// Set the border color based on the game statusbar
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
	var clockLayout = clock.TimeFormat(model.Options.TimeFormat)
	v.ClockDisplay.SetText(time.Now().Format(clockLayout))
}

// updateMenuText updates the menu text based on the game state
func updateMenuText(menu *tview.TextView, status GameStatus) {
	// Define the instructions
	instructions := []MenuBar.MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

	// Update the description for the "S" key based on the game statusbar
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

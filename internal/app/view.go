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

// View represents the main UI structure of the application.
type View struct {
	App                   *tview.Application // The main tview application instance.
	MainFlex              *tview.Flex        // The main container for the UI layout.
	PlayerPanelsContainer *tview.Flex        // Container for player panels.
	PlayerPanels          []*tview.Flex      // List of individual player panels.
	TopMenu               *tview.TextView    // The top menu bar.
	BottomMenu            *tview.TextView    // The bottom menu bar.
	StatusPanel           *tview.Flex        // Panel displaying the current game status.
	ClockDisplay          *tview.TextView    // Text view for displaying the clock.
	OptionsScreen         *tview.Grid        // Grid layout for the options screen.
	AboutScreen           *tview.Flex        // Flex layout for the about screen.
	MessageChan           chan<- Message     // Channel for sending messages to the application.
	CurrentScreen         string             // Tracks the currently displayed screen.
}

// NewView initializes and returns a new View instance.
// It sets up the main UI components and applies the current color palette.
func NewView(model *Model, msgChan chan<- Message) *View {
	app := tview.NewApplication()
	Palette.ApplyColorPalette(model.CurrentColorPalette)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	topFlex := createTopFlex(model)
	mainFlex.AddItem(topFlex, 1, 0, false)

	playerPanelsContainer, playerPanels := createPlayerPanels(model)
	mainFlex.AddItem(playerPanelsContainer, 0, 1, false)

	optionsScreen := createOptionsScreen(model, msgChan)
	aboutScreen := AboutPanel.Create(model.CurrentColorPalette.White)

	statusPanel := StatusPanel.Create(string(model.GameStatus), model.CurrentColorPalette.Cyan, model.CurrentColorPalette.Black)
	mainFlex.AddItem(statusPanel, 3, 0, false)

	bottomMenu := createBottomMenu(model.GameStatus)
	mainFlex.AddItem(bottomMenu, 1, 0, false)

	return &View{
		App:                   app,
		MainFlex:              mainFlex,
		PlayerPanelsContainer: playerPanelsContainer,
		PlayerPanels:          playerPanels,
		TopMenu:               topFlex.GetItem(0).(*tview.TextView),
		BottomMenu:            bottomMenu,
		StatusPanel:           statusPanel,
		ClockDisplay:          topFlex.GetItem(4).(*tview.TextView),
		OptionsScreen:         optionsScreen,
		AboutScreen:           aboutScreen,
		MessageChan:           msgChan,
		CurrentScreen:         "", // Initialize with an empty screen.
	}
}

// Render updates the UI based on the current model state.
// It refreshes player panels, status panel, and menu text, and switches screens as needed.
func (v *View) Render(model *Model) {
	if model.CurrentScreen != v.CurrentScreen {
		v.CurrentScreen = model.CurrentScreen
		v.PlayerPanelsContainer.Clear()
		switch model.CurrentScreen {
		case "options":
			v.PlayerPanelsContainer.AddItem(v.OptionsScreen, 0, 1, false)
		case "about":
			v.PlayerPanelsContainer.AddItem(v.AboutScreen, 0, 1, false)
		default:
			for _, panel := range v.PlayerPanels {
				v.PlayerPanelsContainer.AddItem(panel, 0, 1, false)
			}
		}
	}

	updatePlayerPanels(model.Players, v.PlayerPanels, model)
	updateStatusPanel(v.StatusPanel, string(model.GameStatus), model)
	updateMenuText(v.BottomMenu, model.GameStatus)
}

// UpdateClock updates the clock display with the current time.
// The time format is determined by the model's options.
func (v *View) UpdateClock(model *Model) {
	currentTime := time.Now().Format(clock.TimeFormat(model.Options.TimeFormat))
	if v.ClockDisplay.GetText(false) != currentTime {
		v.ClockDisplay.SetText(currentTime)
	}
}

// updateStatusPanel updates the status panel with the current game status.
// It also changes the border color based on the game status.
func updateStatusPanel(panel *tview.Flex, status string, model *Model) {
	StatusPanel.UpdateWithGameTime(panel, status, model.TotalGameTime)

	switch model.GameStatus {
	case GameNotStarted:
		panel.SetBorderColor(model.CurrentColorPalette.Cyan)
	case GameInProgress:
		panel.SetBorderColor(model.CurrentColorPalette.Green)
	case GamePaused:
		panel.SetBorderColor(model.CurrentColorPalette.Yellow)
	}
}

// updateMenuText updates the bottom menu text based on the current game status.
// It modifies the description of menu options dynamically.
func updateMenuText(menu *tview.TextView, status GameStatus) {
	instructions := []MenuBar.MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

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

	var menuString strings.Builder
	for i, option := range instructions {
		if i > 0 {
			menuString.WriteString("   ")
		}
		menuString.WriteString("[white]" + option.Key + "[d:] " + option.Description)
	}
	menu.SetText(menuString.String())
}

// createTopFlex creates the top flex layout containing the menu, name display, and clock.
func createTopFlex(model *Model) *tview.Flex {
	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	topMenu := MenuBar.CreateMenuBar([]MenuBar.MenuOption{
		{Key: "O", Description: "Options"},
		{Key: "A", Description: "About"},
	}).SetDynamicColors(true)
	topFlex.AddItem(topMenu, 0, 1, false)

	topFlex.AddItem(tview.NewBox(), 0, 1, false)

	nameDisplay := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("[white]" + model.Options.Rules[model.Options.Default].Name + "[-]")
	topFlex.AddItem(nameDisplay, 0, 1, false)

	topFlex.AddItem(tview.NewBox(), 0, 1, false)

	hClock := clock.Display(model.Options.TimeFormat, model.CurrentColorPalette.White)
	topFlex.AddItem(hClock, 10, 0, false)

	return topFlex
}

// createPlayerPanels creates the player panels and their container.
// Each panel is assigned a color from a predefined list.
func createPlayerPanels(model *Model) (*tview.Flex, []*tview.Flex) {
	container := tview.NewFlex().SetDirection(tview.FlexColumn)
	playerPanels := make([]*tview.Flex, len(model.Players))
	colors := []string{"blue", "yellow", "green", "red"}

	for i, player := range model.Players {
		panel := CreatePlayerPanel(player, colors[i%len(colors)], model)
		playerPanels[i] = panel
		container.AddItem(panel, 0, 1, false)
	}
	return container, playerPanels
}

// createBottomMenu creates the bottom menu bar and initializes its text.
func createBottomMenu(status GameStatus) *tview.TextView {
	menu := MenuBar.CreateMenuBar(nil).SetDynamicColors(true)
	updateMenuText(menu, status)
	return menu
}

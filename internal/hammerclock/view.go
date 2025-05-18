package hammerclock

import (
	"strings"
	"time"

	"github.com/rivo/tview"
	"hammerclock/internal/hammerclock/common"
	"hammerclock/internal/hammerclock/palette"
	"hammerclock/internal/hammerclock/ui"
)

// View represents the main UI structure of the application.
type View struct {
	App                   *tview.Application    // The main tview application instance.
	MainFlex              *tview.Flex           // The main container for the UI layout.
	PlayerPanelsContainer *tview.Flex           // Container for player panels.
	PlayerPanels          []*tview.Flex         // List of individual player panels.
	TopMenu               *tview.TextView       // The top menu bar.
	BottomMenu            *tview.TextView       // The bottom menu bar.
	StatusPanel           *tview.Flex           // Panel displaying the current game status.
	ClockDisplay          *tview.TextView       // Text view for displaying the clock.
	OptionsScreen         *tview.Grid           // Grid layout for the options screen.
	AboutScreen           *tview.Flex           // Flex layout for the about screen.
	MessageChan           chan<- common.Message // Channel for sending messages to the application.
	CurrentScreen         string                // Tracks the currently displayed screen.
}

// NewView initializes and returns a new View instance.
// It sets up the main UI components and applies the current color palette.
func NewView(model *common.Model, msgChan chan<- common.Message) *View {
	app := tview.NewApplication()
	palette.ApplyColorPalette(model.CurrentColorPalette)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	topFlex := createTopFlex(model)
	mainFlex.AddItem(topFlex, 1, 0, false)

	playerPanelsContainer, playerPanels := createPlayerPanels(model)
	mainFlex.AddItem(playerPanelsContainer, 0, 1, false)

	optionsScreen := ui.CreateOptionsScreen(model, msgChan)
	aboutScreen := ui.CreateAboutPanel(model.CurrentColorPalette.White)

	statusPanel := ui.CreateStatusPanel(string(model.GameStatus), model.CurrentColorPalette.Cyan, model.CurrentColorPalette.Black)
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
func (v *View) Render(model *common.Model) {
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

	ui.UpdatePlayerPanels(model.Players, v.PlayerPanels, model)
	updateStatusPanel(v.StatusPanel, string(model.GameStatus), model)
	updateMenuText(v.BottomMenu, model.GameStatus)
}

// UpdateClock updates the clock display with the current time.
// The time format is determined by the model's options.
func (v *View) UpdateClock(model *common.Model) {
	currentTime := time.Now().Format(ui.TimeFormat(model.Options.TimeFormat))
	if v.ClockDisplay.GetText(false) != currentTime {
		v.ClockDisplay.SetText(currentTime)
	}
}

// updateStatusPanel updates the status panel with the current game status.
// It also changes the border color based on the game status.
func updateStatusPanel(panel *tview.Flex, status string, model *common.Model) {
	ui.UpdateWithGameTime(panel, status, model.TotalGameTime)

	switch model.GameStatus {
	case gameNotStarted:
		panel.SetBorderColor(model.CurrentColorPalette.Cyan)
	case gameInProgress:
		panel.SetBorderColor(model.CurrentColorPalette.Green)
	case gamePaused:
		panel.SetBorderColor(model.CurrentColorPalette.Yellow)
	}
}

// updateMenuText updates the bottom menu text based on the current game status.
// It modifies the description of menu options dynamically.
func updateMenuText(menu *tview.TextView, status common.GameStatus) {
	instructions := []ui.MenuOption{
		{Key: "S", Description: "Start Game"},
		{Key: "E", Description: "End Game"},
		{Key: "SPACE", Description: "Switch Turns"},
		{Key: "P", Description: "Next Phase"},
		{Key: "B", Description: "Previous Phase"},
		{Key: "Q", Description: "Quit"},
	}

	for i := range instructions {
		if instructions[i].Key == "S" {
			switch status {
			case gameInProgress:
				instructions[i].Description = "Pause Game"
			case gamePaused:
				instructions[i].Description = "Resume Game"
			}
		}
	}

	var menuString strings.Builder
	for i, option := range instructions {
		if i > 0 {
			menuString.WriteString("   ")
		}

		// Special case for End Game option - dimmed and only visible when game started
		if option.Key == "E" {
			if status == gameNotStarted {
				// Skip the End Game option when game hasn't started
				continue
			}
			// Show dimmed when game is started
			menuString.WriteString("[#888888]" + option.Key + "[d:] " + option.Description)
		} else {
			menuString.WriteString("[white]" + option.Key + "[d:] " + option.Description)
		}
	}
	menu.SetText(menuString.String())
}

// createTopFlex creates the top flex layout containing the menu, name display, and clock.
func createTopFlex(model *common.Model) *tview.Flex {
	topFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	topMenu := ui.CreateMenuBar([]ui.MenuOption{
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

	hClock := ui.Display(model.Options.TimeFormat, model.CurrentColorPalette.White)
	topFlex.AddItem(hClock, 10, 0, false)

	return topFlex
}

// createPlayerPanels creates the player panels and their container.
// Each panel is assigned a color from a predefined list.
func createPlayerPanels(model *common.Model) (*tview.Flex, []*tview.Flex) {
	container := tview.NewFlex().SetDirection(tview.FlexColumn)
	playerPanels := make([]*tview.Flex, len(model.Players))
	colors := []string{"blue", "yellow", "green", "red"}

	for i, player := range model.Players {
		panel := ui.CreatePlayerPanel(player, colors[i%len(colors)], model)
		playerPanels[i] = panel
		container.AddItem(panel, 0, 1, false)
	}
	return container, playerPanels
}

// createBottomMenu creates the bottom menu bar and initializes its text.
func createBottomMenu(status common.GameStatus) *tview.TextView {
	menu := ui.CreateMenuBar(nil).SetDynamicColors(true)
	updateMenuText(menu, status)
	return menu
}

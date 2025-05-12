package app

import (
	"fmt"
	"hammerclock/internal/app/about"
	"hammerclock/internal/app/status"
	"strconv"
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
	OptionsScreen         *tview.Grid
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
	nameDisplay.SetText("[white]" + model.Options.Rules[model.Options.Default].Name + "[-]")
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

	phaseBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)
	// Create a box for the current phase
	if model.Options.Rules[model.Options.Default].OneTurnForAllPlayers != true {
		phaseBox = tview.NewTextView().
			SetText(fmt.Sprintf("Phase: %s", model.Phases[player.CurrentPhase])).
			SetTextAlign(tview.AlignCenter).
			SetTextColor(model.CurrentColorPalette.White)
	}
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
		if model.Options.Rules[model.Options.Default].OneTurnForAllPlayers != true {
			phaseBox.SetText(fmt.Sprintf("Phase: %s", model.Phases[player.CurrentPhase]))
		} else {
			phaseBox.SetText("")
		}

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
func createOptionsScreen(model *Model) *tview.Grid {
	optionsPanel := tview.NewGrid().
		SetRows(0, 0).
		SetColumns(0, 0).
		SetBorders(true)

	optionsBox := tview.NewFlex().
		SetDirection(tview.FlexRow)

	HelpContentBox := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(model.CurrentColorPalette.White).
		SetDynamicColors(true)

	currentRulesetContentBox := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(model.CurrentColorPalette.White).
		SetDynamicColors(true)

	var rulesetsNames []string
	for _, ruleset := range model.Options.Rules {
		rulesetsNames = append(rulesetsNames, ruleset.Name)
	}

	rulesetBox := tview.NewDropDown().
		SetLabel("Select ruleset(press Enter): ").
		SetOptions(rulesetsNames, func(option string, index int) {
			model.Options.Default = index
			model.Phases = model.Options.Rules[index].Phases
		}).SetCurrentOption(model.Options.Default).
		SetLabelColor(model.CurrentColorPalette.White)

	playerCountBox := tview.NewInputField().
		SetLabel("Player Count: ").
		SetText(fmt.Sprintf("%d", model.Options.PlayerCount)).
		SetLabelColor(model.CurrentColorPalette.White).
		SetFieldWidth(3).
		SetChangedFunc(func(text string) {
			count, err := strconv.Atoi(text)
			if err == nil && count > 0 {
				model.Options.PlayerCount = count
			}
		})

	playerNamesBox := tview.NewInputField().
		SetLabel("Player Names: ").
		SetText(strings.Join(model.Options.PlayerNames, ", ")).
		SetLabelColor(model.CurrentColorPalette.White).
		SetFieldWidth(20).
		SetChangedFunc(func(text string) {
			names := strings.Split(text, ",")
			for i := 0; i < len(names); i++ {
				names[i] = strings.TrimSpace(names[i])
			}
			model.Options.PlayerNames = names
		})

	// Create a dropdown for color palettes
	colorPalettes := GetColorPalettes()
	colorPaletteBox := tview.NewDropDown().
		SetLabel("Select color palette(press Enter): ").
		SetOptions(colorPalettes, func(option string, index int) {
			model.Options.ColorPalette = option
			model.CurrentColorPalette = GetColorPaletteByName(option)
		}).SetCurrentOption(GetColorPaletteIndexByName(model.Options.ColorPalette)).
		SetLabelColor(model.CurrentColorPalette.White)

	timeFormatBox := tview.NewDropDown().
		SetLabel("Select time format(press Enter): ").
		SetOptions([]string{"AMPM", "24-hour"}, func(option string, index int) {
			model.Options.TimeFormat = option
		}).SetCurrentOption(timeFormatToIndex(model.Options.TimeFormat)).
		SetLabelColor(model.CurrentColorPalette.White)

	oneTurnForAllPlayersBox := tview.NewDropDown().
		SetLabel("One Turn For All Players(press Enter): ").
		SetOptions([]string{"true", "false"}, func(option string, index int) {
			model.Options.Rules[model.Options.Default].OneTurnForAllPlayers = index == 0
		}).SetCurrentOption(boolToIndex(model.Options.Rules[model.Options.Default].OneTurnForAllPlayers)).
		SetLabelColor(model.CurrentColorPalette.White)

	optionsBox.AddItem(rulesetBox, 0, 1, false).
		AddItem(playerCountBox, 0, 1, false).
		AddItem(playerNamesBox, 0, 1, false).
		AddItem(colorPaletteBox, 0, 1, false).
		AddItem(timeFormatBox, 0, 1, false).
		AddItem(oneTurnForAllPlayersBox, 0, 1, false)

	optionsPanel.AddItem(optionsBox, 0, 0, 1, 2, 0, 0, false)

	var HelpContent strings.Builder
	HelpContent.WriteString("\n")
	HelpContent.WriteString("\n [b]Press [-]O[b] to return to the main screen")
	HelpContentBox.SetText(HelpContent.String())

	var currentRuleset strings.Builder
	currentRuleset.WriteString(" [b]Name of the ruleset:[-] " + model.Options.Rules[model.Options.Default].Name + "\n\n")
	currentRuleset.WriteString(" [b]Player Count:[-] " + fmt.Sprintf("%d", model.Options.PlayerCount) + "\n\n")
	currentRuleset.WriteString(" [b]Players:[-]\n")
	for i, name := range model.Players {
		currentRuleset.WriteString(fmt.Sprintf("  %d. %s\n", i+1, name))
	}
	currentRuleset.WriteString(" [b]Phases:[-]\n")
	for i, phase := range model.Phases {
		currentRuleset.WriteString(fmt.Sprintf("  %d. %s\n", i+1, phase))
	}
	currentRuleset.WriteString("\n")
	currentRuleset.WriteString(" [b]One Turn For All Players:[-] " + fmt.Sprintf("%t", model.Options.Rules[model.Options.Default].OneTurnForAllPlayers) + "\n\n")
	currentRuleset.WriteString(" [b]Color Palette:[-] " + model.Options.ColorPalette + "\n\n")
	currentRuleset.WriteString(" [b]Time Format:[-] " + model.Options.TimeFormat + "\n\n")
	currentRulesetContentBox.SetText(currentRuleset.String())

	optionsPanel.AddItem(currentRulesetContentBox, 1, 0, 3, 2, 0, 0, false)
	optionsPanel.AddItem(HelpContentBox, 4, 0, 1, 2, 0, 0, false)

	optionsPanel.SetBorder(true)
	optionsPanel.SetTitle(" options ")
	optionsPanel.SetBorderColor(model.CurrentColorPalette.Cyan)
	optionsPanel.SetBackgroundColor(model.CurrentColorPalette.Black)

	return optionsPanel
}

func boolToIndex(players bool) int {
	if players {
		return 0 // true
	}
	return 1 // false

}

func timeFormatToIndex(format string) int {
	if format == "AMPM" {
		return 0
	}
	return 1 // Default to 24-hour format
}

func GetColorPaletteIndexByName(palette string) int {
	for i, name := range GetColorPalettes() {
		if name == palette {
			return i
		}
	}
	return 0 // Default to the first palette if not found
}

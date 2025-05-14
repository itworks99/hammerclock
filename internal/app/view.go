package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"hammerclock/internal/app/about"
	"hammerclock/internal/app/clock"
	"hammerclock/internal/app/statusbar"

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
	hClock := clock.DisplayClock(model.Options.TimeFormat, model.CurrentColorPalette.White)
	topFlex.AddItem(hClock, 10, 0, false)

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

	// Create statusbar panel
	statusPanel := statusbar.Status(string(model.GameStatus), model.CurrentColorPalette.Cyan, model.CurrentColorPalette.Black)

	// Add statusbar panels to the main layout
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

// updateRulesetContent updates the content of the ruleset display
func updateRulesetContent(model *Model, textView *tview.Flex) {
	var leftText, rightText strings.Builder

	// Build left column content
	leftText.WriteString(fmt.Sprintf(
		" [b]Name of the ruleset:[-] %s\n\n [b]Player Count:[-] %d\n\n [b]Players:[-]\n",
		model.Options.Rules[model.Options.Default].Name,
		model.Options.PlayerCount,
	))
	for i, name := range model.Players {
		leftText.WriteString(fmt.Sprintf(" %d. %s\n", i+1, name.Name))
	}
	leftText.WriteString(fmt.Sprintf(
		"\n [b]One Turn For All Players:[-] %t\n\n [b]Color Palette:[-] %s\n",
		model.Options.Rules[model.Options.Default].OneTurnForAllPlayers,
		model.Options.ColorPalette,
	))

	// Inline color palette display
	palette := model.CurrentColorPalette
	leftText.WriteString(" [b]Palette:[-] ")
	colorBlocks := []struct {
		Name  string
		Color tcell.Color
	}{
		{"Blue", palette.Blue},
		{"Cyan", palette.Cyan},
		{"White", palette.White},
		{"DimWhite", palette.DimWhite},
		{"Yellow", palette.Yellow},
		{"Green", palette.Green},
		{"Red", palette.Red},
		{"Black", palette.Black},
	}
	for _, c := range colorBlocks {
		leftText.WriteString(fmt.Sprintf("[#%06x]█[-]", uint32(c.Color.TrueColor())))
	}
	leftText.WriteString("\n\n")

	leftText.WriteString(fmt.Sprintf(
		" [b]Time Format:[-] %s\n\n",
		model.Options.TimeFormat,
	))

	// Build right column content
	rightText.WriteString(" [b]Phases:[-]\n")
	for i, phase := range model.Phases {
		rightText.WriteString(fmt.Sprintf("  %d. %s\n", i+1, phase))
	}

	leftColumn := createTextColumn(leftText.String(), model.CurrentColorPalette.White)
	rightColumn := createTextColumn(rightText.String(), model.CurrentColorPalette.White)

	// Create grid layout
	grid := tview.NewGrid().
		AddItem(leftColumn, 0, 0, 1, 1, 0, 0, false).
		AddItem(rightColumn, 0, 1, 1, 1, 0, 0, false)

	// Clear and update the text view
	textView.Clear()
	textView.AddItem(grid, 0, 1, false)
}

// createTextColumn creates a text column with the given text
func createTextColumn(text string, color tcell.Color) *tview.TextView {
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetTextColor(color).
		SetDynamicColors(true).
		SetText(text)
}

// createOptionsScreen creates the options screen with various settings
func createOptionsScreen(model *Model) *tview.Grid {
	optionsPanel := tview.NewGrid().
		SetRows(10).
		SetColumns(0).
		SetBorders(true)

	optionsBox := tview.NewFlex().SetDirection(tview.FlexRow)
	currentRulesetContentBox := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Cache color palettes to avoid repeated calls
	colorPalettes := GetColorPalettes()

	// Create dropdown for rulesets
	rulesetBox := tview.NewDropDown().
		SetLabel("Select rules: ").
		SetOptions(getRulesetNames(model.Options.Rules), func(option string, index int) {
			model.Options.Default = index
			model.Phases = model.Options.Rules[index].Phases
			updateRulesetContent(model, currentRulesetContentBox)
		}).
		SetCurrentOption(model.Options.Default).
		SetLabelColor(model.CurrentColorPalette.White)

	// Create input field for player count
	playerCountBox := tview.NewInputField().
		SetLabel("Players: ").
		SetText(strconv.Itoa(model.Options.PlayerCount)).
		SetLabelColor(model.CurrentColorPalette.White).
		SetFieldWidth(1).
		SetChangedFunc(func(text string) {
			if count, err := strconv.Atoi(text); err == nil && count > 0 {
				model.Options.PlayerCount = count
			}
		})

	// Create player name input fields
	playerNamesBox := createPlayerNameFields(model, currentRulesetContentBox)

	// Create dropdown for color palettes
	colorPaletteBox := tview.NewDropDown().
		SetLabel("Select color palette: ").
		SetOptions(colorPalettes, func(option string, index int) {
			model.Options.ColorPalette = option
			model.CurrentColorPalette = GetColorPaletteByName(option)
			updateRulesetContent(model, currentRulesetContentBox)
		}).
		SetCurrentOption(GetColorPaletteIndexByName(model.Options.ColorPalette)).
		SetLabelColor(model.CurrentColorPalette.White)

	// Create dropdown for time format
	timeFormatBox := tview.NewDropDown().
		SetLabel("Select time format: ").
		SetOptions([]string{"AMPM", "24-hour"}, func(option string, index int) {
			model.Options.TimeFormat = option
			updateRulesetContent(model, currentRulesetContentBox)
		}).
		SetCurrentOption(timeFormatToIndex(model.Options.TimeFormat)).
		SetLabelColor(model.CurrentColorPalette.White)

	// Create checkbox for "One Turn For All Players"
	oneTurnForAllPlayersBox := tview.NewCheckbox().
		SetLabel("One Turn For All Players: ").
		SetChecked(model.Options.Rules[model.Options.Default].OneTurnForAllPlayers).
		SetChangedFunc(func(checked bool) {
			model.Options.Rules[model.Options.Default].OneTurnForAllPlayers = checked
			updateRulesetContent(model, currentRulesetContentBox)
		}).
		SetLabelColor(model.CurrentColorPalette.White)

	// Add components to options box
	optionsBox.AddItem(rulesetBox, 0, 1, false).
		AddItem(playerCountBox, 0, 1, false).
		AddItem(playerNamesBox, 0, 1, false).
		AddItem(colorPaletteBox, 0, 1, false).
		AddItem(timeFormatBox, 0, 1, false).
		AddItem(oneTurnForAllPlayersBox, 0, 1, false)

	// Add options box and help content to options panel
	optionsPanel.AddItem(optionsBox, 0, 0, 1, 2, 0, 0, false)

	helpContentBox := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White).
		SetDynamicColors(true).
		SetText("[b]Use mouse to change setting\n Press [-]O[b] to return to the main screen")

	updateRulesetContent(model, currentRulesetContentBox)

	optionsPanel.AddItem(currentRulesetContentBox, 1, 0, 3, 2, 0, 0, false)
	optionsPanel.AddItem(helpContentBox, 4, 0, 1, 2, 0, 0, false)

	optionsPanel.SetBorder(true).
		SetTitle(" options ").
		SetBorderColor(model.CurrentColorPalette.Cyan).
		SetBackgroundColor(model.CurrentColorPalette.Black)

	return optionsPanel
}

// createPlayerNameFields creates input fields for player names
func createPlayerNameFields(model *Model, currentRulesetContentBox *tview.Flex) *tview.Grid {
	playerNamesFlex := tview.NewGrid().
		SetRows(1).
		SetColumns(0).
		SetBorders(false)

	// Preallocate player names slice
	if len(model.Options.PlayerNames) < model.Options.PlayerCount {
		model.Options.PlayerNames = append(model.Options.PlayerNames, make([]string, model.Options.PlayerCount-len(model.Options.PlayerNames))...)
	}

	for i := 0; i < model.Options.PlayerCount; i++ {
		label := ""
		if i == 0 {
			label = "Player names: "
		}
		playerNamesFlex.AddItem(
			tview.NewInputField().
				SetLabel(label).
				SetText(model.Options.PlayerNames[i]).
				SetLabelColor(model.CurrentColorPalette.White).
				SetFieldWidth(10).
				SetChangedFunc(func(idx int) func(string) {
					return func(text string) {
						model.Options.PlayerNames[idx] = strings.TrimSpace(text)
						updateRulesetContent(model, currentRulesetContentBox)
					}
				}(i)),
			1, i, 1, 1, 0, 0, false,
		)
	}

	return playerNamesFlex
}

// getRulesetNames returns the names of the rulesets
func getRulesetNames(rules []Rules) []string {
	names := make([]string, len(rules))
	for i, ruleset := range rules {
		names[i] = ruleset.Name
	}
	return names
}

// timeFormatToIndex converts the time format string to an index
func timeFormatToIndex(format string) int {
	if format == "AMPM" {
		return 0
	}
	return 1 // Default to 24-hour format
}

// GetColorPaletteIndexByName returns the index of the color palette by name
func GetColorPaletteIndexByName(palette string) int {
	for i, name := range GetColorPalettes() {
		if name == palette {
			return i
		}
	}
	return 0 // Default to the first palette if not found
}

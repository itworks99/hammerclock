package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"hammerclock/components/hammerclock/Palette"
	"hammerclock/components/hammerclock/Rules"
	"hammerclock/internal/app/clock"
)

// createOptionsScreen creates the options screen with various settings
func createOptionsScreen(model *Model, msgChan chan<- Message) *tview.Grid {
	optionsPanel := tview.NewGrid().
		SetRows(10).
		SetColumns(0).
		SetBorders(true)

	optionsBox := tview.NewFlex().SetDirection(tview.FlexRow)
	currentRulesetContentBox := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Cache color palettes to avoid repeated calls
	colorPalettes := Palette.GetColorPalettes()

	// Create dropdown for rulesets
	rulesetBox := tview.NewDropDown().
		SetLabel("Select rules: ").
		SetOptions(Rules.GetRulesetNames(model.Options.Rules), nil).
		SetCurrentOption(model.Options.Default).
		SetLabelColor(model.CurrentColorPalette.White)
	// Set the changed function after initialization
	rulesetBox.SetSelectedFunc(func(option string, index int) {
		msgChan <- &SetRulesetMsg{Index: index}
	})

	// Create input field for player count
	playerCountBox := tview.NewInputField().
		SetLabel("Players: ").
		SetText(strconv.Itoa(model.Options.PlayerCount)).
		SetLabelColor(model.CurrentColorPalette.White).
		SetFieldWidth(1)

	// Set the changed function after initialization, not during
	playerCountBox.SetChangedFunc(func(text string) {
		if count, err := strconv.Atoi(text); err == nil && count > 0 {
			msgChan <- &SetPlayerCountMsg{Count: count}
		}
	})

	// Create player name input fields
	playerNamesBox := createPlayerNameFields(model, msgChan)

	// Create dropdown for color palettes
	colorPaletteBox := tview.NewDropDown().
		SetLabel("Select color palette: ").
		SetOptions(colorPalettes, nil).
		SetCurrentOption(Palette.GetColorPaletteIndexByName(model.Options.ColorPalette)).
		SetLabelColor(model.CurrentColorPalette.White)
	// Set the changed function after initialization
	colorPaletteBox.SetSelectedFunc(func(option string, index int) {
		msgChan <- &SetColorPaletteMsg{Name: option}
	})

	// Create dropdown for time format
	timeFormatBox := tview.NewDropDown().
		SetLabel("Select time format: ").
		SetOptions([]string{"AMPM", "24-hour"}, nil).
		SetCurrentOption(clock.TimeFormatToIndex(model.Options.TimeFormat)).
		SetLabelColor(model.CurrentColorPalette.White)
	// Set the changed function after initialization
	timeFormatBox.SetSelectedFunc(func(option string, index int) {
		msgChan <- &SetTimeFormatMsg{Format: option}
	})

	// Create checkbox for "One Turn For All Players"
	oneTurnForAllPlayersBox := tview.NewCheckbox().
		SetLabel("One Turn For All Players: ").
		SetChecked(model.Options.Rules[model.Options.Default].OneTurnForAllPlayers).
		SetLabelColor(model.CurrentColorPalette.White)
	// Set the changed function after initialization
	oneTurnForAllPlayersBox.SetChangedFunc(func(checked bool) {
		msgChan <- &SetOneTurnForAllPlayersMsg{Value: checked}
	})

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

	// Add a message handler to update content on model changes
	updateRulesetContent(model, currentRulesetContentBox)

	// Observe model changes and update UI accordingly
	// This would be handled by the Render function when model updates

	optionsPanel.AddItem(currentRulesetContentBox, 1, 0, 3, 2, 0, 0, false)
	optionsPanel.AddItem(helpContentBox, 4, 0, 1, 2, 0, 0, false)

	optionsPanel.SetBorder(true).
		SetTitle(" options ").
		SetBorderColor(model.CurrentColorPalette.Cyan).
		SetBackgroundColor(model.CurrentColorPalette.Black)

	return optionsPanel
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

// createPlayerNameFields creates input fields for player names
func createPlayerNameFields(model *Model, msgChan chan<- Message) *tview.Grid {
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

		// Create the input field without setting the changed function initially
		inputField := tview.NewInputField().
			SetLabel(label).
			SetText(model.Options.PlayerNames[i]).
			SetLabelColor(model.CurrentColorPalette.White).
			SetFieldWidth(10)

		// Store index in a closure to avoid variable capture issues
		idx := i
		inputField.SetChangedFunc(func(text string) {
			msgChan <- &SetPlayerNameMsg{
				Index: idx,
				Name:  strings.TrimSpace(text),
			}
		})

		playerNamesFlex.AddItem(
			inputField,
			1, i, 1, 1, 0, 0, false)
	}

	return playerNamesFlex
}

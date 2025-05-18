package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"hammerclock/internal/hammerclock/common"
)

// CreatePlayerPanel creates a player panel
func CreatePlayerPanel(player *common.Player, color string, model *common.Model) *tview.Flex {
	panel := tview.NewFlex().SetDirection(tview.FlexRow)
	upper := tview.NewFlex().SetDirection(tview.FlexRow)
	lower := tview.NewFlex().SetDirection(tview.FlexRow)

	playerName := tview.NewTextView().
		SetText("\nPlayer: " + player.Name).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)
	elapsedTime := tview.NewTextView().
		SetText(fmt.Sprintf("Time Elapsed: %v", player.TimeElapsed)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)
	horizontalDivider := tview.NewTextView().
		SetText(strings.Repeat("─", 30)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.DimWhite)
	currentTurnAndPhase := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)

	setPhaseText := func() {
		if !model.Options.Rules[model.Options.Default].OneTurnForAllPlayers {
			currentTurnAndPhase.SetText(fmt.Sprintf("Turn: %d | Phase: %s", player.TurnCount, model.Phases[player.CurrentPhase]))
		} else {
			currentTurnAndPhase.SetText(fmt.Sprintf("Turn: %d", player.TurnCount))
		}
	}
	setPhaseText()

	upper.AddItem(playerName, 2, 1, false).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(elapsedTime, 1, 1, false).
		AddItem(horizontalDivider, 1, 0, false).
		AddItem(currentTurnAndPhase, 1, 1, false).
		AddItem(tview.NewBox(), 0, 1, false)

	logTitle := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("\nAction Log:").
		SetTextColor(model.CurrentColorPalette.White)

	// Creating a scrollable log view
	logView := createLogView()

	// Set initial content if any exists
	if len(player.ActionLog) > 0 {
		// Use LogPanel.SetLogContent to consistently format log entries
		SetLogContent(logView, player.ActionLog)
	}

	// CreateAboutPanel a container with the log view
	logContainer := createLogContainer(logView)
	lower.AddItem(logTitle, 3, 0, false)
	lower.AddItem(logContainer, 0, 1, true)

	borderColor := model.CurrentColorPalette.Black
	switch color {
	case "blue":
		borderColor = model.CurrentColorPalette.Blue
	case "yellow":
		borderColor = model.CurrentColorPalette.Yellow
	case "green":
		borderColor = model.CurrentColorPalette.Green
	case "red":
		borderColor = model.CurrentColorPalette.Red
	}

	panel.AddItem(upper, 7, 0, false)
	panel.AddItem(lower, 0, 3, true)
	panel.SetBorder(true).
		SetBackgroundColor(model.CurrentColorPalette.Black).
		SetBorderColor(borderColor)
	horizontalDivider.SetTextColor(borderColor)

	// Add mouse capture for smooth player selection
	panel.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick {
			// Only select if not already selected
			if !player.IsTurn {
				for _, p := range model.Players {
					if p == player {
						// Set this player as active
						for i := range model.Players {
							model.Players[i].IsTurn = false
						}
						player.IsTurn = true

						// Set focus on the panel to show double-line border
						panel.Focus(func(p tview.Primitive) {
							// Delegate function - we don't need to do anything here
						})
						break
					}
				}
			}
		}
		return action, event
	})

	return panel
}

// UpdatePlayerPanels updates the player panels with the current player data
func UpdatePlayerPanels(players []*common.Player, panels []*tview.Flex, model *common.Model) {
	for i, player := range players {
		currentPlayerPanel := panels[i].GetItem(0).(*tview.Flex)
		gameInfoBox := currentPlayerPanel.GetItem(0).(*tview.TextView)
		elapsedTimeBox := currentPlayerPanel.GetItem(2).(*tview.TextView)
		horizontalDivider := currentPlayerPanel.GetItem(3).(*tview.TextView)
		currentTurnAndPhase := currentPlayerPanel.GetItem(4).(*tview.TextView)

		elapsedTimeBox.SetText(fmt.Sprintf("Time Elapsed: %v", player.TimeElapsed))
		if !model.Options.Rules[model.Options.Default].OneTurnForAllPlayers {
			currentTurnAndPhase.SetText(fmt.Sprintf("Turn: %d | Phase: %s", player.TurnCount, model.Phases[player.CurrentPhase]))
		} else {
			currentTurnAndPhase.SetText(fmt.Sprintf("Turn: %d", player.TurnCount))
		}

		if !model.GameStarted {
			panels[i].SetTitle("")
			gameInfoBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			elapsedTimeBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			currentTurnAndPhase.SetTextColor(model.CurrentColorPalette.DimWhite)
			panels[i].Blur() // Remove focus
		} else if player.IsTurn {
			panels[i].SetTitle(" ACTIVE TURN ")
			gameInfoBox.SetTextColor(model.CurrentColorPalette.White)
			elapsedTimeBox.SetTextColor(model.CurrentColorPalette.White)
			currentTurnAndPhase.SetTextColor(model.CurrentColorPalette.White)
			// Set focus to get double-line border
			panels[i].Focus(func(p tview.Primitive) {
				// Delegate function - we don't need to do anything here
			})
		} else {
			panels[i].SetTitle("")
			gameInfoBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			elapsedTimeBox.SetTextColor(model.CurrentColorPalette.DimWhite)
			currentTurnAndPhase.SetTextColor(model.CurrentColorPalette.DimWhite)
			panels[i].Blur() // Remove focus
		}
		horizontalDivider.SetTextColor(panels[i].GetBorderColor())

		lower := panels[i].GetItem(1).(*tview.Flex)
		if lower != nil && lower.GetItemCount() > 1 {
			logContainer := lower.GetItem(1).(*tview.Flex)
			// The log container has the log view as its only item now
			logView := logContainer.GetItem(0).(*tview.TextView)

			// Update log panel content
			SetLogContent(logView, player.ActionLog)
		}
	}
}

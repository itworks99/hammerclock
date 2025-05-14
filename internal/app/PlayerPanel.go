package app

import (
	"fmt"
	"strings"
	"time"

	"hammerclock/internal/app/LogPanel"

	"github.com/rivo/tview"
)

func CreatePlayerPanel(player *Player, color string, model *Model) *tview.Flex {
	panel := tview.NewFlex().SetDirection(tview.FlexRow)
	upper := tview.NewFlex().SetDirection(tview.FlexRow)
	lower := tview.NewFlex().SetDirection(tview.FlexRow)

	name := tview.NewTextView().
		SetText("\nPlayer: " + player.Name).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)
	elapsedTime := tview.NewTextView().
		SetText(fmt.Sprintf("Time Elapsed: %v", player.TimeElapsed)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)
	line := tview.NewTextView().
		SetText(strings.Repeat("─", 30)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.DimWhite)
	phase := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextColor(model.CurrentColorPalette.White)

	setPhaseText := func() {
		if !model.Options.Rules[model.Options.Default].OneTurnForAllPlayers {
			phase.SetText(fmt.Sprintf("Turn: %d | Phase: %s", player.TurnCount, model.Phases[player.CurrentPhase]))
		} else {
			phase.SetText(fmt.Sprintf("Turn: %d", player.TurnCount))
		}
	}
	setPhaseText()

	upper.AddItem(name, 2, 1, false).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(elapsedTime, 1, 1, false).
		AddItem(line, 1, 0, false).
		AddItem(phase, 1, 1, false).
		AddItem(tview.NewBox(), 0, 1, false)

	logTitle := tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetText("\nAction Log:").
		SetTextColor(model.CurrentColorPalette.White)

	// Use the logpanel package to create a scrollable log view
	logView := LogPanel.CreateLogView()

	// Set initial content if any exists
	if len(player.ActionLog) > 0 {
		var b strings.Builder
		for _, entry := range player.ActionLog {
			b.WriteString(fmt.Sprintf("[%s] Turn %d, Phase %s: %s\n",
				entry.DateTime.Format("15:04:05"), entry.Turn, entry.Phase, entry.Message))
		}
		logView.SetText(b.String())
	}

	// Create a container with the log view
	logContainer := LogPanel.CreateLogContainer(logView)
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
	line.SetTextColor(borderColor)
	return panel
}

func updatePlayerPanels(players []*Player, panels []*tview.Flex, model *Model) {
	for i, player := range players {
		upper := panels[i].GetItem(0).(*tview.Flex)
		name := upper.GetItem(0).(*tview.TextView)
		playerTime := upper.GetItem(2).(*tview.TextView)
		line := upper.GetItem(3).(*tview.TextView)
		phase := upper.GetItem(4).(*tview.TextView)

		playerTime.SetText(fmt.Sprintf("Time Elapsed: %v", player.TimeElapsed))
		if !model.Options.Rules[model.Options.Default].OneTurnForAllPlayers {
			phase.SetText(fmt.Sprintf("Turn: %d | Phase: %s", player.TurnCount, model.Phases[player.CurrentPhase]))
		} else {
			phase.SetText(fmt.Sprintf("Turn: %d", player.TurnCount))
		}

		if !model.GameStarted {
			panels[i].SetTitle("")
			name.SetTextColor(model.CurrentColorPalette.DimWhite)
			playerTime.SetTextColor(model.CurrentColorPalette.DimWhite)
			phase.SetTextColor(model.CurrentColorPalette.DimWhite)
		} else if player.IsTurn {
			panels[i].SetTitle(" ACTIVE TURN ")
			name.SetTextColor(model.CurrentColorPalette.White)
			playerTime.SetTextColor(model.CurrentColorPalette.White)
			phase.SetTextColor(model.CurrentColorPalette.White)
		} else {
			panels[i].SetTitle("")
			name.SetTextColor(model.CurrentColorPalette.DimWhite)
			playerTime.SetTextColor(model.CurrentColorPalette.DimWhite)
			phase.SetTextColor(model.CurrentColorPalette.DimWhite)
		}
		line.SetTextColor(panels[i].GetBorderColor())

		lower := panels[i].GetItem(1).(*tview.Flex)
		if lower != nil && lower.GetItemCount() > 1 {
			logContainer := lower.GetItem(1).(*tview.Flex)
			// The log container has the log view as its only item now
			logView := logContainer.GetItem(0).(*tview.TextView)

			// Use the logpanel package to update log content
			LogPanel.SetLogContent(logView, player.ActionLog)
		}
	}
}

// AddLogEntry adds a log entry to a player's action log
func AddLogEntry(player *Player, format string, args ...any) {
	currentPhase := ""
	if player.CurrentPhase < len(DefaultOptions.Rules[DefaultOptions.Default].Phases) && player.CurrentPhase >= 0 {
		currentPhase = DefaultOptions.Rules[DefaultOptions.Default].Phases[player.CurrentPhase]
	}

	logEntry := LogPanel.LogEntry{
		DateTime:   time.Now(),
		PlayerName: player.Name,
		Turn:       player.TurnCount,
		Phase:      currentPhase,
		Message:    fmt.Sprintf(format, args...),
	}

	player.ActionLog = append(player.ActionLog, logEntry)
}

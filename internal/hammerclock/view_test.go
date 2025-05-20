package hammerclock

import (
	"testing"

	"github.com/rivo/tview"
	"hammerclock/internal/hammerclock/common"
	"hammerclock/internal/hammerclock/options"
	"hammerclock/internal/hammerclock/rules"
)

var testModel = &common.Model{
	Players: []*common.Player{
		{Name: "Player 1", IsTurn: true},
		{Name: "Player 2"},
	},
	Phases:        []string{"Setup", "Movement", "Shooting", "Melee", "End"},
	GameStatus:    gameNotStarted,
	CurrentScreen: "main",
	Options: options.Options{
		TimeFormat: "24h",
		Rules: []rules.Rules{
			{
				Name:                 "Default Rules",
				Phases:               []string{"Setup", "Movement", "Shooting", "Melee", "End"},
				OneTurnForAllPlayers: true,
			},
		},
		Default: 0,
	},
	TotalGameTime: 0,
}

func TestNewView(t *testing.T) {
	model := testModel

	view := NewView(model, make(chan common.Message, 10))

	if view == nil || view.App == nil || view.MainView == nil {
		t.Fatal("View or essential components are nil")
	}
	if len(view.PlayerPanels) != len(model.Players) {
		t.Errorf("Expected %d player panels, got %d", len(model.Players), len(view.PlayerPanels))
	}
}

func TestUpdateClock(t *testing.T) {
	model := testModel
	view := NewView(model, make(chan common.Message, 10))

	view.UpdateClock(model)
	if view.ClockDisplay.GetText(false) == "" {
		t.Error("Clock display is empty after UpdateClock")
	}
}

func TestRender(t *testing.T) {
	model := testModel
	view := NewView(model, make(chan common.Message, 10))

	view.Render(model)
	if view.CurrentScreen != "main" {
		t.Errorf("Expected 'main', got '%s'", view.CurrentScreen)
	}
}

func TestCreateEndGameConfirmationModal(t *testing.T) {

	view := NewView(testModel, make(chan common.Message, 10))
	modal := CreateEndGameConfirmationModal(view)

	if modal == nil {
		t.Error("Modal creation failed")
	}
}

func TestShowConfirmationModal(t *testing.T) {
	view := NewView(testModel, make(chan common.Message, 10))
	modal := tview.NewModal().SetText("Test Modal")
	ShowConfirmationModal(view, modal)
}

func TestRestoreMainView(t *testing.T) {
	view := NewView(testModel, make(chan common.Message, 10))
	view.RestoreMainView()
}

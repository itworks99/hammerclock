package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"strings"
	"time"
)

type Player struct {
	name         string
	timeElapsed  time.Duration
	isTurn       bool
	currentPhase int
	armyList     []Unit
}

type Unit struct {
	Name   string
	Points int
}

type HCPanel struct {
	xStart int
	xEnd   int
	yStart int
	yEnd   int
	style  tcell.Style
	screen tcell.Screen
}

type Settings struct {
	Name                 string   `json:"name"`
	Default              bool     `json:"default"`
	PlayerCount          int      `json:"playerCount"`
	Phases               []string `json:"phases"`
	OneTurnForAllPlayers bool     `json:"oneTurnForAllPlayers"`
}

var defaultRulesFile = "defaultRules.json"

var phases []string

var menuOptions = []string{
	"1: Load rules",
	"2: Settings",
	"A: About",
}

var instructions = []string{
	"S: Start Game",
	"SPACE: Switch Turns",
	"P: Next Phase",
	"B: Previous Phase",
	"Q: Quit",
}

func readSettingsFile() Settings {
	var defaultSettings Settings
	byteValue, err := os.ReadFile(defaultRulesFile)
	if err != nil {
		fmt.Println("Error reading settings file:", err)
	}
	err = json.Unmarshal(byteValue, &defaultSettings)
	if err != nil {
		fmt.Println("Error processing settings file:", err)
	}
	return defaultSettings
}

func drawMenu(screen tcell.Screen, position string) {
	menuStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	keyStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	_, height := screen.Size()

	var options = menuOptions
	y := 0
	if position == "bottom" {
		y = height - 1
		options = instructions
	}

	// Draw menu bar
	for i, option := range options {
		parts := strings.SplitN(option, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		item := parts[1]

		// Draw key
		for j, char := range key {
			screen.SetContent(j+i*20, y, char, nil, keyStyle)
		}

		// Draw separator
		screen.SetContent(len(key)+i*20, y, ' ', nil, keyStyle)
		screen.SetContent(len(key)+1+i*20, y, ' ', nil, keyStyle)

		// Draw item
		for j, char := range item {
			screen.SetContent(len(key)+1+j+i*20, y, char, nil, menuStyle)
		}
	}
}

// NewHCPanel creates a new panel with the given parameters
func NewHCPanel(screen tcell.Screen, xStart, xEnd, yStart, yEnd int, style tcell.Style) *HCPanel {
	return &HCPanel{
		xStart: xStart,
		xEnd:   xEnd,
		yStart: yStart,
		yEnd:   yEnd,
		style:  style,
		screen: screen,
	}
}

// Draw renders the borders of the HCPanel on the screen using the specified style and the provided isTurn state.
func (p *HCPanel) Draw(isTurn bool) {
	borderStyle := p.style.Reverse(isTurn)
	for x := p.xStart; x < p.xEnd; x++ {
		p.screen.SetContent(x, p.yStart, tcell.RuneHLine, nil, borderStyle)
		p.screen.SetContent(x, p.yEnd-1, tcell.RuneHLine, nil, borderStyle)
	}
	for y := p.yStart; y < p.yEnd; y++ {
		p.screen.SetContent(p.xStart, y, tcell.RuneVLine, nil, borderStyle)
		p.screen.SetContent(p.xEnd-1, y, tcell.RuneVLine, nil, borderStyle)
	}

	// Corners
	p.screen.SetContent(p.xStart, p.yStart, tcell.RuneULCorner, nil, borderStyle)
	p.screen.SetContent(p.xEnd-1, p.yStart, tcell.RuneURCorner, nil, borderStyle)
	p.screen.SetContent(p.xStart, p.yEnd-1, tcell.RuneLLCorner, nil, borderStyle)
	p.screen.SetContent(p.xEnd-1, p.yEnd-1, tcell.RuneLRCorner, nil, borderStyle)
}

// DrawPlayerInfo draws the player information inside the panel
func (p *HCPanel) DrawPlayerInfo(player *Player) {
	width := p.xEnd - p.xStart

	// Wrap player name
	playerNameLines := "Player: " + player.name
	for i, line := range playerNameLines {
		p.screen.SetContent(p.xStart+2+i, p.yStart+2, line, nil, p.style)
	}

	// Draw player time elapsed
	timeElapsed := fmt.Sprintf("Time Elapsed: %v", player.timeElapsed)
	for i, char := range timeElapsed {
		p.screen.SetContent(p.xStart+(width-len(timeElapsed))/2+i, p.yStart+len(playerNameLines)+3, char, nil, p.style)
	}

	// Draw current phase
	currentPhase := fmt.Sprintf("Phase: %s", phases[player.currentPhase])
	for i, char := range currentPhase {
		p.screen.SetContent(p.xStart+(width-len(currentPhase))/2+i, p.yStart+len(playerNameLines)+5, char, nil, p.style)
	}
}

// Dummy function for loading BattleScribe lists (to be implemented)
func loadArmyListFromFile(player *Player) {
	// For simplicity, let's load a fake army list
	player.armyList = []Unit{
		{"Space Marine", 20},
		{"Tactical Squad", 80},
		{"Predator Tank", 150},
	}
}

func drawScreen(screen tcell.Screen, players []*Player, gameStarted bool) {
	screen.Clear()
	width, height := screen.Size()
	numberOfPlayers := len(players)

	// Define styles for each player
	playerStyles := []tcell.Style{
		tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack),
		tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack),
		tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack),
		tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack),
	}

	drawMenu(screen, "top")
	drawMenu(screen, "bottom")

	// Calculate the width of each panel
	panelWidth := width / numberOfPlayers

	// Draw each player's panel using HCPanel
	for i, player := range players {
		xStart := i * panelWidth
		xEnd := xStart + panelWidth

		// Create a panel for this player
		panel := NewHCPanel(screen, xStart, xEnd, 1, height-1, playerStyles[i%len(playerStyles)])

		// Draw the panel and player info
		panel.Draw(player.isTurn)
		panel.DrawPlayerInfo(player)
	}

	// Display game status
	status := "Game Not Started"
	if gameStarted {
		status = "Game In Progress"
	}
	for i, char := range status {
		screen.SetContent((width-len(status))/2+i, height-3, char, nil, tcell.StyleDefault)
	}

	screen.Show()
}

func main() {
	defaultSettings := readSettingsFile()
	numberOfPlayers := defaultSettings.PlayerCount
	phases = defaultSettings.Phases

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Error creating screen:", err)
		return
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		fmt.Println("Error initializing screen:", err)
		return
	}

	// Create players
	players := make([]*Player, numberOfPlayers)
	for i := 0; i < numberOfPlayers; i++ {
		players[i] = &Player{name: fmt.Sprintf("Player %d", i+1), timeElapsed: 0, isTurn: i == 0, currentPhase: 0}
	}

	gameStarted := false

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if gameStarted {
					for _, player := range players {
						if player.isTurn {
							player.timeElapsed += 1 * time.Second
						}
					}
					err := screen.PostEvent(nil)
					if err != nil {
						return
					}
				}
			case <-quit:
				return
			}
		}
	}()

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			// Check for Alt modifier and ignore those events to prevent terminal scrolling
			if ev.Modifiers()&tcell.ModAlt != 0 {
				continue
			}
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(quit)
				return
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'q', 'Q':
					screen.Fill(' ', tcell.StyleDefault)
					close(quit)
					return
				case 's', 'S':
					if !gameStarted {
						gameStarted = true
					}
					drawScreen(screen, players, gameStarted)
					screen.Show()
				case '1':
					fmt.Println("Load Army List selected") // Placeholder for file-loading functionality
					drawScreen(screen, players, gameStarted)
					screen.Show()
				case '2':
					fmt.Println("Build Army List selected") // Placeholder for manual input
					drawScreen(screen, players, gameStarted)
					screen.Show()
				case '3':
					fmt.Println("Start Game selected")
					drawScreen(screen, players, gameStarted)
					screen.Show()
				case ' ':
					// Switch turns
					for _, player := range players {
						player.isTurn = !player.isTurn
						if player.isTurn {
							player.currentPhase = 0
						}
					}
					drawScreen(screen, players, gameStarted)
					screen.Show()
				case 'p', 'P':
					// Move forward in the phase
					for _, player := range players {
						if player.isTurn {
							player.currentPhase = (player.currentPhase + 1) % len(phases)
						}
					}
					drawScreen(screen, players, gameStarted)
					screen.Show()
				case 'b', 'B':
					// Move backward in the phase
					for _, player := range players {
						if player.isTurn {
							player.currentPhase = (player.currentPhase - 1 + len(phases)) % len(phases)
						}
					}
					drawScreen(screen, players, gameStarted)
					screen.Show()
				case 'l', 'L':
					// Load army lists
					for _, player := range players {
						if player.isTurn {
							loadArmyListFromFile(player)
						}
					}
					drawScreen(screen, players, gameStarted)
					screen.Show()
				}
			default:
				continue
			}
		case *tcell.EventResize:
			screen.Sync()
			drawScreen(screen, players, gameStarted)
			screen.Show()
		case nil:
			// Handle nil events posted by the ticker
			drawScreen(screen, players, gameStarted)
			screen.Show()
		}
	}
}

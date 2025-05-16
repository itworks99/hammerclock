package palette

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ColorPalette contains all the colors used in the application
type ColorPalette struct {
	Blue     tcell.Color
	Cyan     tcell.Color
	White    tcell.Color
	DimWhite tcell.Color
	Yellow   tcell.Color
	Green    tcell.Color
	Red      tcell.Color
	Black    tcell.Color
}

// K9sPalette K9s color palette
var K9sPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(36, 96, 146),   // Dark blue for backgrounds
	Cyan:     tcell.NewRGBColor(0, 183, 235),   // Cyan for highlights
	White:    tcell.NewRGBColor(255, 255, 255), // White for primary text
	DimWhite: tcell.NewRGBColor(180, 180, 180), // Dimmed white for inactive panels
	Yellow:   tcell.NewRGBColor(253, 185, 19),  // Yellow for warnings
	Green:    tcell.NewRGBColor(0, 200, 83),    // Green for success/active states
	Red:      tcell.NewRGBColor(255, 0, 0),     // Red for errors/critical states
	Black:    tcell.NewRGBColor(0, 0, 0),       // Black for default backgrounds
}

// DraculaPalette Dracula color palette
var DraculaPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(189, 147, 249), // Purple
	Cyan:     tcell.NewRGBColor(139, 233, 253), // Cyan
	White:    tcell.NewRGBColor(248, 248, 242), // Foreground
	DimWhite: tcell.NewRGBColor(174, 174, 169), // Dimmed foreground
	Yellow:   tcell.NewRGBColor(241, 250, 140), // Yellow
	Green:    tcell.NewRGBColor(80, 250, 123),  // Green
	Red:      tcell.NewRGBColor(255, 85, 85),   // Red
	Black:    tcell.NewRGBColor(40, 42, 54),    // Background
}

// MonokaiPalette Monokai color palette
var MonokaiPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(102, 217, 239), // Blue
	Cyan:     tcell.NewRGBColor(102, 217, 239), // Blue (same as Blue for Monokai)
	White:    tcell.NewRGBColor(248, 248, 242), // Foreground
	DimWhite: tcell.NewRGBColor(174, 174, 169), // Dimmed foreground
	Yellow:   tcell.NewRGBColor(230, 219, 116), // Yellow
	Green:    tcell.NewRGBColor(166, 226, 46),  // Green
	Red:      tcell.NewRGBColor(249, 38, 114),  // Red
	Black:    tcell.NewRGBColor(39, 40, 34),    // Background
}

// WarhammerPalette represents the color theme for Warhammer 40K
var WarhammerPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(38, 57, 132),   // Ultramarine Blue
	Cyan:     tcell.NewRGBColor(23, 155, 215),  // Tyranid Blue
	White:    tcell.NewRGBColor(255, 250, 240), // Imperial White
	DimWhite: tcell.NewRGBColor(180, 170, 150), // Bone Color
	Yellow:   tcell.NewRGBColor(245, 180, 26),  // Imperial Gold
	Green:    tcell.NewRGBColor(0, 120, 50),    // Dark Angels Green
	Red:      tcell.NewRGBColor(190, 0, 0),     // Blood Angels Red
	Black:    tcell.NewRGBColor(10, 10, 10),    // Abaddon Black
}

// KillTeamPalette represents the color theme for Kill Team
var KillTeamPalette = ColorPalette{
	Blue:     tcell.NewRGBColor(63, 81, 153),   // Night Lords Blue
	Cyan:     tcell.NewRGBColor(0, 169, 157),   // Tactical Turquoise
	White:    tcell.NewRGBColor(230, 230, 230), // Tactical White
	DimWhite: tcell.NewRGBColor(150, 150, 150), // Urban Gray
	Yellow:   tcell.NewRGBColor(255, 193, 0),   // Warning Yellow
	Green:    tcell.NewRGBColor(76, 99, 25),    // Camo Green
	Red:      tcell.NewRGBColor(200, 40, 40),   // Target Red
	Black:    tcell.NewRGBColor(5, 5, 5),       // Shadow Black
}

// ColorPalettes returns a list of available color palettes
func ColorPalettes() []string {
	return []string{
		"k9s",
		"dracula",
		"monokai",
		"warhammer",
		"killteam",
	}
}

// ColorPaletteByName returns the color palette for the given name
func ColorPaletteByName(name string) ColorPalette {
	switch name {
	case "dracula":
		return DraculaPalette
	case "monokai":
		return MonokaiPalette
	case "warhammer":
		return WarhammerPalette
	case "killteam":
		return KillTeamPalette
	default: // "k9s" or any other value defaults to k9s
		return K9sPalette
	}
}

// ApplyColorPalette applies the color palette to tview styles
func ApplyColorPalette(palette ColorPalette) {
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

// ColorPaletteIndexByName returns the index of the color palette by name
func ColorPaletteIndexByName(palette string) int {
	for i, name := range ColorPalettes() {
		if name == palette {
			return i
		}
	}
	return 0 // Default to the first palette if not found
}

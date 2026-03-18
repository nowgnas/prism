package tab

import "github.com/charmbracelet/lipgloss"

// Palette is the ordered list of tab accent colors.
var Palette = []lipgloss.Color{
	"#E05C5C", // 0 Crimson
	"#E0A05C", // 1 Amber
	"#5CC4A0", // 2 Jade
	"#5C8CE0", // 3 Cobalt
	"#A05CE0", // 4 Violet
	"#E05CA0", // 5 Rose
}

// ColorForIndex returns the tab color for a given index (wraps around the palette).
func ColorForIndex(idx int) lipgloss.Color {
	if len(Palette) == 0 {
		return "#FFFFFF"
	}
	return Palette[idx%len(Palette)]
}

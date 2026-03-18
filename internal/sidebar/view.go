package sidebar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/nowgnas/prism/internal/tab"
)

var (
	sidebarBg = lipgloss.Color("#1A1A2E")
	dimColor  = lipgloss.Color("#444466")
	textColor = lipgloss.Color("#CCCCDD")

	sidebarStyle = lipgloss.NewStyle().
			Background(sidebarBg).
			PaddingLeft(1).
			PaddingRight(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888899")).
			Background(sidebarBg).
			Bold(true).
			PaddingLeft(1)

	dividerStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Background(sidebarBg)
)

// RenderSidebar returns the full sidebar as a single string.
func RenderSidebar(m *Model, tabs []*tab.Model) string {
	lines := make([]string, 0, m.Height)

	// Header
	header := headerStyle.Width(m.Width).Render("PRISM")
	lines = append(lines, header)

	divider := dividerStyle.Width(m.Width).Render(strings.Repeat("─", m.Width))
	lines = append(lines, divider)

	// Tab entries
	for i, t := range tabs {
		lines = append(lines, renderTabEntry(m, i, t, i == m.ActiveTab))
	}

	// Pad remaining height
	for len(lines) < m.Height {
		lines = append(lines, sidebarStyle.Width(m.Width).Render(""))
	}

	return strings.Join(lines[:m.Height], "\n")
}

func renderTabEntry(m *Model, idx int, t *tab.Model, active bool) string {
	iconStr := statusStyle(t.Status, t.Color)

	// Tab index number, colored with tab color
	numStr := lipgloss.NewStyle().
		Foreground(t.Color).
		Background(sidebarBg).
		Bold(active).
		Render(fmt.Sprintf("%d", idx+1))

	// Tab name truncated to fit
	maxNameLen := m.Width - 6 // account for icon + number + spaces
	name := t.Name
	if len(name) > maxNameLen {
		name = name[:maxNameLen-1] + "…"
	}

	var nameStyle lipgloss.Style
	if active {
		nameStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(sidebarBg).
			Bold(true)
	} else {
		nameStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Background(sidebarBg)
	}
	nameStr := nameStyle.Render(name)

	content := fmt.Sprintf(" %s %s %s", numStr, iconStr, nameStr)

	var rowStyle lipgloss.Style
	if active {
		rowStyle = lipgloss.NewStyle().
			Background(sidebarBg).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(t.Color).
			Width(m.Width - 1)
	} else {
		rowStyle = lipgloss.NewStyle().
			Background(sidebarBg).
			Width(m.Width)
	}

	return rowStyle.Render(content)
}

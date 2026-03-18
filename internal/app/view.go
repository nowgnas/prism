package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/nowgnas/prism/internal/layout"
	"github.com/nowgnas/prism/internal/tab"
	"github.com/nowgnas/prism/internal/terminal"
)


// View renders the full TUI.
func (m *AppModel) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "initializing…"
	}
	if len(m.Tabs) == 0 {
		return "no tabs"
	}

	sidebarView := m.Sidebar.View(m.Tabs)
	tabView := m.renderActiveTab()

	// Status bar
	statusBar := m.renderStatusBar()

	// Combine sidebar + tab area side by side, then status bar below
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, tabView)

	if m.showHelp {
		return m.renderHelpOverlay()
	}
	if m.Mode == ModeRename {
		prompt := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCDD")).
			Render("Rename tab: ") + m.Rename.View()
		lines := strings.Split(body, "\n")
		if len(lines) > 0 {
			lines[len(lines)-1] = prompt
		}
		return strings.Join(lines, "\n")
	}

	return body + "\n" + statusBar
}

// renderActiveTab renders the current tab's pane grid.
func (m *AppModel) renderActiveTab() string {
	t := m.activeTab()
	if t == nil {
		return ""
	}
	tabW, tabH := m.tabAreaDimensions()
	return renderTabPanes(t, tabW, tabH, m.Spinner.View())
}

// renderTabPanes renders all panes of a tab within the given area.
func renderTabPanes(t *tab.Model, w, h int, spinnerStr string) string {
	n := len(t.Terminals)
	if n == 0 {
		return lipgloss.NewStyle().Width(w).Height(h).Render("no terminals")
	}

	// If zoomed, render only the zoomed pane full-size
	if t.ZoomedPane >= 0 && t.ZoomedPane < n {
		return renderPane(t.Terminals[t.ZoomedPane], layout.Rect{0, 0, w, h}, t.Color, spinnerStr)
	}

	rects := layout.ComputeLayout(n, w, h)
	if len(rects) == 0 {
		return ""
	}

	// Build a character canvas
	canvas := makeCanvas(w, h)

	for i, term := range t.Terminals {
		if i >= len(rects) {
			break
		}
		paneStr := renderPane(term, rects[i], t.Color, spinnerStr)
		blitToCanvas(canvas, rects[i].X, rects[i].Y, paneStr)
	}

	return canvasToString(canvas, w, h)
}

// renderPane renders a single terminal pane with a border.
func renderPane(term *terminal.Model, rect layout.Rect, tabColor lipgloss.Color, spinnerStr string) string {
	inner := layout.InnerRect(rect)
	if inner.W <= 0 || inner.H <= 0 {
		return ""
	}

	// Get terminal content
	var content string
	if term.Screen != nil {
		content = term.Screen.RenderANSI()
	} else {
		content = ""
	}

	// Truncate/pad content to inner dimensions
	content = fitContent(content, inner.W, inner.H)

	// Choose border color
	var borderColor lipgloss.Color
	if term.Focused {
		borderColor = tabColor
	} else {
		borderColor = "#444444"
	}

	// Status indicator in top-right of border
	statusStr := termStatusStr(term, tabColor, spinnerStr)

	// Build border title
	title := fmt.Sprintf(" %d%s ", term.ID, statusStr)

	paneStyle := lipgloss.NewStyle().
		Width(inner.W).
		Height(inner.H).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	_ = title // lipgloss border title support varies; use simple border for now
	return paneStyle.Render(content)
}

func termStatusStr(term *terminal.Model, tabColor lipgloss.Color, spinnerStr string) string {
	switch term.Status {
	case terminal.StatusRunning:
		return " " + lipgloss.NewStyle().Foreground(tabColor).Render(spinnerStr)
	case terminal.StatusAI:
		return " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Render(spinnerStr)
	case terminal.StatusError:
		return " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333")).Render("✕")
	}
	return ""
}

// renderStatusBar renders the bottom status line.
func (m *AppModel) renderStatusBar() string {
	var hint string
	switch m.Mode {
	case ModePrefix:
		hint = "^B — awaiting command (? for help)"
	case ModeRename:
		hint = "Rename: Enter to confirm, Esc to cancel"
	default:
		hint = "^B → prefix  |  ^B ? → help  |  ^B q → quit"
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888899")).
		Background(lipgloss.Color("#111122")).
		Width(m.Width).
		Padding(0, 1)

	return style.Render(hint)
}

// renderHelpOverlay returns the help screen.
func (m *AppModel) renderHelpOverlay() string {
	lines := []string{
		"",
		"  PRISM — Key Bindings",
		"  " + strings.Repeat("─", 36),
		"",
		"  ^B n      New tab",
		"  ^B w      Close current tab",
		"  ^B ,      Rename tab",
		"  ^B Tab    Next tab",
		"  ^B 1-9    Jump to tab N",
		"",
		"  ^B \"      Add pane (horizontal split)",
		"  ^B %      Add pane (vertical split)",
		"  ^B x      Close current pane",
		"  ^B z      Zoom / unzoom pane",
		"  ^B h/j/k/l  Move pane focus",
		"",
		"  ^B ?      Toggle this help",
		"  ^B q      Quit",
		"",
		"  Press any key to close",
		"",
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(0, 2).
		Foreground(lipgloss.Color("#CCCCDD")).
		Background(lipgloss.Color("#111122")).
		Render(strings.Join(lines, "\n"))

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, box)
}

// --- Canvas helpers ---

type canvasLine []rune

func makeCanvas(w, h int) []canvasLine {
	c := make([]canvasLine, h)
	for i := range c {
		c[i] = make(canvasLine, w)
		for j := range c[i] {
			c[i][j] = ' '
		}
	}
	return c
}

func blitToCanvas(canvas []canvasLine, x, y int, s string) {
	lines := strings.Split(s, "\n")
	for dy, line := range lines {
		row := y + dy
		if row < 0 || row >= len(canvas) {
			continue
		}
		col := x
		for _, r := range line {
			if col >= len(canvas[row]) {
				break
			}
			canvas[row][col] = r
			col++
		}
	}
}

func canvasToString(canvas []canvasLine, w, h int) string {
	rows := make([]string, len(canvas))
	for i, line := range canvas {
		rows[i] = string(line)
	}
	return strings.Join(rows, "\n")
}

// fitContent truncates/pads content to exactly w×h cells.
func fitContent(content string, w, h int) string {
	// Strip ANSI sequences for simple fitting; real ANSI-aware fitting is done
	// by the vt.Screen which already has the right dimensions.
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		visible := visibleLen(line)
		if visible > w {
			lines[i] = truncateVisible(line, w)
		}
		// Don't pad — lipgloss handles width
	}
	// Trim to h lines
	if len(lines) > h {
		lines = lines[:h]
	}
	return strings.Join(lines, "\n")
}

// visibleLen returns the number of visible (non-ANSI-escape) characters.
func visibleLen(s string) int {
	inEsc := false
	count := 0
	for _, r := range s {
		if inEsc {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		if r == '\x1b' {
			inEsc = true
			continue
		}
		count++
	}
	return count
}

// truncateVisible truncates a string to n visible characters.
func truncateVisible(s string, n int) string {
	var out strings.Builder
	inEsc := false
	count := 0
	for _, r := range s {
		if inEsc {
			out.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEsc = false
			}
			continue
		}
		if r == '\x1b' {
			inEsc = true
			out.WriteRune(r)
			continue
		}
		if count >= n {
			break
		}
		out.WriteRune(r)
		count++
	}
	return out.String()
}

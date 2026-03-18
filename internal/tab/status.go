package tab

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/nowgnas/prism/internal/terminal"
)

// DeriveStatus returns the highest-priority status from a list of terminal statuses.
// Priority: Error > AI > Running > Idle
func DeriveStatus(terminals []*terminal.Model) terminal.StatusType {
	best := terminal.StatusIdle
	for _, t := range terminals {
		if t.Status > best {
			best = t.Status
		}
	}
	return best
}

// StatusIcon returns the display icon and color for a given status + tab color.
func StatusIcon(status terminal.StatusType, tabColor lipgloss.Color) (icon string, color lipgloss.Color) {
	switch status {
	case terminal.StatusIdle:
		return "○", "#555555"
	case terminal.StatusRunning:
		return "●", tabColor
	case terminal.StatusAI:
		return "◈", "#00FFFF"
	case terminal.StatusError:
		return "✕", "#FF3333"
	}
	return "○", "#555555"
}

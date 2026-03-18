package sidebar

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/nowgnas/prism/internal/tab"
	"github.com/nowgnas/prism/internal/terminal"
)

// Model holds the sidebar rendering state.
type Model struct {
	Width     int
	Height    int
	ActiveTab int
}

// New creates a new sidebar model.
func New(width, height int) Model {
	return Model{Width: width, Height: height}
}

// SetActive updates the currently selected tab index.
func (m *Model) SetActive(idx int) {
	m.ActiveTab = idx
}

// Resize updates dimensions.
func (m *Model) Resize(width, height int) {
	m.Width = width
	m.Height = height
}

// View renders the sidebar given the current tab list.
func (m *Model) View(tabs []*tab.Model) string {
	return RenderSidebar(m, tabs)
}

// statusStyle returns a styled status icon string.
func statusStyle(status terminal.StatusType, color lipgloss.Color) string {
	icon, c := tab.StatusIcon(status, color)
	return lipgloss.NewStyle().Foreground(c).Render(icon)
}

package tab

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nowgnas/prism/internal/layout"
	"github.com/nowgnas/prism/internal/messages"
	"github.com/nowgnas/prism/internal/terminal"
)

const maxTerminals = 6

// Model represents a single tab containing 1–6 terminal panes.
type Model struct {
	ID          int
	Name        string
	Color       lipgloss.Color
	Terminals   []*terminal.Model
	ActiveTerm  int
	Status      terminal.StatusType
	ZoomedPane  int  // -1 means no zoom
	Shell       string
	nextTermID  int
}

// New creates a new tab with one terminal already started.
func New(id int, shell string, cols, rows int) (*Model, tea.Cmd, error) {
	m := &Model{
		ID:         id,
		Color:      ColorForIndex(id),
		Shell:      shell,
		ZoomedPane: -1,
		nextTermID: 1,
	}

	cmd, err := m.addTerminal(cols, rows)
	if err != nil {
		return nil, nil, err
	}

	// Set default name from cwd
	m.Name = defaultName()

	return m, cmd, nil
}

// addTerminal adds a new terminal pane to this tab.
func (m *Model) addTerminal(cols, rows int) (tea.Cmd, error) {
	if len(m.Terminals) >= maxTerminals {
		return nil, fmt.Errorf("maximum of %d terminals per tab", maxTerminals)
	}
	id := m.nextTermID
	m.nextTermID++
	t := terminal.New(id, m.ID, cols, rows, m.Shell)
	readCmd, err := t.Start()
	if err != nil {
		return nil, err
	}
	m.Terminals = append(m.Terminals, t)
	return readCmd, nil
}

// AddPane adds a new pane and returns the read loop command.
func (m *Model) AddPane(cols, rows int) (tea.Cmd, error) {
	n := len(m.Terminals) + 1
	rects := layout.ComputeLayout(n, cols, rows)
	if len(rects) == 0 {
		return nil, fmt.Errorf("no rects computed")
	}
	// Resize existing terminals to new layout
	for i, t := range m.Terminals {
		inner := layout.InnerRect(rects[i])
		_ = t.Resize(inner.W, inner.H)
	}
	// Add new terminal with last rect dimensions
	newRect := layout.InnerRect(rects[len(rects)-1])
	return m.addTerminal(newRect.W, newRect.H)
}

// ClosePane removes a terminal pane by terminal ID.
func (m *Model) ClosePane(termID int) {
	for i, t := range m.Terminals {
		if t.ID == termID {
			if t.PTY != nil {
				t.PTY.Close()
			}
			m.Terminals = append(m.Terminals[:i], m.Terminals[i+1:]...)
			if m.ActiveTerm >= len(m.Terminals) && len(m.Terminals) > 0 {
				m.ActiveTerm = len(m.Terminals) - 1
			}
			return
		}
	}
}

// ActiveTerminal returns the currently focused terminal, or nil.
func (m *Model) ActiveTerminal() *terminal.Model {
	if m.ActiveTerm < 0 || m.ActiveTerm >= len(m.Terminals) {
		return nil
	}
	return m.Terminals[m.ActiveTerm]
}

// MoveFocus moves focus by delta (+1 / -1) within the pane list.
func (m *Model) MoveFocus(delta int) {
	if len(m.Terminals) == 0 {
		return
	}
	m.ActiveTerm = (m.ActiveTerm + delta + len(m.Terminals)) % len(m.Terminals)
	for i, t := range m.Terminals {
		t.Focused = (i == m.ActiveTerm)
	}
}

// FocusDirection moves focus in a vim-like direction (h/j/k/l).
func (m *Model) FocusDirection(dir string) {
	switch dir {
	case "h", "k":
		m.MoveFocus(-1)
	case "l", "j":
		m.MoveFocus(+1)
	}
}

// ToggleZoom zooms or un-zooms the active pane.
func (m *Model) ToggleZoom() {
	if m.ZoomedPane >= 0 {
		m.ZoomedPane = -1
	} else {
		m.ZoomedPane = m.ActiveTerm
	}
}

// Update processes messages directed at this tab.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case messages.PTYOutputMsg:
		if msg.TabID != m.ID {
			break
		}
		t := m.findTerminal(msg.TerminalID)
		if t != nil {
			cmd := t.HandleOutput(msg.Data)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		m.Status = DeriveStatus(m.Terminals)

	case messages.PTYExitMsg:
		if msg.TabID != m.ID {
			break
		}
		t := m.findTerminal(msg.TerminalID)
		if t != nil {
			t.HandleExit()
		}
		m.Status = DeriveStatus(m.Terminals)

	case messages.TickMsg:
		for _, t := range m.Terminals {
			t.UpdateStatusFromPoll()
		}
		m.Status = DeriveStatus(m.Terminals)
	}

	return tea.Batch(cmds...)
}

// Resize recomputes layout and resizes all terminals.
func (m *Model) Resize(totalCols, totalRows int) {
	n := len(m.Terminals)
	if n == 0 {
		return
	}
	var rects []layout.Rect
	if m.ZoomedPane >= 0 {
		rects = layout.ComputeLayout(1, totalCols, totalRows)
		for i, t := range m.Terminals {
			if i == m.ZoomedPane {
				inner := layout.InnerRect(rects[0])
				_ = t.Resize(inner.W, inner.H)
			}
		}
		return
	}
	rects = layout.ComputeLayout(n, totalCols, totalRows)
	for i, t := range m.Terminals {
		if i < len(rects) {
			inner := layout.InnerRect(rects[i])
			_ = t.Resize(inner.W, inner.H)
		}
	}
}

// findTerminal returns the terminal model by ID.
func (m *Model) findTerminal(id int) *terminal.Model {
	for _, t := range m.Terminals {
		if t.ID == id {
			return t
		}
	}
	return nil
}

// defaultName returns the base name of the current working directory.
func defaultName() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "terminal"
	}
	base := filepath.Base(cwd)
	if base == "" || base == "." || base == "/" {
		return "terminal"
	}
	return strings.ToLower(base)
}

package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nowgnas/prism/internal/messages"
)

// Update is the BubbleTea update function for AppModel.
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Sidebar.Resize(m.Config.SidebarWidth, m.Height)
		if len(m.Tabs) == 0 {
			// First resize — create the initial tab now that we have dimensions
			cmd := m.addTab()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		} else {
			tabW, tabH := m.tabAreaDimensions()
			for _, t := range m.Tabs {
				t.Resize(tabW, tabH)
			}
		}

	case spinner.TickMsg:
		var spCmd tea.Cmd
		m.Spinner, spCmd = m.Spinner.Update(msg)
		cmds = append(cmds, spCmd)

	case messages.TickMsg:
		for _, t := range m.Tabs {
			t.Update(msg)
		}
		cmds = append(cmds, m.pollTick())

	case messages.PTYOutputMsg:
		for _, t := range m.Tabs {
			if t.ID == msg.TabID {
				cmd := t.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				break
			}
		}

	case messages.PTYExitMsg:
		for _, t := range m.Tabs {
			if t.ID == msg.TabID {
				t.Update(msg)
				break
			}
		}

	case tea.KeyMsg:
		cmd := m.handleKey(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleKey processes a key message based on the current mode.
func (m *AppModel) handleKey(msg tea.KeyMsg) tea.Cmd {
	switch m.Mode {
	case ModeRename:
		return m.handleRenameKey(msg)
	case ModeHelp:
		m.Mode = ModeNormal
		m.showHelp = false
		return nil
	case ModePrefix:
		return m.handlePrefixKey(msg)
	default:
		return m.handleNormalKey(msg)
	}
}

// handleNormalKey handles keys in normal (pass-through) mode.
func (m *AppModel) handleNormalKey(msg tea.KeyMsg) tea.Cmd {
	// Ctrl+B activates prefix mode
	if msg.Type == tea.KeyCtrlB {
		m.Mode = ModePrefix
		return nil
	}
	// All other keys go to the active terminal
	t := m.activeTab()
	if t == nil {
		return nil
	}
	term := t.ActiveTerminal()
	if term == nil {
		return nil
	}
	term.Write([]byte(msg.String()))
	return nil
}

// handlePrefixKey handles keys after Ctrl+B prefix.
func (m *AppModel) handlePrefixKey(msg tea.KeyMsg) tea.Cmd {
	m.Mode = ModeNormal // reset after handling one key

	keyStr := msg.String()

	switch {
	case keyStr == "ctrl+b":
		// Double Ctrl+B — pass one Ctrl+B to terminal
		t := m.activeTab()
		if t != nil {
			if term := t.ActiveTerminal(); term != nil {
				term.Write([]byte{0x02}) // Ctrl+B byte
			}
		}

	case keyStr == "n":
		return m.addTab()

	case keyStr == "w":
		idx := m.ActiveTab
		m.removeTab(idx)
		if len(m.Tabs) == 0 {
			return tea.Quit
		}

	case keyStr == "q":
		return tea.Quit

	case keyStr == ",":
		if t := m.activeTab(); t != nil {
			m.Mode = ModeRename
			m.Rename.SetValue(t.Name)
			m.Rename.Focus()
		}

	case keyStr == "tab":
		if len(m.Tabs) > 0 {
			m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
			m.Sidebar.SetActive(m.ActiveTab)
		}

	case keyStr == "\"" || keyStr == "%":
		// Both split types add a pane (layout auto-adjusts)
		if t := m.activeTab(); t != nil {
			tabW, tabH := m.tabAreaDimensions()
			cmd, err := t.AddPane(tabW, tabH)
			if err == nil && cmd != nil {
				return cmd
			}
		}

	case keyStr == "x":
		if t := m.activeTab(); t != nil {
			if term := t.ActiveTerminal(); term != nil {
				t.ClosePane(term.ID)
				if len(t.Terminals) == 0 {
					m.removeTab(m.ActiveTab)
					if len(m.Tabs) == 0 {
						return tea.Quit
					}
				}
			}
		}

	case keyStr == "z":
		if t := m.activeTab(); t != nil {
			t.ToggleZoom()
		}

	case keyStr == "h", keyStr == "j", keyStr == "k", keyStr == "l":
		if t := m.activeTab(); t != nil {
			t.FocusDirection(keyStr)
		}

	case keyStr == "?":
		m.Mode = ModeHelp
		m.showHelp = true

	case len(keyStr) == 1 && keyStr[0] >= '1' && keyStr[0] <= '9':
		idx := int(keyStr[0] - '1')
		if idx < len(m.Tabs) {
			m.ActiveTab = idx
			m.Sidebar.SetActive(m.ActiveTab)
		}
	}

	return nil
}

// handleRenameKey handles input in tab rename mode.
func (m *AppModel) handleRenameKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEnter:
		if t := m.activeTab(); t != nil {
			name := strings.TrimSpace(m.Rename.Value())
			if name != "" {
				t.Name = name
			}
		}
		m.Rename.Blur()
		m.Mode = ModeNormal
		return nil
	case tea.KeyEsc:
		m.Rename.Blur()
		m.Mode = ModeNormal
		return nil
	}
	var cmd tea.Cmd
	m.Rename, cmd = m.Rename.Update(msg)
	return cmd
}


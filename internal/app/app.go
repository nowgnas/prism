package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nowgnas/prism/internal/config"
	"github.com/nowgnas/prism/internal/messages"
	"github.com/nowgnas/prism/internal/sidebar"
	"github.com/nowgnas/prism/internal/tab"
)

// Mode represents the current input mode of the app.
type Mode int

const (
	ModeNormal  Mode = iota // forwarding input to active terminal
	ModePrefix              // Ctrl+B prefix was pressed, waiting for command
	ModeRename              // user is typing a new tab name
	ModeHelp                // help overlay visible
)

// AppModel is the root BubbleTea model.
type AppModel struct {
	Tabs      []*tab.Model
	ActiveTab int
	Sidebar   sidebar.Model
	Width     int
	Height    int
	Config    config.Config
	KeyMap    KeyMap
	Mode      Mode
	Spinner   spinner.Model
	Rename    textinput.Model
	nextTabID int
	showHelp  bool
}

// New creates an initial AppModel with one tab.
func New(cfg config.Config) (*AppModel, error) {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "tab name"
	ti.CharLimit = 40

	m := &AppModel{
		Config:    cfg,
		KeyMap:    DefaultKeyMap(),
		Spinner:   sp,
		Rename:    ti,
		nextTabID: 1,
		Mode:      ModeNormal,
	}

	// We don't know dimensions yet; use sensible defaults until WindowSizeMsg
	return m, nil
}

// Init is called once by BubbleTea on startup.
func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.pollTick(),
		// First tab will be created on first WindowSizeMsg
	)
}

// pollTick returns a command that fires a TickMsg after the poll interval.
func (m *AppModel) pollTick() tea.Cmd {
	interval := time.Duration(m.Config.PollInterval) * time.Millisecond
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return messages.TickMsg{}
	})
}

// activeTab returns the currently active tab model, or nil.
func (m *AppModel) activeTab() *tab.Model {
	if m.ActiveTab < 0 || m.ActiveTab >= len(m.Tabs) {
		return nil
	}
	return m.Tabs[m.ActiveTab]
}

// addTab creates a new tab with the current terminal dimensions.
func (m *AppModel) addTab() tea.Cmd {
	tabW, tabH := m.tabAreaDimensions()
	id := m.nextTabID
	m.nextTabID++
	t, readCmd, err := tab.New(id, m.Config.Shell, tabW, tabH)
	if err != nil {
		// Non-fatal: log to stderr in real usage
		fmt.Printf("error creating tab: %v\n", err)
		return nil
	}
	m.Tabs = append(m.Tabs, t)
	m.ActiveTab = len(m.Tabs) - 1
	m.Sidebar.SetActive(m.ActiveTab)
	return readCmd
}

// removeTab closes and removes a tab by index.
func (m *AppModel) removeTab(idx int) {
	if idx < 0 || idx >= len(m.Tabs) {
		return
	}
	t := m.Tabs[idx]
	for _, term := range t.Terminals {
		if term.PTY != nil {
			term.PTY.Close()
		}
	}
	m.Tabs = append(m.Tabs[:idx], m.Tabs[idx+1:]...)
	if len(m.Tabs) == 0 {
		return
	}
	if m.ActiveTab >= len(m.Tabs) {
		m.ActiveTab = len(m.Tabs) - 1
	}
	m.Sidebar.SetActive(m.ActiveTab)
}

// tabAreaDimensions returns the W×H available for terminal panes.
func (m *AppModel) tabAreaDimensions() (int, int) {
	w := m.Width - m.Config.SidebarWidth
	if w < 10 {
		w = 10
	}
	h := m.Height
	if h < 5 {
		h = 5
	}
	return w, h
}

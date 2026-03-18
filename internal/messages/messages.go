package messages

import "os"

// PTYOutputMsg carries raw bytes read from a PTY.
type PTYOutputMsg struct {
	TerminalID int
	TabID      int
	Data       []byte
}

// PTYExitMsg signals that a PTY process has exited.
type PTYExitMsg struct {
	TerminalID int
	TabID      int
}

// TerminalResizeMsg requests a PTY resize.
type TerminalResizeMsg struct {
	TerminalID int
	TabID      int
	Cols, Rows int
}

// StatusUpdateMsg signals a terminal status change.
type StatusUpdateMsg struct {
	TerminalID int
	TabID      int
	Status     int // maps to StatusType
}

// TabRenameMsg carries a new name for the active tab.
type TabRenameMsg struct {
	TabID int
	Name  string
}

// NewTabMsg requests a new tab.
type NewTabMsg struct{}

// CloseTabMsg requests closing a tab.
type CloseTabMsg struct {
	TabID int
}

// NewPaneMsg requests adding a pane to the active tab.
type NewPaneMsg struct {
	TabID int
}

// ClosePaneMsg requests closing a pane.
type ClosePaneMsg struct {
	TabID      int
	TerminalID int
}

// ZoomPaneMsg toggles zoom on the active pane.
type ZoomPaneMsg struct {
	TabID      int
	TerminalID int
}

// FocusPaneMsg moves focus to a specific pane.
type FocusPaneMsg struct {
	TabID      int
	TerminalID int
}

// ProcessStartedMsg signals that a child process started in a terminal.
type ProcessStartedMsg struct {
	TerminalID int
	TabID      int
	PID        int
	Process    *os.Process
}

// TickMsg drives periodic polling for process status.
type TickMsg struct{}

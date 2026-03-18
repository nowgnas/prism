package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nowgnas/prism/internal/messages"
	"github.com/nowgnas/prism/pkg/vt"
)

// StatusType represents the activity state of a terminal.
type StatusType int

const (
	StatusIdle    StatusType = iota // ○ grey
	StatusRunning                   // ● tab color (spinner)
	StatusAI                        // ◈ cyan (spinner)
	StatusError                     // ✕ red
)

// Model represents a single PTY terminal pane.
type Model struct {
	ID      int
	TabID   int
	PTY     *os.File
	Cmd     *exec.Cmd
	Screen  *vt.Screen
	Cols    int
	Rows    int
	Focused bool
	Status  StatusType
	Shell   string

	// Process tracking
	ChildPID     int
	lastPollTime time.Time
}

// New creates a new terminal model. Call Start() to actually launch the PTY.
func New(id, tabID, cols, rows int, shell string) *Model {
	return &Model{
		ID:     id,
		TabID:  tabID,
		Cols:   cols,
		Rows:   rows,
		Shell:  shell,
		Screen: vt.NewScreen(cols, rows),
		Status: StatusIdle,
	}
}

// Start launches the PTY subprocess.
func (m *Model) Start() (tea.Cmd, error) {
	ptmx, cmd, err := Start(m.Shell, m.Cols, m.Rows, nil)
	if err != nil {
		return nil, fmt.Errorf("terminal %d start: %w", m.ID, err)
	}
	m.PTY = ptmx
	m.Cmd = cmd
	m.ChildPID = cmd.Process.Pid
	return ReadLoop(ptmx, m.TabID, m.ID), nil
}

// HandleOutput processes raw PTY output, updates the screen, and detects status.
func (m *Model) HandleOutput(data []byte) tea.Cmd {
	// Parse OSC 133 events
	event := ParseOSC133(data)
	switch event {
	case OSC133CommandStart:
		m.Status = StatusRunning
	case OSC133CommandEnd:
		m.Status = StatusIdle
	}

	// Write to virtual screen
	m.Screen.Write(data)

	// Schedule next read
	return ReadLoop(m.PTY, m.TabID, m.ID)
}

// HandleExit marks the terminal as errored/exited.
func (m *Model) HandleExit() {
	m.Status = StatusError
}

// Resize adjusts the terminal dimensions.
func (m *Model) Resize(cols, rows int) error {
	m.Cols = cols
	m.Rows = rows
	m.Screen.Resize(cols, rows)
	if m.PTY != nil {
		return Resize(m.PTY, cols, rows)
	}
	return nil
}

// Write sends input bytes to the PTY (e.g., keyboard input).
func (m *Model) Write(p []byte) {
	if m.PTY != nil {
		m.PTY.Write(p)
	}
}

// PollStatus polls the child process tree to detect running/AI processes.
func (m *Model) PollStatus() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond)
		return messages.TickMsg{}
	}
}

// UpdateStatusFromPoll checks child processes of the shell and updates status.
func (m *Model) UpdateStatusFromPoll() {
	if m.ChildPID <= 0 {
		return
	}
	children := getChildProcesses(m.ChildPID)
	if len(children) == 0 {
		if m.Status == StatusRunning || m.Status == StatusAI {
			m.Status = StatusIdle
		}
		return
	}
	// Check for AI tools first (higher priority)
	for _, name := range children {
		if IsAIProcess(name) {
			m.Status = StatusAI
			return
		}
	}
	m.Status = StatusRunning
}

// getChildProcesses returns the names of direct children of the given PID.
func getChildProcesses(pid int) []string {
	// Read /proc on Linux or use pgrep on macOS
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/task/%d/children", pid, pid))
	if err == nil {
		// Linux: parse child PIDs
		var names []string
		for _, p := range strings.Fields(string(data)) {
			childPID, _ := strconv.Atoi(p)
			if childPID > 0 {
				if name := processName(childPID); name != "" {
					names = append(names, name)
				}
			}
		}
		return names
	}
	// macOS: use pgrep
	return pgrepChildren(pid)
}

func processName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func pgrepChildren(ppid int) []string {
	// Use ps to get children on macOS
	cmd := exec.Command("ps", "-o", "pid=,comm=", fmt.Sprintf("--ppid=%d", ppid))
	out, err := cmd.Output()
	if err != nil {
		// Try macOS syntax
		cmd2 := exec.Command("ps", "-ax", "-o", "ppid=,comm=")
		out2, err2 := cmd2.Output()
		if err2 != nil {
			return nil
		}
		var names []string
		for _, line := range strings.Split(string(out2), "\n") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				parentPID, _ := strconv.Atoi(parts[0])
				if parentPID == ppid {
					names = append(names, parts[1])
				}
			}
		}
		return names
	}
	var names []string
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			names = append(names, parts[1])
		}
	}
	return names
}

package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/creack/pty"
)

// Start launches a shell in a new PTY with the given dimensions.
// Returns the PTY file and the started Cmd.
func Start(shell string, cols, rows int, env []string) (*os.File, *exec.Cmd, error) {
	if shell == "" {
		shell = detectShell()
	}

	cmd := exec.Command(shell)
	cmd.Env = buildEnv(env, cols, rows)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	size := &pty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	}

	ptmx, err := pty.StartWithSize(cmd, size)
	if err != nil {
		return nil, nil, fmt.Errorf("pty.StartWithSize: %w", err)
	}

	return ptmx, cmd, nil
}

// Resize sends a SIGWINCH to the PTY with the new dimensions.
func Resize(ptmx *os.File, cols, rows int) error {
	return pty.Setsize(ptmx, &pty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	})
}

func detectShell() string {
	// Try $SHELL env var first
	if s := os.Getenv("SHELL"); s != "" {
		return s
	}
	// Fallback order
	for _, sh := range []string{"/bin/zsh", "/bin/bash", "/bin/sh"} {
		if _, err := os.Stat(sh); err == nil {
			return sh
		}
	}
	return "/bin/sh"
}

func buildEnv(extra []string, cols, rows int) []string {
	env := os.Environ()
	// Override TERM, COLUMNS, LINES
	filtered := make([]string, 0, len(env)+10)
	for _, e := range env {
		if strings.HasPrefix(e, "TERM=") ||
			strings.HasPrefix(e, "COLUMNS=") ||
			strings.HasPrefix(e, "LINES=") {
			continue
		}
		filtered = append(filtered, e)
	}
	filtered = append(filtered,
		"TERM=xterm-256color",
		fmt.Sprintf("COLUMNS=%d", cols),
		fmt.Sprintf("LINES=%d", rows),
	)
	filtered = append(filtered, extra...)
	return filtered
}

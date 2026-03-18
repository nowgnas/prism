//go:build windows

package terminal

import (
	"fmt"
	"os"
	"os/exec"
)

// Start is not supported on Windows — PTY requires Unix.
func Start(shell string, cols, rows int, env []string) (*os.File, *exec.Cmd, error) {
	return nil, nil, fmt.Errorf("prism does not support Windows PTY")
}

// Resize is a no-op on Windows.
func Resize(ptmx *os.File, cols, rows int) error {
	return nil
}

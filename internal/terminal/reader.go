package terminal

import (
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nowgnas/prism/internal/messages"
)

const readBufSize = 4096

// ReadLoop returns a tea.Cmd that continuously reads from the PTY
// and sends PTYOutputMsg messages. It exits when the PTY is closed.
func ReadLoop(ptmx *os.File, tabID, termID int) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, readBufSize)
		n, err := ptmx.Read(buf)
		if err != nil {
			if err == io.EOF {
				return messages.PTYExitMsg{TabID: tabID, TerminalID: termID}
			}
			return messages.PTYExitMsg{TabID: tabID, TerminalID: termID}
		}
		data := make([]byte, n)
		copy(data, buf[:n])
		return messages.PTYOutputMsg{
			TabID:      tabID,
			TerminalID: termID,
			Data:       data,
		}
	}
}

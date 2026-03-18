package vt

import (
	"bytes"
	"strings"
	"sync"
	"unicode/utf8"
)

// Cell represents a single character cell in the terminal grid.
type Cell struct {
	Char       rune
	FG, BG     Color
	Bold       bool
	Underline  bool
	Reverse    bool
	Dim        bool
	Italic     bool
	Blink      bool
}

// Color represents an ANSI terminal color.
type Color struct {
	Mode  ColorMode
	Index uint8   // for 256-color
	R, G, B uint8 // for truecolor
}

// ColorMode distinguishes color types.
type ColorMode int

const (
	ColorDefault ColorMode = iota
	ColorANSI
	Color256
	ColorRGB
)

// DefaultCell returns a blank cell with default colors.
func DefaultCell() Cell {
	return Cell{Char: ' '}
}

// Screen is a VT100/ANSI terminal state machine with a fixed-size cell grid.
type Screen struct {
	mu     sync.Mutex
	cols   int
	rows   int
	cells  [][]Cell
	cursor Cursor
	saved  Cursor
	attrs  Cell // current rendering attributes

	// Scroll region
	scrollTop    int
	scrollBottom int

	// OSC/escape parsing state
	inEscape  bool
	inCSI     bool
	inOSC     bool
	escBuf    bytes.Buffer
	oscBuf    bytes.Buffer

	// OSC 133 command tracking
	CommandStart bool // received OSC 133 ; A (prompt start)
	CommandEnd   bool // received OSC 133 ; C (command start / end)

	// Alt screen
	altScreen   [][]Cell
	usingAlt    bool
	mainCells   [][]Cell
	mainCursor  Cursor
}

// Cursor tracks the current position and state.
type Cursor struct {
	X, Y    int
	Visible bool
}

// NewScreen creates a new Screen with the given dimensions.
func NewScreen(cols, rows int) *Screen {
	s := &Screen{
		cols:         cols,
		rows:         rows,
		scrollTop:    0,
		scrollBottom: rows - 1,
		cursor:       Cursor{Visible: true},
	}
	s.cells = makeGrid(cols, rows)
	return s
}

func makeGrid(cols, rows int) [][]Cell {
	g := make([][]Cell, rows)
	for i := range g {
		g[i] = make([]Cell, cols)
		for j := range g[i] {
			g[i][j] = DefaultCell()
		}
	}
	return g
}

// Resize resizes the screen, preserving as much content as possible.
func (s *Screen) Resize(cols, rows int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newCells := makeGrid(cols, rows)
	for r := 0; r < rows && r < s.rows; r++ {
		for c := 0; c < cols && c < s.cols; c++ {
			newCells[r][c] = s.cells[r][c]
		}
	}
	s.cells = newCells
	s.cols = cols
	s.rows = rows
	s.scrollTop = 0
	s.scrollBottom = rows - 1
	if s.cursor.Y >= rows {
		s.cursor.Y = rows - 1
	}
	if s.cursor.X >= cols {
		s.cursor.X = cols - 1
	}
}

// Write processes a raw byte slice (may contain ANSI escape sequences).
func (s *Screen) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.process(p)
	return len(p), nil
}

// Render returns the current screen contents as a plain string (no ANSI).
func (s *Screen) Render() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sb strings.Builder
	for r := 0; r < s.rows; r++ {
		for c := 0; c < s.cols; c++ {
			ch := s.cells[r][c].Char
			if ch == 0 {
				ch = ' '
			}
			sb.WriteRune(ch)
		}
		if r < s.rows-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// RenderANSI returns the current screen with ANSI color codes reconstructed.
func (s *Screen) RenderANSI() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sb strings.Builder
	prevCell := DefaultCell()

	for r := 0; r < s.rows; r++ {
		for c := 0; c < s.cols; c++ {
			cell := s.cells[r][c]
			if cellAttrsEqual(cell, prevCell) {
				// same attributes — just write char
			} else {
				sb.WriteString(cellToANSI(cell))
				prevCell = cell
			}
			ch := cell.Char
			if ch == 0 {
				ch = ' '
			}
			sb.WriteRune(ch)
		}
		if r < s.rows-1 {
			sb.WriteString("\033[0m\n")
			prevCell = DefaultCell()
		}
	}
	sb.WriteString("\033[0m")
	return sb.String()
}

// CursorPos returns the current cursor position.
func (s *Screen) CursorPos() (x, y int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cursor.X, s.cursor.Y
}

// Cols and Rows return screen dimensions.
func (s *Screen) Cols() int { return s.cols }
func (s *Screen) Rows() int { return s.rows }

// --- Internal processing ---

func (s *Screen) process(data []byte) {
	i := 0
	for i < len(data) {
		b := data[i]

		if s.inOSC {
			if b == '\x07' || (b == '\\' && s.escBuf.Len() > 0 && s.escBuf.Bytes()[s.escBuf.Len()-1] == '\x1b') {
				s.handleOSC(s.oscBuf.String())
				s.oscBuf.Reset()
				s.escBuf.Reset()
				s.inOSC = false
				i++
				continue
			}
			s.oscBuf.WriteByte(b)
			i++
			continue
		}

		if s.inCSI {
			s.escBuf.WriteByte(b)
			if b >= 0x40 && b <= 0x7E {
				s.handleCSI(s.escBuf.String())
				s.escBuf.Reset()
				s.inCSI = false
			}
			i++
			continue
		}

		if s.inEscape {
			s.inEscape = false
			switch b {
			case '[':
				s.inCSI = true
				s.escBuf.Reset()
			case ']':
				s.inOSC = true
				s.oscBuf.Reset()
				s.escBuf.Reset()
			case '7':
				s.saved = s.cursor
			case '8':
				s.cursor = s.saved
			case 'M': // reverse index
				if s.cursor.Y == s.scrollTop {
					s.scrollDown(1)
				} else if s.cursor.Y > 0 {
					s.cursor.Y--
				}
			case 'c': // full reset
				s.fullReset()
			case '(':
				i++ // skip charset
			}
			i++
			continue
		}

		switch b {
		case '\x1b':
			s.inEscape = true
			i++
		case '\r':
			s.cursor.X = 0
			i++
		case '\n', '\x0b', '\x0c':
			s.lineFeed()
			i++
		case '\x08': // backspace
			if s.cursor.X > 0 {
				s.cursor.X--
			}
			i++
		case '\x09': // tab
			s.cursor.X = (s.cursor.X/8 + 1) * 8
			if s.cursor.X >= s.cols {
				s.cursor.X = s.cols - 1
			}
			i++
		case '\x07': // bell — ignore
			i++
		case '\x0e', '\x0f': // shift in/out — ignore
			i++
		default:
			// Try to decode UTF-8 rune
			r, size := utf8.DecodeRune(data[i:])
			if r == utf8.RuneError && size == 1 {
				// Bad byte, skip
				i++
				continue
			}
			s.putChar(r)
			i += size
		}
	}
}

func (s *Screen) putChar(r rune) {
	if s.cursor.Y < 0 || s.cursor.Y >= s.rows {
		return
	}
	if s.cursor.X >= s.cols {
		// auto-wrap
		s.cursor.X = 0
		s.lineFeed()
	}
	cell := s.attrs
	cell.Char = r
	s.cells[s.cursor.Y][s.cursor.X] = cell
	s.cursor.X++
}

func (s *Screen) lineFeed() {
	if s.cursor.Y == s.scrollBottom {
		s.scrollUp(1)
	} else if s.cursor.Y < s.rows-1 {
		s.cursor.Y++
	}
}

func (s *Screen) scrollUp(n int) {
	for i := 0; i < n; i++ {
		copy(s.cells[s.scrollTop:s.scrollBottom+1], s.cells[s.scrollTop+1:s.scrollBottom+2])
		s.cells[s.scrollBottom] = make([]Cell, s.cols)
		for j := range s.cells[s.scrollBottom] {
			s.cells[s.scrollBottom][j] = DefaultCell()
		}
	}
}

func (s *Screen) scrollDown(n int) {
	for i := 0; i < n; i++ {
		copy(s.cells[s.scrollTop+1:s.scrollBottom+2], s.cells[s.scrollTop:s.scrollBottom+1])
		s.cells[s.scrollTop] = make([]Cell, s.cols)
		for j := range s.cells[s.scrollTop] {
			s.cells[s.scrollTop][j] = DefaultCell()
		}
	}
}

func (s *Screen) fullReset() {
	s.cells = makeGrid(s.cols, s.rows)
	s.cursor = Cursor{Visible: true}
	s.saved = Cursor{}
	s.attrs = DefaultCell()
	s.scrollTop = 0
	s.scrollBottom = s.rows - 1
}

// handleCSI parses a CSI escape sequence (without the leading ESC [).
func (s *Screen) handleCSI(seq string) {
	if len(seq) == 0 {
		return
	}
	final := seq[len(seq)-1]
	params := seq[:len(seq)-1]

	nums := parseParams(params)
	getParam := func(idx, def int) int {
		if idx < len(nums) && nums[idx] != 0 {
			return nums[idx]
		}
		return def
	}

	switch final {
	case 'A': // cursor up
		n := getParam(0, 1)
		s.cursor.Y -= n
		if s.cursor.Y < 0 {
			s.cursor.Y = 0
		}
	case 'B': // cursor down
		n := getParam(0, 1)
		s.cursor.Y += n
		if s.cursor.Y >= s.rows {
			s.cursor.Y = s.rows - 1
		}
	case 'C': // cursor forward
		n := getParam(0, 1)
		s.cursor.X += n
		if s.cursor.X >= s.cols {
			s.cursor.X = s.cols - 1
		}
	case 'D': // cursor back
		n := getParam(0, 1)
		s.cursor.X -= n
		if s.cursor.X < 0 {
			s.cursor.X = 0
		}
	case 'E': // cursor next line
		n := getParam(0, 1)
		s.cursor.Y += n
		s.cursor.X = 0
		if s.cursor.Y >= s.rows {
			s.cursor.Y = s.rows - 1
		}
	case 'F': // cursor prev line
		n := getParam(0, 1)
		s.cursor.Y -= n
		s.cursor.X = 0
		if s.cursor.Y < 0 {
			s.cursor.Y = 0
		}
	case 'G': // cursor horizontal absolute
		n := getParam(0, 1)
		s.cursor.X = n - 1
		if s.cursor.X < 0 {
			s.cursor.X = 0
		}
	case 'H', 'f': // cursor position
		row := getParam(0, 1)
		col := getParam(1, 1)
		s.cursor.Y = row - 1
		s.cursor.X = col - 1
		if s.cursor.Y < 0 {
			s.cursor.Y = 0
		}
		if s.cursor.Y >= s.rows {
			s.cursor.Y = s.rows - 1
		}
		if s.cursor.X < 0 {
			s.cursor.X = 0
		}
		if s.cursor.X >= s.cols {
			s.cursor.X = s.cols - 1
		}
	case 'J': // erase in display
		n := getParam(0, 0)
		switch n {
		case 0: // erase from cursor to end
			s.eraseFromCursor()
		case 1: // erase from start to cursor
			s.eraseToCursor()
		case 2, 3: // erase entire screen
			s.cells = makeGrid(s.cols, s.rows)
		}
	case 'K': // erase in line
		n := getParam(0, 0)
		switch n {
		case 0: // erase from cursor to end of line
			s.eraseLineRight()
		case 1: // erase from start of line to cursor
			s.eraseLineLeft()
		case 2: // erase entire line
			s.eraseLine(s.cursor.Y)
		}
	case 'L': // insert lines
		n := getParam(0, 1)
		s.insertLines(n)
	case 'M': // delete lines
		n := getParam(0, 1)
		s.deleteLines(n)
	case 'P': // delete characters
		n := getParam(0, 1)
		s.deleteChars(n)
	case 'S': // scroll up
		n := getParam(0, 1)
		s.scrollUp(n)
	case 'T': // scroll down
		n := getParam(0, 1)
		s.scrollDown(n)
	case 'X': // erase characters
		n := getParam(0, 1)
		for i := 0; i < n && s.cursor.X+i < s.cols; i++ {
			s.cells[s.cursor.Y][s.cursor.X+i] = DefaultCell()
		}
	case '@': // insert characters
		n := getParam(0, 1)
		s.insertChars(n)
	case 'd': // line position absolute
		n := getParam(0, 1)
		s.cursor.Y = n - 1
		if s.cursor.Y < 0 {
			s.cursor.Y = 0
		}
		if s.cursor.Y >= s.rows {
			s.cursor.Y = s.rows - 1
		}
	case 'm': // select graphic rendition
		s.handleSGR(nums)
	case 'h', 'l': // mode set/reset (handle private modes)
		if strings.HasPrefix(params, "?") {
			s.handlePrivateMode(params[1:], final == 'h')
		}
	case 'r': // set scroll region
		top := getParam(0, 1)
		bottom := getParam(1, s.rows)
		s.scrollTop = top - 1
		s.scrollBottom = bottom - 1
		s.cursor.X = 0
		s.cursor.Y = 0
	case 's': // save cursor
		s.saved = s.cursor
	case 'u': // restore cursor
		s.cursor = s.saved
	}
}

func (s *Screen) handlePrivateMode(params string, set bool) {
	for _, p := range parseParams(params) {
		switch p {
		case 1049: // alt screen
			if set {
				s.mainCells = s.cells
				s.mainCursor = s.cursor
				s.cells = makeGrid(s.cols, s.rows)
				s.cursor = Cursor{Visible: true}
				s.usingAlt = true
			} else {
				if s.mainCells != nil {
					s.cells = s.mainCells
					s.cursor = s.mainCursor
				}
				s.usingAlt = false
			}
		case 25: // cursor visibility
			s.cursor.Visible = set
		}
	}
}

func (s *Screen) handleSGR(nums []int) {
	if len(nums) == 0 {
		nums = []int{0}
	}
	for i := 0; i < len(nums); i++ {
		n := nums[i]
		switch {
		case n == 0:
			s.attrs = DefaultCell()
		case n == 1:
			s.attrs.Bold = true
		case n == 2:
			s.attrs.Dim = true
		case n == 3:
			s.attrs.Italic = true
		case n == 4:
			s.attrs.Underline = true
		case n == 5:
			s.attrs.Blink = true
		case n == 7:
			s.attrs.Reverse = true
		case n == 22:
			s.attrs.Bold = false
			s.attrs.Dim = false
		case n == 23:
			s.attrs.Italic = false
		case n == 24:
			s.attrs.Underline = false
		case n == 25:
			s.attrs.Blink = false
		case n == 27:
			s.attrs.Reverse = false
		case n >= 30 && n <= 37:
			s.attrs.FG = Color{Mode: ColorANSI, Index: uint8(n - 30)}
		case n == 38: // extended FG
			if i+2 < len(nums) && nums[i+1] == 5 {
				s.attrs.FG = Color{Mode: Color256, Index: uint8(nums[i+2])}
				i += 2
			} else if i+4 < len(nums) && nums[i+1] == 2 {
				s.attrs.FG = Color{Mode: ColorRGB, R: uint8(nums[i+2]), G: uint8(nums[i+3]), B: uint8(nums[i+4])}
				i += 4
			}
		case n == 39:
			s.attrs.FG = Color{}
		case n >= 40 && n <= 47:
			s.attrs.BG = Color{Mode: ColorANSI, Index: uint8(n - 40)}
		case n == 48: // extended BG
			if i+2 < len(nums) && nums[i+1] == 5 {
				s.attrs.BG = Color{Mode: Color256, Index: uint8(nums[i+2])}
				i += 2
			} else if i+4 < len(nums) && nums[i+1] == 2 {
				s.attrs.BG = Color{Mode: ColorRGB, R: uint8(nums[i+2]), G: uint8(nums[i+3]), B: uint8(nums[i+4])}
				i += 4
			}
		case n == 49:
			s.attrs.BG = Color{}
		case n >= 90 && n <= 97: // bright FG
			s.attrs.FG = Color{Mode: ColorANSI, Index: uint8(n - 90 + 8)}
		case n >= 100 && n <= 107: // bright BG
			s.attrs.BG = Color{Mode: ColorANSI, Index: uint8(n - 100 + 8)}
		}
	}
}

func (s *Screen) handleOSC(seq string) {
	// OSC 133 shell integration
	if strings.HasPrefix(seq, "133;A") {
		s.CommandStart = true
		s.CommandEnd = false
	} else if strings.HasPrefix(seq, "133;C") {
		s.CommandEnd = true
		s.CommandStart = false
	}
	// OSC 0 / OSC 2 — window title (ignore or store if needed)
}

// --- Erase helpers ---

func (s *Screen) eraseFromCursor() {
	s.eraseLineRight()
	for r := s.cursor.Y + 1; r < s.rows; r++ {
		s.eraseLine(r)
	}
}

func (s *Screen) eraseToCursor() {
	for r := 0; r < s.cursor.Y; r++ {
		s.eraseLine(r)
	}
	s.eraseLineLeft()
}

func (s *Screen) eraseLineRight() {
	for c := s.cursor.X; c < s.cols; c++ {
		s.cells[s.cursor.Y][c] = DefaultCell()
	}
}

func (s *Screen) eraseLineLeft() {
	for c := 0; c <= s.cursor.X; c++ {
		s.cells[s.cursor.Y][c] = DefaultCell()
	}
}

func (s *Screen) eraseLine(row int) {
	if row >= 0 && row < s.rows {
		for c := range s.cells[row] {
			s.cells[row][c] = DefaultCell()
		}
	}
}

func (s *Screen) insertLines(n int) {
	for i := 0; i < n; i++ {
		copy(s.cells[s.cursor.Y+1:s.scrollBottom+1], s.cells[s.cursor.Y:s.scrollBottom])
		s.cells[s.cursor.Y] = make([]Cell, s.cols)
		for j := range s.cells[s.cursor.Y] {
			s.cells[s.cursor.Y][j] = DefaultCell()
		}
	}
}

func (s *Screen) deleteLines(n int) {
	for i := 0; i < n; i++ {
		copy(s.cells[s.cursor.Y:s.scrollBottom], s.cells[s.cursor.Y+1:s.scrollBottom+1])
		s.cells[s.scrollBottom] = make([]Cell, s.cols)
		for j := range s.cells[s.scrollBottom] {
			s.cells[s.scrollBottom][j] = DefaultCell()
		}
	}
}

func (s *Screen) deleteChars(n int) {
	row := s.cells[s.cursor.Y]
	end := s.cols - n
	copy(row[s.cursor.X:end], row[s.cursor.X+n:])
	for i := end; i < s.cols; i++ {
		row[i] = DefaultCell()
	}
}

func (s *Screen) insertChars(n int) {
	row := s.cells[s.cursor.Y]
	copy(row[s.cursor.X+n:], row[s.cursor.X:s.cols-n])
	for i := s.cursor.X; i < s.cursor.X+n && i < s.cols; i++ {
		row[i] = DefaultCell()
	}
}

// --- Utilities ---

func parseParams(s string) []int {
	if s == "" {
		return nil
	}
	var nums []int
	cur := 0
	hasCur := false
	for _, c := range s {
		if c >= '0' && c <= '9' {
			cur = cur*10 + int(c-'0')
			hasCur = true
		} else if c == ';' {
			nums = append(nums, cur)
			cur = 0
			hasCur = false
		}
	}
	if hasCur || len(s) > 0 {
		nums = append(nums, cur)
	}
	return nums
}

func cellAttrsEqual(a, b Cell) bool {
	return a.FG == b.FG && a.BG == b.BG &&
		a.Bold == b.Bold && a.Underline == b.Underline &&
		a.Reverse == b.Reverse && a.Dim == b.Dim &&
		a.Italic == b.Italic && a.Blink == b.Blink
}

func cellToANSI(cell Cell) string {
	var sb strings.Builder
	sb.WriteString("\033[0")
	if cell.Bold {
		sb.WriteString(";1")
	}
	if cell.Dim {
		sb.WriteString(";2")
	}
	if cell.Italic {
		sb.WriteString(";3")
	}
	if cell.Underline {
		sb.WriteString(";4")
	}
	if cell.Reverse {
		sb.WriteString(";7")
	}
	switch cell.FG.Mode {
	case ColorANSI:
		if cell.FG.Index < 8 {
			sb.WriteString(";3" + string(rune('0'+cell.FG.Index)))
		} else {
			sb.WriteString(";9" + string(rune('0'+cell.FG.Index-8)))
		}
	case Color256:
		sb.WriteString(";38;5;" + itoa(int(cell.FG.Index)))
	case ColorRGB:
		sb.WriteString(";38;2;" + itoa(int(cell.FG.R)) + ";" + itoa(int(cell.FG.G)) + ";" + itoa(int(cell.FG.B)))
	}
	switch cell.BG.Mode {
	case ColorANSI:
		if cell.BG.Index < 8 {
			sb.WriteString(";4" + string(rune('0'+cell.BG.Index)))
		} else {
			sb.WriteString(";10" + string(rune('0'+cell.BG.Index-8)))
		}
	case Color256:
		sb.WriteString(";48;5;" + itoa(int(cell.BG.Index)))
	case ColorRGB:
		sb.WriteString(";48;2;" + itoa(int(cell.BG.R)) + ";" + itoa(int(cell.BG.G)) + ";" + itoa(int(cell.BG.B)))
	}
	sb.WriteString("m")
	return sb.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}

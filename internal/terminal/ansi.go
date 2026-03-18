package terminal

import "strings"

// OSC133Event represents a parsed OSC 133 shell-integration event.
type OSC133Event int

const (
	OSC133None         OSC133Event = iota
	OSC133PromptStart              // A — prompt about to be drawn
	OSC133PromptEnd                // B — prompt finished
	OSC133CommandStart             // C — command about to run
	OSC133CommandEnd               // D — command finished
)

// ParseOSC133 inspects raw PTY output for OSC 133 sequences and returns
// the last detected event. It does NOT strip the sequences from the data.
func ParseOSC133(data []byte) OSC133Event {
	s := string(data)
	last := OSC133None
	idx := 0
	for {
		start := strings.Index(s[idx:], "\x1b]133;")
		if start < 0 {
			break
		}
		start += idx
		end := strings.IndexAny(s[start:], "\x07\x1b")
		if end < 0 {
			break
		}
		end += start
		payload := s[start+6 : end] // skip "\x1b]133;"
		switch {
		case strings.HasPrefix(payload, "A"):
			last = OSC133PromptStart
		case strings.HasPrefix(payload, "B"):
			last = OSC133PromptEnd
		case strings.HasPrefix(payload, "C"):
			last = OSC133CommandStart
		case strings.HasPrefix(payload, "D"):
			last = OSC133CommandEnd
		}
		idx = end + 1
	}
	return last
}

// AIProcessNames is the set of process names that indicate an AI tool is running.
var AIProcessNames = map[string]bool{
	"claude":  true,
	"aider":   true,
	"llm":     true,
	"ollama":  true,
	"sgpt":    true,
	"gpt4all": true,
	"chatgpt": true,
}

// IsAIProcess returns true if the given process name is an AI tool.
func IsAIProcess(name string) bool {
	base := name
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		base = name[idx+1:]
	}
	return AIProcessNames[strings.ToLower(base)]
}

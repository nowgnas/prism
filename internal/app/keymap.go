package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all global key bindings for prism.
// All actions use Ctrl+B as a prefix.
type KeyMap struct {
	NewTab      key.Binding
	CloseTab    key.Binding
	RenameTab   key.Binding
	NextTab     key.Binding
	SplitH      key.Binding // horizontal split (add pane below)
	SplitV      key.Binding // vertical split (add pane right)
	ClosePane   key.Binding
	ZoomPane    key.Binding
	FocusLeft   key.Binding
	FocusDown   key.Binding
	FocusUp     key.Binding
	FocusRight  key.Binding
	Help        key.Binding
	Quit        key.Binding
	Tab1        key.Binding
	Tab2        key.Binding
	Tab3        key.Binding
	Tab4        key.Binding
	Tab5        key.Binding
	Tab6        key.Binding
	Tab7        key.Binding
	Tab8        key.Binding
	Tab9        key.Binding
}

// DefaultKeyMap returns the default prism key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		NewTab:     key.NewBinding(key.WithKeys("n"), key.WithHelp("^B n", "new tab")),
		CloseTab:   key.NewBinding(key.WithKeys("w"), key.WithHelp("^B w", "close tab")),
		RenameTab:  key.NewBinding(key.WithKeys(","), key.WithHelp("^B ,", "rename tab")),
		NextTab:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("^B Tab", "next tab")),
		SplitH:     key.NewBinding(key.WithKeys("\""), key.WithHelp("^B \"", "split horizontal")),
		SplitV:     key.NewBinding(key.WithKeys("%"), key.WithHelp("^B %", "split vertical")),
		ClosePane:  key.NewBinding(key.WithKeys("x"), key.WithHelp("^B x", "close pane")),
		ZoomPane:   key.NewBinding(key.WithKeys("z"), key.WithHelp("^B z", "zoom pane")),
		FocusLeft:  key.NewBinding(key.WithKeys("h"), key.WithHelp("^B h", "focus left")),
		FocusDown:  key.NewBinding(key.WithKeys("j"), key.WithHelp("^B j", "focus down")),
		FocusUp:    key.NewBinding(key.WithKeys("k"), key.WithHelp("^B k", "focus up")),
		FocusRight: key.NewBinding(key.WithKeys("l"), key.WithHelp("^B l", "focus right")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("^B ?", "help")),
		Quit:       key.NewBinding(key.WithKeys("q"), key.WithHelp("^B q", "quit")),
		Tab1:       key.NewBinding(key.WithKeys("1")),
		Tab2:       key.NewBinding(key.WithKeys("2")),
		Tab3:       key.NewBinding(key.WithKeys("3")),
		Tab4:       key.NewBinding(key.WithKeys("4")),
		Tab5:       key.NewBinding(key.WithKeys("5")),
		Tab6:       key.NewBinding(key.WithKeys("6")),
		Tab7:       key.NewBinding(key.WithKeys("7")),
		Tab8:       key.NewBinding(key.WithKeys("8")),
		Tab9:       key.NewBinding(key.WithKeys("9")),
	}
}

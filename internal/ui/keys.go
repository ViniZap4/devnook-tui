package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Back     key.Binding
	Quit     key.Binding
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Top      key.Binding
	Bottom   key.Binding
	HalfUp   key.Binding
	HalfDown key.Binding
	Search   key.Binding
	Help     key.Binding
}

var keys = keyMap{
	Tab:      key.NewBinding(key.WithKeys("tab", "L"), key.WithHelp("tab/L", "next tab")),
	ShiftTab: key.NewBinding(key.WithKeys("shift+tab", "H"), key.WithHelp("S-tab/H", "prev tab")),
	Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	Back:     key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc/q", "back")),
	Quit:     key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("up/k", "up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("down/j", "down")),
	Left:     key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("left/h", "left")),
	Right:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("right/l", "right")),
	Top:      key.NewBinding(key.WithKeys("g"), key.WithHelp("gg", "top")),
	Bottom:   key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "bottom")),
	HalfUp:   key.NewBinding(key.WithKeys("ctrl+u"), key.WithHelp("C-u", "half page up")),
	HalfDown: key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("C-d", "half page down")),
	Search:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}

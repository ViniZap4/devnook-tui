package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Tab   key.Binding
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Up    key.Binding
	Down  key.Binding
}

var keys = keyMap{
	Tab:   key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
	Back:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	Quit:  key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("up/k", "up")),
	Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("down/j", "down")),
}

package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func applyVimNavigation(cursor, scroll *int, gPressed *bool, msg tea.KeyMsg, maxIdx int) {
	if key.Matches(msg, keys.Down) {
		*gPressed = false
		if *cursor < maxIdx {
			*cursor++
		}
	}
	if key.Matches(msg, keys.Up) {
		*gPressed = false
		if *cursor > 0 {
			*cursor--
		}
	}
	if msg.String() == "g" {
		if *gPressed {
			*cursor = 0
			*scroll = 0
			*gPressed = false
		} else {
			*gPressed = true
		}
		return
	}
	if key.Matches(msg, keys.Bottom) {
		*cursor = maxIdx
		*gPressed = false
	}
	if key.Matches(msg, keys.HalfDown) {
		*gPressed = false
		*cursor += 10
		if *cursor > maxIdx {
			*cursor = maxIdx
		}
	}
	if key.Matches(msg, keys.HalfUp) {
		*gPressed = false
		*cursor -= 10
		if *cursor < 0 {
			*cursor = 0
		}
	}
	if msg.String() != "g" {
		*gPressed = false
	}
}

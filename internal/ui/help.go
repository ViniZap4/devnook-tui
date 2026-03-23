package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderHelp(width, height int) string {
	title := titleStyle.Render("Keyboard Shortcuts")

	sections := []struct {
		name  string
		binds []string
	}{
		{"Navigation", []string{
			"j / down     Move down",
			"k / up       Move up",
			"h / left     Go back",
			"l / right    Open / expand",
			"enter        Open selected item",
			"esc          Go back",
			"tab          Next tab",
			"shift+tab    Previous tab",
			"1-4          Jump to tab",
		}},
		{"Scrolling", []string{
			"gg           Go to top",
			"G            Go to bottom",
			"ctrl+d       Half page down",
			"ctrl+u       Half page up",
		}},
		{"Global", []string{
			"?            Toggle help",
			"ctrl+c       Quit",
		}},
		{"Views", []string{
			"Dashboard    Repos, orgs, shortcuts overview",
			"Repo Detail  Files, commits, branches, issues",
			"Issue List   Browse open/closed issues",
			"Issue Detail Read issue body and comments",
			"File Viewer  View file contents with line numbers",
		}},
	}

	var lines []string
	for _, s := range sections {
		lines = append(lines, "")
		lines = append(lines, subtitleStyle.Render(s.name))
		for _, b := range s.binds {
			parts := strings.SplitN(b, " ", 2)
			if len(parts) == 2 {
				key := parts[0]
				desc := strings.TrimSpace(parts[1])
				lines = append(lines, "  "+accentStyle.Render(padRight(key, 14))+mutedStyle.Render(desc))
			}
		}
	}

	lines = append(lines, "")
	lines = append(lines, mutedStyle.Render("Press ? or esc to close"))

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{title}, lines...)...)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func padRight(s string, length int) string {
	for len(s) < length {
		s += " "
	}
	return s
}

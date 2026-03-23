package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type fileViewerModel struct {
	owner   string
	repo    string
	ref     string
	path    string
	name    string
	content string
	binary  bool
	scroll  int
	loading bool
}

func newFileViewerModel(owner, repo, ref, path, name string) fileViewerModel {
	return fileViewerModel{owner: owner, repo: repo, ref: ref, path: path, name: name, loading: true}
}

func (m fileViewerModel) Update(msg tea.Msg) (fileViewerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Down) {
			m.scroll++
		}
		if key.Matches(msg, keys.Up) {
			if m.scroll > 0 {
				m.scroll--
			}
		}
		if key.Matches(msg, keys.HalfDown) {
			m.scroll += 10
		}
		if key.Matches(msg, keys.HalfUp) {
			m.scroll -= 10
			if m.scroll < 0 {
				m.scroll = 0
			}
		}
		if msg.String() == "g" {
			m.scroll = 0
		}
		if key.Matches(msg, keys.Bottom) {
			m.scroll = 99999
		}
	}
	return m, nil
}

func (m fileViewerModel) View(width, height int) string {
	if width == 0 {
		return ""
	}

	cw := min(width-4, 120)
	if cw < 20 {
		cw = 20
	}

	title := titleStyle.Render(m.name)
	pathLine := mutedStyle.Render(m.path)

	if m.loading {
		content := lipgloss.JoinVertical(lipgloss.Left, title, pathLine, "", mutedStyle.Render("Loading..."))
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}

	if m.binary {
		content := lipgloss.JoinVertical(lipgloss.Left, title, pathLine, "",
			mutedStyle.Render("Binary file — cannot display"))
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}

	lines := strings.Split(m.content, "\n")
	totalLines := len(lines)

	// Apply scroll
	start := m.scroll
	if start >= totalLines {
		start = totalLines - 1
	}
	if start < 0 {
		start = 0
	}

	maxVisible := height - 6
	if maxVisible < 3 {
		maxVisible = 3
	}

	end := start + maxVisible
	if end > totalLines {
		end = totalLines
	}

	var rendered []string
	lineNumWidth := len(fmt.Sprintf("%d", totalLines))
	for i := start; i < end; i++ {
		num := fmt.Sprintf("%*d", lineNumWidth, i+1)
		line := lines[i]
		if len(line) > cw-lineNumWidth-3 {
			line = line[:cw-lineNumWidth-6] + "..."
		}
		rendered = append(rendered, fmt.Sprintf("%s %s %s",
			mutedStyle.Render(num),
			mutedStyle.Render("|"),
			textStyle.Render(line)))
	}

	statusLeft := statusBarActiveStyle.Render(fmt.Sprintf(" %s ", m.name))
	pct := 0
	if totalLines > 1 {
		pct = int(float64(start) / float64(totalLines-1) * 100)
	}
	statusRight := statusBarStyle.Render(fmt.Sprintf(" %d lines  %d%%  esc back  j/k scroll  ? help", totalLines, pct))
	statusBar := lipgloss.JoinHorizontal(lipgloss.Top, statusLeft, statusRight)

	content := strings.Join(rendered, "\n")
	page := lipgloss.JoinVertical(lipgloss.Left, title, pathLine, "", content, "", statusBar)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, page)
}

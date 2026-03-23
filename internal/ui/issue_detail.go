package ui

import (
	"fmt"
	"strings"

	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type issueDetailModel struct {
	owner    string
	repo     string
	number   int
	issue    *api.Issue
	comments []api.IssueComment
	scroll   int
	loading  bool
}

func newIssueDetailModel(owner, repo string, number int) issueDetailModel {
	return issueDetailModel{owner: owner, repo: repo, number: number, loading: true}
}

func (m issueDetailModel) Update(msg tea.Msg) (issueDetailModel, tea.Cmd) {
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
	}
	return m, nil
}

func (m issueDetailModel) View(width, height int) string {
	if width == 0 {
		return ""
	}

	cw := min(width-4, 90)
	if cw < 20 {
		cw = 20
	}

	if m.loading {
		content := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(fmt.Sprintf("Issue #%d", m.number)), "",
			mutedStyle.Render("Loading..."))
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}

	if m.issue == nil {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
			errorStyle.Render("Issue not found"))
	}

	var lines []string

	// Title
	stateIcon := successStyle.Render("[Open]")
	if m.issue.State == "closed" {
		stateIcon = errorStyle.Render("[Closed]")
	}
	lines = append(lines, titleStyle.Render(fmt.Sprintf("#%d %s", m.issue.Number, m.issue.Title)))
	lines = append(lines, fmt.Sprintf("%s  %s %s",
		stateIcon, mutedStyle.Render("by"), textStyle.Render(m.issue.Author)))
	lines = append(lines, "")

	// Body
	if m.issue.Body != "" {
		body := wrapText(m.issue.Body, cw)
		lines = append(lines, textStyle.Render(body))
	} else {
		lines = append(lines, mutedStyle.Render("No description provided."))
	}

	// Comments
	if len(m.comments) > 0 {
		lines = append(lines, "")
		lines = append(lines, subtitleStyle.Render(fmt.Sprintf("Comments (%d)", len(m.comments))))
		lines = append(lines, mutedStyle.Render(strings.Repeat("─", min(cw, 60))))

		for _, c := range m.comments {
			lines = append(lines, "")
			lines = append(lines, fmt.Sprintf("%s %s",
				accentStyle.Render(c.Author), mutedStyle.Render("commented:")))
			body := wrapText(c.Body, cw)
			lines = append(lines, textStyle.Render(body))
		}
	}

	// Apply scroll
	if m.scroll > 0 && m.scroll < len(lines) {
		lines = lines[m.scroll:]
	}

	maxLines := height - 4
	if maxLines > 0 && len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	statusLeft := statusBarActiveStyle.Render(fmt.Sprintf(" #%d ", m.number))
	statusRight := statusBarStyle.Render(" esc back  j/k scroll  ctrl+d/u page  ? help")
	statusBar := lipgloss.JoinHorizontal(lipgloss.Top, statusLeft, statusRight)

	content := strings.Join(lines, "\n")
	page := lipgloss.JoinVertical(lipgloss.Left, content, "", statusBar)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, page)
}

func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}
	var result strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if len(line) <= maxWidth {
			result.WriteString(line)
			result.WriteString("\n")
			continue
		}
		words := strings.Fields(line)
		current := ""
		for _, word := range words {
			if len(current)+len(word)+1 > maxWidth {
				result.WriteString(current)
				result.WriteString("\n")
				current = word
			} else if current == "" {
				current = word
			} else {
				current += " " + word
			}
		}
		if current != "" {
			result.WriteString(current)
			result.WriteString("\n")
		}
	}
	return strings.TrimRight(result.String(), "\n")
}

package ui

import (
	"fmt"
	"strings"

	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type issueListModel struct {
	owner    string
	repo     string
	issues   []api.Issue
	cursor   int
	scroll   int
	loading  bool
	gPressed bool
}

func newIssueListModel(owner, repo string) issueListModel {
	return issueListModel{owner: owner, repo: repo, loading: true}
}

func (m issueListModel) Update(msg tea.Msg) (issueListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		maxIdx := len(m.issues) - 1
		if maxIdx < 0 {
			maxIdx = 0
		}

		if key.Matches(msg, keys.Enter) {
			if m.cursor < len(m.issues) {
				issue := m.issues[m.cursor]
				return m, func() tea.Msg {
					return openIssueMsg{owner: m.owner, repo: m.repo, number: issue.Number}
				}
			}
			return m, nil
		}

		applyVimNavigation(&m.cursor, &m.scroll, &m.gPressed, msg, maxIdx)
	}
	return m, nil
}

func (m issueListModel) View(width, height int) string {
	if width == 0 {
		return ""
	}

	cw := min(width-4, 100)
	if cw < 20 {
		cw = 20
	}

	title := titleStyle.Render(fmt.Sprintf("Issues — %s/%s", m.owner, m.repo))

	if m.loading {
		content := lipgloss.JoinVertical(lipgloss.Left, title, "", mutedStyle.Render("Loading issues..."))
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}

	if len(m.issues) == 0 {
		content := lipgloss.JoinVertical(lipgloss.Left, title, "", mutedStyle.Render("No open issues."))
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}

	visible := height - 8
	if visible < 3 {
		visible = 3
	}
	if visible > len(m.issues) {
		visible = len(m.issues)
	}

	scroll := m.scroll
	if m.cursor < scroll {
		scroll = m.cursor
	}
	if m.cursor >= scroll+visible {
		scroll = m.cursor - visible + 1
	}

	var lines []string
	end := scroll + visible
	if end > len(m.issues) {
		end = len(m.issues)
	}

	for i := scroll; i < end; i++ {
		issue := m.issues[i]
		style := normalItemStyle
		prefix := "  "
		if i == m.cursor {
			style = selectedItemStyle
			prefix = "> "
		}

		stateIcon := successStyle.Render("O")
		if issue.State == "closed" {
			stateIcon = errorStyle.Render("C")
		}

		num := mutedStyle.Render(fmt.Sprintf("#%d", issue.Number))
		titleText := issue.Title
		if len(titleText) > 60 {
			titleText = titleText[:57] + "..."
		}

		line := fmt.Sprintf("%s%s %s %s  %s",
			style.Render(prefix), stateIcon, num, style.Render(titleText),
			mutedStyle.Render(issue.Author))
		lines = append(lines, line)
	}

	statusLeft := statusBarActiveStyle.Render(fmt.Sprintf(" %d issues ", len(m.issues)))
	statusRight := statusBarStyle.Render(" esc back  j/k nav  enter open  ? help")
	statusBar := lipgloss.JoinHorizontal(lipgloss.Top, statusLeft, statusRight)

	content := strings.Join(lines, "\n")
	page := lipgloss.JoinVertical(lipgloss.Left, title, "", content, "", statusBar)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, page)
}

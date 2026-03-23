package ui

import (
	"fmt"
	"strings"

	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type repoTab int

const (
	repoTabFiles repoTab = iota
	repoTabCommits
	repoTabBranches
	repoTabIssues
)

type repoDetailModel struct {
	owner    string
	name     string
	repo     *api.Repo
	branches []api.Branch
	commits  []api.Commit
	tree     []api.TreeEntry
	tab      repoTab
	cursor   int
	scroll   int
	loading  bool
	gPressed bool
}

func newRepoDetailModel(owner, name string) repoDetailModel {
	return repoDetailModel{owner: owner, name: name, loading: true}
}

func (m repoDetailModel) listLen() int {
	switch m.tab {
	case repoTabFiles:
		return len(m.tree)
	case repoTabCommits:
		return len(m.commits)
	case repoTabBranches:
		return len(m.branches)
	default:
		return 0
	}
}

func (m repoDetailModel) Update(msg tea.Msg) (repoDetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Tab nav
		if key.Matches(msg, keys.Tab) {
			m.tab = (m.tab + 1) % 4
			m.cursor = 0
			m.scroll = 0
			return m, nil
		}
		if key.Matches(msg, keys.ShiftTab) {
			m.tab = (m.tab + 3) % 4
			m.cursor = 0
			m.scroll = 0
			return m, nil
		}

		switch msg.String() {
		case "1":
			m.tab = repoTabFiles
			m.cursor, m.scroll = 0, 0
		case "2":
			m.tab = repoTabCommits
			m.cursor, m.scroll = 0, 0
		case "3":
			m.tab = repoTabBranches
			m.cursor, m.scroll = 0, 0
		case "4":
			m.tab = repoTabIssues
			m.cursor, m.scroll = 0, 0
			return m, func() tea.Msg {
				return openIssueListMsg{owner: m.owner, repo: m.name}
			}
		}

		maxIdx := m.listLen() - 1
		if maxIdx < 0 {
			maxIdx = 0
		}

		// Enter to open item
		if key.Matches(msg, keys.Enter) {
			switch m.tab {
			case repoTabFiles:
				if m.cursor < len(m.tree) {
					entry := m.tree[m.cursor]
					if entry.Type == "blob" {
						ref := "HEAD"
						if m.repo != nil {
							ref = m.repo.DefaultBranch
						}
						return m, func() tea.Msg {
							return openFileMsg{
								owner: m.owner, repo: m.name,
								ref: ref, path: entry.Path, name: entry.Name,
							}
						}
					}
				}
			}
			return m, nil
		}

		applyVimNavigation(&m.cursor, &m.scroll, &m.gPressed, msg, maxIdx)
	}
	return m, nil
}

func (m repoDetailModel) View(width, height int) string {
	if width == 0 {
		return ""
	}

	cw := min(width-4, 100)
	if cw < 20 {
		cw = 20
	}

	// Header
	header := m.renderHeader(cw)

	// Tabs
	tabs := m.renderTabs()

	// Content
	var content string
	if m.loading {
		content = mutedStyle.Render("Loading repository...")
	} else {
		switch m.tab {
		case repoTabFiles:
			content = m.renderFiles(cw, height-12)
		case repoTabCommits:
			content = m.renderCommits(cw, height-12)
		case repoTabBranches:
			content = m.renderBranches(cw, height-12)
		}
	}

	statusBar := m.renderStatusBar(cw)

	page := lipgloss.JoinVertical(lipgloss.Left, header, "", tabs, "", content, "", statusBar)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, page)
}

func (m repoDetailModel) renderHeader(width int) string {
	if m.repo == nil {
		return titleStyle.Render(fmt.Sprintf("%s/%s", m.owner, m.name))
	}

	title := titleStyle.Render(fmt.Sprintf("%s/%s", m.repo.Owner, m.repo.Name))

	var badges []string
	if m.repo.IsPrivate {
		badges = append(badges, privateBadge.Render())
	} else {
		badges = append(badges, publicBadge.Render())
	}
	if m.repo.IsFork {
		badges = append(badges, mutedStyle.Render("fork"))
	}
	badges = append(badges, mutedStyle.Render(fmt.Sprintf("* %d", m.repo.StarsCount)))
	badges = append(badges, mutedStyle.Render(fmt.Sprintf("Y %d", m.repo.ForksCount)))

	info := strings.Join(badges, "  ")
	desc := ""
	if m.repo.Description != "" {
		desc = "\n" + descriptionStyle.Render(m.repo.Description)
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, info+desc)
}

func (m repoDetailModel) renderTabs() string {
	tabDefs := []struct {
		label string
		t     repoTab
		key   string
	}{
		{"Files", repoTabFiles, "1"},
		{"Commits", repoTabCommits, "2"},
		{"Branches", repoTabBranches, "3"},
		{"Issues", repoTabIssues, "4"},
	}

	var parts []string
	for _, t := range tabDefs {
		label := fmt.Sprintf(" %s %s ", t.key, t.label)
		if m.tab == t.t {
			parts = append(parts, activeTabStyle.Render(label))
		} else {
			parts = append(parts, inactiveTabStyle.Render(label))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (m repoDetailModel) renderFiles(width, maxHeight int) string {
	if len(m.tree) == 0 {
		return mutedStyle.Render("Empty repository. Push some code to get started!")
	}

	visible := maxHeight - 2
	if visible < 3 {
		visible = 3
	}
	if visible > len(m.tree) {
		visible = len(m.tree)
	}

	scroll := m.scroll
	if m.cursor < scroll {
		scroll = m.cursor
	}
	if m.cursor >= scroll+visible {
		scroll = m.cursor - visible + 1
	}

	var lines []string
	header := fmt.Sprintf(" %s  %-40s  %s",
		mutedStyle.Render(" "),
		mutedStyle.Render("Name"),
		mutedStyle.Render("Size"),
	)
	lines = append(lines, header)
	lines = append(lines, mutedStyle.Render(strings.Repeat("─", min(width, 60))))

	end := scroll + visible
	if end > len(m.tree) {
		end = len(m.tree)
	}

	for i := scroll; i < end; i++ {
		entry := m.tree[i]
		icon := "  "
		if entry.Type == "tree" {
			icon = accentStyle.Render("D ")
		} else {
			icon = mutedStyle.Render("F ")
		}

		name := entry.Name
		if len(name) > 40 {
			name = name[:37] + "..."
		}

		size := ""
		if entry.Type == "blob" && entry.Size > 0 {
			size = formatSize(entry.Size)
		}

		style := normalItemStyle
		prefix := "  "
		if i == m.cursor {
			style = selectedItemStyle
			prefix = "> "
		}

		line := fmt.Sprintf("%s%s%-40s  %s",
			style.Render(prefix), icon, style.Render(name), mutedStyle.Render(size))
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m repoDetailModel) renderCommits(width, maxHeight int) string {
	if len(m.commits) == 0 {
		return mutedStyle.Render("No commits yet.")
	}

	visible := maxHeight - 2
	if visible < 3 {
		visible = 3
	}
	if visible > len(m.commits) {
		visible = len(m.commits)
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
	if end > len(m.commits) {
		end = len(m.commits)
	}

	for i := scroll; i < end; i++ {
		c := m.commits[i]
		style := normalItemStyle
		prefix := "  "
		if i == m.cursor {
			style = selectedItemStyle
			prefix = "> "
		}

		hash := accentStyle.Render(c.ShortHash)
		msg := c.Message
		if len(msg) > 50 {
			msg = msg[:47] + "..."
		}

		line := fmt.Sprintf("%s%s %s %s",
			style.Render(prefix), hash, style.Render(msg), mutedStyle.Render(c.Author))
		lines = append(lines, line)
	}

	if len(m.commits) > visible {
		pct := float64(m.cursor) / float64(len(m.commits)-1) * 100
		lines = append(lines, "")
		lines = append(lines, mutedStyle.Render(fmt.Sprintf(" %d/%d (%.0f%%)", m.cursor+1, len(m.commits), pct)))
	}

	return strings.Join(lines, "\n")
}

func (m repoDetailModel) renderBranches(width, maxHeight int) string {
	if len(m.branches) == 0 {
		return mutedStyle.Render("No branches yet.")
	}

	var lines []string
	for i, b := range m.branches {
		style := normalItemStyle
		prefix := "  "
		if i == m.cursor {
			style = selectedItemStyle
			prefix = "> "
		}

		name := style.Render(b.Name)
		extra := ""
		if b.IsDefault {
			extra = successStyle.Render(" default")
		}
		if b.IsHead {
			extra += accentStyle.Render(" HEAD")
		}

		lines = append(lines, fmt.Sprintf("%s%s%s", style.Render(prefix), name, extra))
	}
	return strings.Join(lines, "\n")
}

func (m repoDetailModel) renderStatusBar(width int) string {
	left := statusBarActiveStyle.Render(fmt.Sprintf(" %s/%s ", m.owner, m.name))
	help := " esc back  1-4 tabs  j/k nav  enter open  ? help"
	right := statusBarStyle.Render(help)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fK", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1fM", float64(bytes)/(1024*1024))
}

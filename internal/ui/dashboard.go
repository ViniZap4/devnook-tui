package ui

import (
	"fmt"
	"strings"

	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	tabOverview tab = iota
	tabRepos
	tabShortcuts
)

type dashboardModel struct {
	user      *api.User
	repos     []api.Repo
	shortcuts []api.Shortcut
	tab       tab
	cursor    int
	scroll    int
	loading   bool
	gPressed  bool // for gg motion
}

func newDashboardModel() dashboardModel {
	return dashboardModel{loading: true}
}

func (m dashboardModel) listLen() int {
	switch m.tab {
	case tabRepos:
		return len(m.repos)
	case tabShortcuts:
		return len(m.shortcuts)
	default:
		return 0
	}
}

func (m dashboardModel) Update(msg tea.Msg) (dashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Tab navigation
		if key.Matches(msg, keys.Tab) {
			m.tab = (m.tab + 1) % 3
			m.cursor = 0
			m.scroll = 0
			m.gPressed = false
			return m, nil
		}
		if key.Matches(msg, keys.ShiftTab) {
			m.tab = (m.tab + 2) % 3
			m.cursor = 0
			m.scroll = 0
			m.gPressed = false
			return m, nil
		}

		// Number keys for tabs
		switch msg.String() {
		case "1":
			m.tab = tabOverview
			m.cursor = 0
			m.scroll = 0
		case "2":
			m.tab = tabRepos
			m.cursor = 0
			m.scroll = 0
		case "3":
			m.tab = tabShortcuts
			m.cursor = 0
			m.scroll = 0
		}

		maxIdx := m.listLen() - 1
		if maxIdx < 0 {
			maxIdx = 0
		}

		// Vim motions for list navigation
		if key.Matches(msg, keys.Down) {
			m.gPressed = false
			if m.cursor < maxIdx {
				m.cursor++
			}
		}
		if key.Matches(msg, keys.Up) {
			m.gPressed = false
			if m.cursor > 0 {
				m.cursor--
			}
		}

		// gg - go to top
		if msg.String() == "g" {
			if m.gPressed {
				m.cursor = 0
				m.scroll = 0
				m.gPressed = false
			} else {
				m.gPressed = true
			}
			return m, nil
		}

		// G - go to bottom
		if key.Matches(msg, keys.Bottom) {
			m.cursor = maxIdx
			m.gPressed = false
		}

		// Ctrl+D - half page down
		if key.Matches(msg, keys.HalfDown) {
			m.gPressed = false
			m.cursor += 10
			if m.cursor > maxIdx {
				m.cursor = maxIdx
			}
		}

		// Ctrl+U - half page up
		if key.Matches(msg, keys.HalfUp) {
			m.gPressed = false
			m.cursor -= 10
			if m.cursor < 0 {
				m.cursor = 0
			}
		}

		// Reset g if any other key was pressed
		if msg.String() != "g" {
			m.gPressed = false
		}
	}
	return m, nil
}

func (m dashboardModel) View(width, height int) string {
	if width == 0 {
		return ""
	}

	contentWidth := width - 4
	if contentWidth > 100 {
		contentWidth = 100
	}
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Header
	header := m.renderHeader(contentWidth)

	// Tabs
	tabs := m.renderTabs(contentWidth)

	// Content
	var content string
	switch m.tab {
	case tabOverview:
		content = m.overviewView(contentWidth)
	case tabRepos:
		content = m.reposView(contentWidth, height-10)
	case tabShortcuts:
		content = m.shortcutsView(contentWidth, height-10)
	}

	// Status bar
	statusBar := m.renderStatusBar(contentWidth)

	page := lipgloss.JoinVertical(lipgloss.Left,
		header, "", tabs, "", content, "", statusBar,
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, page)
}

func (m dashboardModel) renderHeader(width int) string {
	title := titleStyle.Render("Dev Nook")

	if m.loading {
		return lipgloss.JoinVertical(lipgloss.Left,
			title,
			mutedStyle.Render("Loading..."),
		)
	}

	userInfo := ""
	if m.user != nil {
		name := m.user.FullName
		if name == "" {
			name = m.user.Username
		}
		userInfo = textStyle.Render(fmt.Sprintf("Welcome, %s", name)) +
			" " + mutedStyle.Render(fmt.Sprintf("@%s", m.user.Username))
	}

	stats := mutedStyle.Render(fmt.Sprintf(
		"%s %d repos  %s %d shortcuts",
		accentStyle.Render("*"),
		len(m.repos),
		accentStyle.Render("*"),
		len(m.shortcuts),
	))

	return lipgloss.JoinVertical(lipgloss.Left, title, userInfo, stats)
}

func (m dashboardModel) renderTabs(width int) string {
	tabs := []struct {
		label string
		t     tab
		key   string
	}{
		{"Overview", tabOverview, "1"},
		{"Repos", tabRepos, "2"},
		{"Shortcuts", tabShortcuts, "3"},
	}

	var parts []string
	for _, t := range tabs {
		label := fmt.Sprintf(" %s %s ", t.key, t.label)
		if m.tab == t.t {
			parts = append(parts, activeTabStyle.Render(label))
		} else {
			parts = append(parts, inactiveTabStyle.Render(label))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (m dashboardModel) renderStatusBar(width int) string {
	left := statusBarActiveStyle.Render(" DEV NOOK ")
	var help string
	switch m.tab {
	case tabOverview:
		help = " 1/2/3 tabs  tab/S-tab next/prev  ctrl+c quit"
	default:
		help = " j/k up/down  gg top  G bottom  C-d/C-u half page  tab next  ctrl+c quit"
	}
	right := statusBarStyle.Render(help)
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m dashboardModel) overviewView(width int) string {
	if m.loading {
		return mutedStyle.Render("Loading dashboard data...")
	}

	var sections []string

	// Recent repos
	repoHeader := subtitleStyle.Render("Recent Repositories")
	sections = append(sections, repoHeader)

	if len(m.repos) == 0 {
		sections = append(sections, mutedStyle.Render("  No repositories yet."))
	} else {
		limit := 5
		if len(m.repos) < limit {
			limit = len(m.repos)
		}
		for _, r := range m.repos[:limit] {
			badge := publicBadge
			if r.IsPrivate {
				badge = privateBadge
			}
			line := fmt.Sprintf("  %s/%s %s",
				mutedStyle.Render(r.Owner),
				textStyle.Render(r.Name),
				badge.Render(),
			)
			sections = append(sections, line)
			if r.Description != "" {
				sections = append(sections, fmt.Sprintf("    %s", descriptionStyle.Render(r.Description)))
			}
		}
		if len(m.repos) > 5 {
			sections = append(sections, mutedStyle.Render(fmt.Sprintf("  ... and %d more", len(m.repos)-5)))
		}
	}

	sections = append(sections, "")

	// Recent shortcuts
	scHeader := subtitleStyle.Render("Shortcuts")
	sections = append(sections, scHeader)

	if len(m.shortcuts) == 0 {
		sections = append(sections, mutedStyle.Render("  No shortcuts yet."))
	} else {
		limit := 5
		if len(m.shortcuts) < limit {
			limit = len(m.shortcuts)
		}
		for _, sc := range m.shortcuts[:limit] {
			line := fmt.Sprintf("  %s %s %s",
				accentStyle.Render("->"),
				textStyle.Render(sc.Title),
				descriptionStyle.Render(sc.URL),
			)
			sections = append(sections, line)
		}
	}

	return strings.Join(sections, "\n")
}

func (m dashboardModel) reposView(width, maxHeight int) string {
	if len(m.repos) == 0 {
		return mutedStyle.Render("No repositories yet. Create one from the web client!")
	}

	visible := maxHeight - 2
	if visible < 3 {
		visible = 3
	}
	if visible > len(m.repos) {
		visible = len(m.repos)
	}

	// Adjust scroll window
	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+visible {
		m.scroll = m.cursor - visible + 1
	}

	var lines []string
	header := fmt.Sprintf(" %s  %-30s  %-8s  %s",
		mutedStyle.Render("#"),
		mutedStyle.Render("Repository"),
		mutedStyle.Render("Vis"),
		mutedStyle.Render("Description"),
	)
	lines = append(lines, header)
	lines = append(lines, mutedStyle.Render(strings.Repeat("─", min(width, 80))))

	end := m.scroll + visible
	if end > len(m.repos) {
		end = len(m.repos)
	}

	for i := m.scroll; i < end; i++ {
		r := m.repos[i]
		vis := "public"
		if r.IsPrivate {
			vis = "private"
		}

		name := fmt.Sprintf("%s/%s", r.Owner, r.Name)
		if len(name) > 30 {
			name = name[:27] + "..."
		}

		desc := r.Description
		if len(desc) > 30 {
			desc = desc[:27] + "..."
		}

		style := normalItemStyle
		prefix := "  "
		if i == m.cursor {
			style = selectedItemStyle
			prefix = "> "
		}

		line := fmt.Sprintf("%s%-30s  %-8s  %s",
			style.Render(prefix),
			style.Render(name),
			descriptionStyle.Render(vis),
			descriptionStyle.Render(desc),
		)
		lines = append(lines, line)
	}

	// Scroll indicator
	if len(m.repos) > visible {
		pct := float64(m.cursor) / float64(len(m.repos)-1) * 100
		lines = append(lines, "")
		lines = append(lines, mutedStyle.Render(fmt.Sprintf(" %d/%d (%.0f%%)", m.cursor+1, len(m.repos), pct)))
	}

	return strings.Join(lines, "\n")
}

func (m dashboardModel) shortcutsView(width, maxHeight int) string {
	if len(m.shortcuts) == 0 {
		return mutedStyle.Render("No shortcuts yet. Add some from the web client!")
	}

	visible := maxHeight - 2
	if visible < 3 {
		visible = 3
	}
	if visible > len(m.shortcuts) {
		visible = len(m.shortcuts)
	}

	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+visible {
		m.scroll = m.cursor - visible + 1
	}

	var lines []string
	header := fmt.Sprintf(" %s  %-25s  %s",
		mutedStyle.Render("#"),
		mutedStyle.Render("Title"),
		mutedStyle.Render("URL"),
	)
	lines = append(lines, header)
	lines = append(lines, mutedStyle.Render(strings.Repeat("─", min(width, 70))))

	end := m.scroll + visible
	if end > len(m.shortcuts) {
		end = len(m.shortcuts)
	}

	for i := m.scroll; i < end; i++ {
		sc := m.shortcuts[i]

		title := sc.Title
		if len(title) > 25 {
			title = title[:22] + "..."
		}
		url := sc.URL
		if len(url) > 40 {
			url = url[:37] + "..."
		}

		style := normalItemStyle
		prefix := "  "
		if i == m.cursor {
			style = selectedItemStyle
			prefix = "> "
		}

		line := fmt.Sprintf("%s%-25s  %s",
			style.Render(prefix),
			style.Render(title),
			descriptionStyle.Render(url),
		)
		lines = append(lines, line)
	}

	if len(m.shortcuts) > visible {
		pct := float64(m.cursor) / float64(len(m.shortcuts)-1) * 100
		lines = append(lines, "")
		lines = append(lines, mutedStyle.Render(fmt.Sprintf(" %d/%d (%.0f%%)", m.cursor+1, len(m.shortcuts), pct)))
	}

	return strings.Join(lines, "\n")
}

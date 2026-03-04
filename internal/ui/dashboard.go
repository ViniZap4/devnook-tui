package ui

import (
	"fmt"

	"github.com/ViniZap4/devnook-tui/internal/api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dashboardModel struct {
	user      *api.User
	repos     []api.Repo
	shortcuts []api.Shortcut
	tab       int
}

func (m dashboardModel) Update(msg tea.Msg) (dashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.tab = 0
		case "2":
			m.tab = 1
		case "3":
			m.tab = 2
		}
	}
	return m, nil
}

func (m dashboardModel) View(width, height int) string {
	header := titleStyle.Render("Dev Nook")

	userInfo := ""
	if m.user != nil {
		userInfo = fmt.Sprintf("Welcome, %s (%s)", m.user.FullName, m.user.Username)
	}

	tabs := []string{"Overview", "Repos", "Shortcuts"}
	var tabRow string
	for i, t := range tabs {
		if i == m.tab {
			tabRow += activeTabStyle.Render(fmt.Sprintf(" [%d] %s ", i+1, t))
		} else {
			tabRow += inactiveTabStyle.Render(fmt.Sprintf(" [%d] %s ", i+1, t))
		}
	}

	var content string
	switch m.tab {
	case 0:
		content = m.overviewView()
	case 1:
		content = m.reposView()
	case 2:
		content = m.shortcutsView()
	}

	help := mutedStyle.Render("1/2/3: switch tab  ctrl+c: quit")

	page := lipgloss.JoinVertical(lipgloss.Left,
		header, userInfo, "", tabRow, "", content, "", help,
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
		boxStyle.Render(page),
	)
}

func (m dashboardModel) overviewView() string {
	return fmt.Sprintf(
		"%s\n%s",
		subtitleStyle.Render(fmt.Sprintf("Repositories: %d", len(m.repos))),
		subtitleStyle.Render(fmt.Sprintf("Shortcuts: %d", len(m.shortcuts))),
	)
}

func (m dashboardModel) reposView() string {
	if len(m.repos) == 0 {
		return mutedStyle.Render("No repositories yet.")
	}
	s := ""
	for _, r := range m.repos {
		vis := "public"
		if r.IsPrivate {
			vis = "private"
		}
		s += fmt.Sprintf("  %s/%s (%s)\n", r.Owner, r.Name, vis)
		if r.Description != "" {
			s += fmt.Sprintf("    %s\n", mutedStyle.Render(r.Description))
		}
	}
	return s
}

func (m dashboardModel) shortcutsView() string {
	if len(m.shortcuts) == 0 {
		return mutedStyle.Render("No shortcuts yet.")
	}
	s := ""
	for _, sc := range m.shortcuts {
		s += fmt.Sprintf("  %s — %s\n", sc.Title, mutedStyle.Render(sc.URL))
	}
	return s
}

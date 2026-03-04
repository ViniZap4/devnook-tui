package ui

import (
	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/ViniZap4/devnook-tui/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	loginView view = iota
	dashboardView
	reposView
	shortcutsView
)

type Model struct {
	client    *api.Client
	config    *config.Config
	view      view
	width     int
	height    int
	login     loginModel
	dashboard dashboardModel
	err       error
}

func NewModel(client *api.Client, cfg *config.Config) Model {
	m := Model{
		client: client,
		config: cfg,
		view:   loginView,
	}
	m.login = newLoginModel()

	if cfg.Token != "" {
		m.view = dashboardView
	}

	return m
}

func (m Model) Init() tea.Cmd {
	if m.config.Token != "" {
		return m.fetchDashboard()
	}
	return checkSetup(m.client)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case loginSuccessMsg:
		m.config.Token = msg.token
		m.config.Save()
		m.view = dashboardView
		return m, m.fetchDashboard()

	case dashboardDataMsg:
		m.dashboard.user = msg.user
		m.dashboard.repos = msg.repos
		m.dashboard.shortcuts = msg.shortcuts
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	switch m.view {
	case loginView:
		m.login, cmd = m.login.Update(msg, m.client)
	case dashboardView:
		m.dashboard, cmd = m.dashboard.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.view {
	case loginView:
		return m.login.View(m.width, m.height)
	case dashboardView:
		return m.dashboard.View(m.width, m.height)
	default:
		return "unknown view"
	}
}

// Messages
type loginSuccessMsg struct{ token string }
type errMsg struct{ err error }

type dashboardDataMsg struct {
	user      *api.User
	repos     []api.Repo
	shortcuts []api.Shortcut
}

func (m Model) fetchDashboard() tea.Cmd {
	return func() tea.Msg {
		user, err := m.client.GetCurrentUser()
		if err != nil {
			return errMsg{err}
		}
		repos, _ := m.client.ListRepos()
		shortcuts, _ := m.client.ListShortcuts()
		return dashboardDataMsg{user: user, repos: repos, shortcuts: shortcuts}
	}
}

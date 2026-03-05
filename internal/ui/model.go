package ui

import (
	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/ViniZap4/devnook-tui/internal/config"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type view int

const (
	loginView view = iota
	dashboardView
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
	m.dashboard = newDashboardModel()

	if cfg.Token != "" {
		m.view = dashboardView
	}

	return m
}

func (m Model) Init() tea.Cmd {
	if m.config.Token != "" {
		return m.validateAndFetch()
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
		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}

	case loginSuccessMsg:
		m.config.Token = msg.token
		m.config.Save()
		m.client.SetToken(msg.token)
		m.dashboard.user = msg.user
		m.view = dashboardView
		return m, m.fetchDashboard()

	case dashboardDataMsg:
		m.dashboard.user = msg.user
		m.dashboard.repos = msg.repos
		m.dashboard.shortcuts = msg.shortcuts
		m.dashboard.loading = false
		return m, nil

	case tokenInvalidMsg:
		// Token expired/invalid — show login
		m.config.Token = ""
		m.config.Save()
		m.client.SetToken("")
		m.view = loginView
		return m, checkSetup(m.client)

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
type loginSuccessMsg struct {
	token string
	user  *api.User
}
type errMsg struct{ err error }
type tokenInvalidMsg struct{}

type dashboardDataMsg struct {
	user      *api.User
	repos     []api.Repo
	shortcuts []api.Shortcut
}

// validateAndFetch checks if the saved token is still valid, then fetches dashboard data.
func (m Model) validateAndFetch() tea.Cmd {
	return func() tea.Msg {
		user, err := m.client.GetCurrentUser()
		if err != nil {
			return tokenInvalidMsg{}
		}
		repos, _ := m.client.ListRepos()
		shortcuts, _ := m.client.ListShortcuts()
		return dashboardDataMsg{user: user, repos: repos, shortcuts: shortcuts}
	}
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

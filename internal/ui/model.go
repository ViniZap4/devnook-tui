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
	repoDetailView
	issueListView
	issueDetailView
	fileViewerView
	helpView
)

type Model struct {
	client      *api.Client
	config      *config.Config
	view        view
	prevView    view
	width       int
	height      int
	login       loginModel
	dashboard   dashboardModel
	repoDetail  repoDetailModel
	issueList   issueListModel
	issueDetail issueDetailModel
	fileViewer  fileViewerModel
	err         error
}

func NewModel(client *api.Client, cfg *config.Config) Model {
	m := Model{
		client: client,
		config: cfg,
		view:   loginView,
	}
	m.login = newLoginModel()
	m.dashboard = newDashboardModel()
	m.repoDetail = repoDetailModel{}
	m.issueList = issueListModel{}
	m.issueDetail = issueDetailModel{}
	m.fileViewer = fileViewerModel{}

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

		// Global back key — Escape or Backspace go back
		if m.view != loginView && m.view != dashboardView {
			if key.Matches(msg, keys.Back) {
				return m.goBack()
			}
		}

		// ? for help toggle
		if key.Matches(msg, keys.Help) && m.view != loginView {
			if m.view == helpView {
				m.view = m.prevView
			} else {
				m.prevView = m.view
				m.view = helpView
			}
			return m, nil
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
		m.dashboard.orgs = msg.orgs
		m.dashboard.loading = false
		return m, nil

	case tokenInvalidMsg:
		m.config.Token = ""
		m.config.Save()
		m.client.SetToken("")
		m.view = loginView
		return m, checkSetup(m.client)

	case openRepoMsg:
		m.repoDetail = newRepoDetailModel(msg.owner, msg.name)
		m.view = repoDetailView
		return m, m.fetchRepoDetail(msg.owner, msg.name)

	case repoDetailDataMsg:
		m.repoDetail.repo = msg.repo
		m.repoDetail.branches = msg.branches
		m.repoDetail.commits = msg.commits
		m.repoDetail.tree = msg.tree
		m.repoDetail.loading = false
		return m, nil

	case openIssueListMsg:
		m.issueList = newIssueListModel(msg.owner, msg.repo)
		m.view = issueListView
		return m, m.fetchIssueList(msg.owner, msg.repo)

	case issueListDataMsg:
		m.issueList.issues = msg.issues
		m.issueList.loading = false
		return m, nil

	case openIssueMsg:
		m.issueDetail = newIssueDetailModel(msg.owner, msg.repo, msg.number)
		m.view = issueDetailView
		return m, m.fetchIssueDetail(msg.owner, msg.repo, msg.number)

	case issueDetailDataMsg:
		m.issueDetail.issue = msg.issue
		m.issueDetail.comments = msg.comments
		m.issueDetail.loading = false
		return m, nil

	case openFileMsg:
		m.fileViewer = newFileViewerModel(msg.owner, msg.repo, msg.ref, msg.path, msg.name)
		m.view = fileViewerView
		return m, m.fetchFile(msg.owner, msg.repo, msg.ref, msg.path)

	case fileDataMsg:
		m.fileViewer.content = msg.content
		m.fileViewer.binary = msg.binary
		m.fileViewer.loading = false
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
	case repoDetailView:
		m.repoDetail, cmd = m.repoDetail.Update(msg)
	case issueListView:
		m.issueList, cmd = m.issueList.Update(msg)
	case issueDetailView:
		m.issueDetail, cmd = m.issueDetail.Update(msg)
	case fileViewerView:
		m.fileViewer, cmd = m.fileViewer.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.view {
	case loginView:
		return m.login.View(m.width, m.height)
	case dashboardView:
		return m.dashboard.View(m.width, m.height)
	case repoDetailView:
		return m.repoDetail.View(m.width, m.height)
	case issueListView:
		return m.issueList.View(m.width, m.height)
	case issueDetailView:
		return m.issueDetail.View(m.width, m.height)
	case fileViewerView:
		return m.fileViewer.View(m.width, m.height)
	case helpView:
		return renderHelp(m.width, m.height)
	default:
		return "unknown view"
	}
}

func (m Model) goBack() (Model, tea.Cmd) {
	switch m.view {
	case repoDetailView:
		m.view = dashboardView
	case issueListView:
		m.view = repoDetailView
	case issueDetailView:
		m.view = issueListView
	case fileViewerView:
		m.view = repoDetailView
	default:
		m.view = dashboardView
	}
	return m, nil
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
	orgs      []api.Organization
}

type openRepoMsg struct {
	owner, name string
}

type repoDetailDataMsg struct {
	repo     *api.Repo
	branches []api.Branch
	commits  []api.Commit
	tree     []api.TreeEntry
}

type openIssueListMsg struct {
	owner, repo string
}

type issueListDataMsg struct {
	issues []api.Issue
}

type openIssueMsg struct {
	owner, repo string
	number      int
}

type issueDetailDataMsg struct {
	issue    *api.Issue
	comments []api.IssueComment
}

type openFileMsg struct {
	owner, repo, ref, path, name string
}

type fileDataMsg struct {
	content string
	binary  bool
}

// Fetch commands

func (m Model) validateAndFetch() tea.Cmd {
	return func() tea.Msg {
		user, err := m.client.GetCurrentUser()
		if err != nil {
			return tokenInvalidMsg{}
		}
		repos, _ := m.client.ListRepos()
		shortcuts, _ := m.client.ListShortcuts()
		orgs, _ := m.client.ListOrgs()
		return dashboardDataMsg{user: user, repos: repos, shortcuts: shortcuts, orgs: orgs}
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
		orgs, _ := m.client.ListOrgs()
		return dashboardDataMsg{user: user, repos: repos, shortcuts: shortcuts, orgs: orgs}
	}
}

func (m Model) fetchRepoDetail(owner, name string) tea.Cmd {
	return func() tea.Msg {
		repo, err := m.client.GetRepo(owner, name)
		if err != nil {
			return errMsg{err}
		}
		branches, _ := m.client.GetBranches(owner, name)
		commits, _ := m.client.GetCommits(owner, name)
		tree, _ := m.client.GetTree(owner, name, repo.DefaultBranch, "")
		return repoDetailDataMsg{repo: repo, branches: branches, commits: commits, tree: tree}
	}
}

func (m Model) fetchIssueList(owner, repo string) tea.Cmd {
	return func() tea.Msg {
		issues, err := m.client.ListIssues(owner, repo, "open")
		if err != nil {
			return errMsg{err}
		}
		return issueListDataMsg{issues: issues}
	}
}

func (m Model) fetchIssueDetail(owner, repo string, number int) tea.Cmd {
	return func() tea.Msg {
		issue, err := m.client.GetIssue(owner, repo, number)
		if err != nil {
			return errMsg{err}
		}
		comments, _ := m.client.GetIssueComments(owner, repo, number)
		return issueDetailDataMsg{issue: issue, comments: comments}
	}
}

func (m Model) fetchFile(owner, repo, ref, path string) tea.Cmd {
	return func() tea.Msg {
		blob, err := m.client.GetBlob(owner, repo, ref, path)
		if err != nil {
			return errMsg{err}
		}
		return fileDataMsg{content: blob.Content, binary: blob.Binary}
	}
}

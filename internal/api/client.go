package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

// Request types

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// Response types

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
	IsAdmin   bool   `json:"is_admin"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Repo struct {
	ID            int64    `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	IsPrivate     bool     `json:"is_private"`
	IsFork        bool     `json:"is_fork"`
	DefaultBranch string   `json:"default_branch"`
	Topics        []string `json:"topics"`
	StarsCount    int      `json:"stars_count"`
	ForksCount    int      `json:"forks_count"`
	OrgID         *int64   `json:"org_id"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type Shortcut struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
	Color   string `json:"color"`
}

type SetupStatus struct {
	NeedsSetup bool `json:"needs_setup"`
}

type Branch struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	IsHead    bool   `json:"is_head"`
}

type TreeEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Mode string `json:"mode"`
	Size int64  `json:"size"`
}

type BlobContent struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	Content string `json:"content"`
	Binary  bool   `json:"binary"`
}

type Commit struct {
	Hash      string `json:"hash"`
	ShortHash string `json:"short_hash"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Email     string `json:"email"`
	Date      string `json:"date"`
}

type Issue struct {
	ID       int64  `json:"id"`
	RepoID   int64  `json:"repo_id"`
	Number   int    `json:"number"`
	AuthorID int64  `json:"author_id"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	State    string `json:"state"`
}

type IssueComment struct {
	ID       int64  `json:"id"`
	AuthorID int64  `json:"author_id"`
	Author   string `json:"author"`
	Body     string `json:"body"`
}

type Organization struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type PullRequest struct {
	ID         int64  `json:"id"`
	Number     int    `json:"number"`
	Author     string `json:"author"`
	Title      string `json:"title"`
	State      string `json:"state"`
	HeadBranch string `json:"head_branch"`
	BaseBranch string `json:"base_branch"`
}

type Notification struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Read  bool   `json:"read"`
	Link  string `json:"link"`
}

// Auth

func (c *Client) CheckSetup() (bool, error) {
	resp, err := c.get("/api/v1/auth/setup")
	if err != nil {
		return false, err
	}
	var s SetupStatus
	if err := json.Unmarshal(resp, &s); err != nil {
		return false, err
	}
	return s.NeedsSetup, nil
}

func (c *Client) Setup(username, email, password, fullName string) (*AuthResponse, error) {
	resp, err := c.post("/api/v1/auth/setup", RegisterRequest{
		Username: username, Email: email, Password: password, FullName: fullName,
	})
	if err != nil {
		return nil, err
	}
	var auth AuthResponse
	if err := json.Unmarshal(resp, &auth); err != nil {
		return nil, err
	}
	c.token = auth.Token
	return &auth, nil
}

func (c *Client) Login(username, password string) (*AuthResponse, error) {
	resp, err := c.post("/api/v1/auth/login", LoginRequest{Username: username, Password: password})
	if err != nil {
		return nil, err
	}
	var auth AuthResponse
	if err := json.Unmarshal(resp, &auth); err != nil {
		return nil, err
	}
	c.token = auth.Token
	return &auth, nil
}

func (c *Client) Register(username, email, password, fullName string) (*AuthResponse, error) {
	resp, err := c.post("/api/v1/auth/register", RegisterRequest{
		Username: username, Email: email, Password: password, FullName: fullName,
	})
	if err != nil {
		return nil, err
	}
	var auth AuthResponse
	if err := json.Unmarshal(resp, &auth); err != nil {
		return nil, err
	}
	c.token = auth.Token
	return &auth, nil
}

// Users

func (c *Client) GetCurrentUser() (*User, error) {
	resp, err := c.get("/api/v1/users/me")
	if err != nil {
		return nil, err
	}
	var u User
	if err := json.Unmarshal(resp, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

// Repos

func (c *Client) ListRepos() ([]Repo, error) {
	resp, err := c.get("/api/v1/repos")
	if err != nil {
		return nil, err
	}
	var repos []Repo
	if err := json.Unmarshal(resp, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *Client) GetRepo(owner, name string) (*Repo, error) {
	resp, err := c.get(fmt.Sprintf("/api/v1/repos/%s/%s", owner, name))
	if err != nil {
		return nil, err
	}
	var repo Repo
	if err := json.Unmarshal(resp, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}

func (c *Client) CreateRepo(name, description string, isPrivate bool) error {
	_, err := c.post("/api/v1/repos", map[string]any{
		"name": name, "description": description, "is_private": isPrivate,
	})
	return err
}

// Git browsing

func (c *Client) GetTree(owner, name, ref, path string) ([]TreeEntry, error) {
	url := fmt.Sprintf("/api/v1/repos/%s/%s/tree/%s/%s", owner, name, ref, path)
	resp, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var entries []TreeEntry
	if err := json.Unmarshal(resp, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *Client) GetBlob(owner, name, ref, path string) (*BlobContent, error) {
	url := fmt.Sprintf("/api/v1/repos/%s/%s/blob/%s/%s", owner, name, ref, path)
	resp, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var blob BlobContent
	if err := json.Unmarshal(resp, &blob); err != nil {
		return nil, err
	}
	return &blob, nil
}

func (c *Client) GetBranches(owner, name string) ([]Branch, error) {
	resp, err := c.get(fmt.Sprintf("/api/v1/repos/%s/%s/branches", owner, name))
	if err != nil {
		return nil, err
	}
	var branches []Branch
	if err := json.Unmarshal(resp, &branches); err != nil {
		return nil, err
	}
	return branches, nil
}

func (c *Client) GetCommits(owner, name string) ([]Commit, error) {
	resp, err := c.get(fmt.Sprintf("/api/v1/repos/%s/%s/commits", owner, name))
	if err != nil {
		return nil, err
	}
	var commits []Commit
	if err := json.Unmarshal(resp, &commits); err != nil {
		return nil, err
	}
	return commits, nil
}

// Issues

func (c *Client) ListIssues(owner, repo, state string) ([]Issue, error) {
	url := fmt.Sprintf("/api/v1/repos/%s/%s/issues?state=%s", owner, repo, state)
	resp, err := c.get(url)
	if err != nil {
		return nil, err
	}
	var issues []Issue
	if err := json.Unmarshal(resp, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func (c *Client) GetIssue(owner, repo string, number int) (*Issue, error) {
	resp, err := c.get(fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d", owner, repo, number))
	if err != nil {
		return nil, err
	}
	var issue Issue
	if err := json.Unmarshal(resp, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

func (c *Client) GetIssueComments(owner, repo string, number int) ([]IssueComment, error) {
	resp, err := c.get(fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/comments", owner, repo, number))
	if err != nil {
		return nil, err
	}
	var comments []IssueComment
	if err := json.Unmarshal(resp, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Client) CreateIssue(owner, repo, title, body string) error {
	_, err := c.post(fmt.Sprintf("/api/v1/repos/%s/%s/issues", owner, repo),
		map[string]string{"title": title, "body": body})
	return err
}

func (c *Client) CreateIssueComment(owner, repo string, number int, body string) error {
	_, err := c.post(fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/comments", owner, repo, number),
		map[string]string{"body": body})
	return err
}

// Pull Requests

func (c *Client) ListPullRequests(owner, repo, state string) ([]PullRequest, error) {
	resp, err := c.get(fmt.Sprintf("/api/v1/repos/%s/%s/pulls?state=%s", owner, repo, state))
	if err != nil {
		return nil, err
	}
	var prs []PullRequest
	if err := json.Unmarshal(resp, &prs); err != nil {
		return nil, err
	}
	return prs, nil
}

// Organizations

func (c *Client) ListOrgs() ([]Organization, error) {
	resp, err := c.get("/api/v1/orgs")
	if err != nil {
		return nil, err
	}
	var orgs []Organization
	if err := json.Unmarshal(resp, &orgs); err != nil {
		return nil, err
	}
	return orgs, nil
}

// Shortcuts

func (c *Client) ListShortcuts() ([]Shortcut, error) {
	resp, err := c.get("/api/v1/shortcuts")
	if err != nil {
		return nil, err
	}
	var shortcuts []Shortcut
	if err := json.Unmarshal(resp, &shortcuts); err != nil {
		return nil, err
	}
	return shortcuts, nil
}

// Notifications

func (c *Client) ListNotifications() ([]Notification, error) {
	resp, err := c.get("/api/v1/notifications?unread=true")
	if err != nil {
		return nil, err
	}
	var notifications []Notification
	if err := json.Unmarshal(resp, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

// HTTP helpers

func (c *Client) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) post(path string, payload any) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("api error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

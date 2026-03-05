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
	ID            int64  `json:"id"`
	Owner         string `json:"owner"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	IsPrivate     bool   `json:"is_private"`
	DefaultBranch string `json:"default_branch"`
	OrgID         *int64 `json:"org_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
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

func (c *Client) post(path string, payload interface{}) ([]byte, error) {
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

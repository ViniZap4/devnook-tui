package ui

import (
	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loginModel struct {
	inputs    []textinput.Model
	focused   int
	loading   bool
	err       string
	setupMode bool // true = first-run setup, false = login
	checking  bool // checking if setup is needed
}

func newLoginModel() loginModel {
	return loginModel{checking: true}
}

func (m loginModel) initInputs() loginModel {
	if m.setupMode {
		username := textinput.New()
		username.Placeholder = "username"
		username.Focus()

		email := textinput.New()
		email.Placeholder = "email"

		password := textinput.New()
		password.Placeholder = "password"
		password.EchoMode = textinput.EchoPassword
		password.EchoCharacter = '*'

		fullName := textinput.New()
		fullName.Placeholder = "full name"

		m.inputs = []textinput.Model{username, email, password, fullName}
	} else {
		username := textinput.New()
		username.Placeholder = "username"
		username.Focus()

		password := textinput.New()
		password.Placeholder = "password"
		password.EchoMode = textinput.EchoPassword
		password.EchoCharacter = '*'

		m.inputs = []textinput.Model{username, password}
	}
	return m
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

// checkSetupMsg is sent after checking if the server needs first-run setup
type checkSetupMsg struct {
	needsSetup bool
	err        string
}

func checkSetup(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		needs, err := client.CheckSetup()
		if err != nil {
			return checkSetupMsg{err: err.Error()}
		}
		return checkSetupMsg{needsSetup: needs}
	}
}

func (m loginModel) Update(msg tea.Msg, client *api.Client) (loginModel, tea.Cmd) {
	switch msg := msg.(type) {
	case checkSetupMsg:
		m.checking = false
		if msg.err != "" {
			m.err = msg.err
			m.setupMode = false
			m = m.initInputs()
			return m, nil
		}
		m.setupMode = msg.needsSetup
		m = m.initInputs()
		return m, nil

	case tea.KeyMsg:
		if m.checking {
			return m, nil
		}
		switch msg.String() {
		case "tab", "shift+tab":
			if msg.String() == "shift+tab" {
				m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
			} else {
				m.focused = (m.focused + 1) % len(m.inputs)
			}
			for i := range m.inputs {
				if i == m.focused {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, nil
		case "enter":
			if m.loading {
				return m, nil
			}
			if m.setupMode {
				return m.submitSetup(client)
			}
			return m.submitLogin(client)
		}

	case loginErrMsg:
		m.loading = false
		m.err = msg.err
		return m, nil
	}

	if m.checking || len(m.inputs) == 0 {
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

func (m loginModel) submitLogin(client *api.Client) (loginModel, tea.Cmd) {
	username := m.inputs[0].Value()
	password := m.inputs[1].Value()
	if username == "" || password == "" {
		m.err = "username and password are required"
		return m, nil
	}
	m.loading = true
	m.err = ""
	return m, func() tea.Msg {
		token, err := client.Login(username, password)
		if err != nil {
			return loginErrMsg{err: err.Error()}
		}
		return loginSuccessMsg{token: token}
	}
}

func (m loginModel) submitSetup(client *api.Client) (loginModel, tea.Cmd) {
	username := m.inputs[0].Value()
	email := m.inputs[1].Value()
	password := m.inputs[2].Value()
	fullName := m.inputs[3].Value()
	if username == "" || email == "" || password == "" {
		m.err = "username, email, and password are required"
		return m, nil
	}
	m.loading = true
	m.err = ""
	return m, func() tea.Msg {
		token, err := client.Setup(username, email, password, fullName)
		if err != nil {
			return loginErrMsg{err: err.Error()}
		}
		return loginSuccessMsg{token: token}
	}
}

type loginErrMsg struct{ err string }

func (m loginModel) View(width, height int) string {
	title := titleStyle.Render("Dev Nook")

	if m.checking {
		content := lipgloss.JoinVertical(lipgloss.Center,
			title, "", mutedStyle.Render("Connecting to server..."),
		)
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
	}

	var subtitle string
	if m.setupMode {
		subtitle = subtitleStyle.Render("First-run setup — create admin account")
	} else {
		subtitle = subtitleStyle.Render("Sign in to your account")
	}

	var fields string
	labels := []string{"username", "password"}
	if m.setupMode {
		labels = []string{"username", "email", "password", "full name"}
	}
	for i, input := range m.inputs {
		style := inputStyle
		if i == m.focused {
			style = focusedInputStyle
		}
		label := mutedStyle.Render(labels[i])
		fields += label + "\n" + style.Render(input.View()) + "\n"
	}

	var action string
	if m.setupMode {
		action = "[ Create Admin ]"
		if m.loading {
			action = "[ Creating... ]"
		}
	} else {
		action = "[ Login ]"
		if m.loading {
			action = "[ Logging in... ]"
		}
	}

	errText := ""
	if m.err != "" {
		errText = errorStyle.Render(m.err)
	}

	help := mutedStyle.Render("tab: next field  shift+tab: prev  enter: submit  ctrl+c: quit")

	content := lipgloss.JoinVertical(lipgloss.Center,
		title, subtitle, "", fields, action, errText, "", help,
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

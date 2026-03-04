package ui

import (
	"github.com/ViniZap4/devnook-tui/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loginModel struct {
	inputs  []textinput.Model
	focused int
	loading bool
	err     string
}

func newLoginModel() loginModel {
	usernameInput := textinput.New()
	usernameInput.Placeholder = "username"
	usernameInput.Focus()

	passwordInput := textinput.New()
	passwordInput.Placeholder = "password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '*'

	return loginModel{
		inputs: []textinput.Model{usernameInput, passwordInput},
	}
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg, client *api.Client) (loginModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			m.focused = (m.focused + 1) % len(m.inputs)
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

	case loginErrMsg:
		m.loading = false
		m.err = msg.err
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

type loginErrMsg struct{ err string }

func (m loginModel) View(width, height int) string {
	title := titleStyle.Render("Dev Nook")
	subtitle := subtitleStyle.Render("Sign in to your account")

	var fields string
	for i, input := range m.inputs {
		style := inputStyle
		if i == m.focused {
			style = focusedInputStyle
		}
		fields += style.Render(input.View()) + "\n"
	}

	action := "[ Login ]"
	if m.loading {
		action = "[ Logging in... ]"
	}

	errText := ""
	if m.err != "" {
		errText = errorStyle.Render(m.err)
	}

	help := mutedStyle.Render("tab: switch field  enter: submit  ctrl+c: quit")

	content := lipgloss.JoinVertical(lipgloss.Center,
		title, subtitle, "", fields, action, errText, "", help,
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

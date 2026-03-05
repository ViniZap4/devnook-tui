package ui

import "github.com/charmbracelet/lipgloss"

// Theme colors — configurable per-theme
type themeColors struct {
	primary   lipgloss.Color
	secondary lipgloss.Color
	accent    lipgloss.Color
	text      lipgloss.Color
	textDim   lipgloss.Color
	muted     lipgloss.Color
	error     lipgloss.Color
	success   lipgloss.Color
	warning   lipgloss.Color
	border    lipgloss.Color
	surface   lipgloss.Color
	bg        lipgloss.Color
}

var colors = themeColors{
	primary:   lipgloss.Color("#7aa2f7"),
	secondary: lipgloss.Color("#bb9af7"),
	accent:    lipgloss.Color("#7dcfff"),
	text:      lipgloss.Color("#c0caf5"),
	textDim:   lipgloss.Color("#565f89"),
	muted:     lipgloss.Color("#565f89"),
	error:     lipgloss.Color("#f7768e"),
	success:   lipgloss.Color("#9ece6a"),
	warning:   lipgloss.Color("#e0af68"),
	border:    lipgloss.Color("#3b4261"),
	surface:   lipgloss.Color("#1a1b26"),
	bg:        lipgloss.Color("#16161e"),
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.primary).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colors.secondary)

	accentStyle = lipgloss.NewStyle().
			Foreground(colors.accent)

	errorStyle = lipgloss.NewStyle().
			Foreground(colors.error)

	successStyle = lipgloss.NewStyle().
			Foreground(colors.success)

	warningStyle = lipgloss.NewStyle().
			Foreground(colors.warning)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colors.muted)

	textStyle = lipgloss.NewStyle().
			Foreground(colors.text)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.border).
			Padding(1, 2)

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.primary).
			Background(lipgloss.Color("#292e42")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(colors.muted).
				Padding(0, 2)

	tabGapStyle = lipgloss.NewStyle().
			Foreground(colors.border)

	// List styles
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colors.primary).
				Bold(true)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(colors.text)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(colors.textDim)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(colors.textDim).
			Background(lipgloss.Color("#1a1b26"))

	statusBarActiveStyle = lipgloss.NewStyle().
				Foreground(colors.bg).
				Background(colors.primary).
				Bold(true).
				Padding(0, 1)

	// Input styles
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.muted).
			Padding(0, 1).
			Width(40)

	focusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colors.primary).
				Padding(0, 1).
				Width(40)

	// Repo/item card style
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.border).
			Padding(0, 1)

	cardSelectedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colors.primary).
				Padding(0, 1)

	// Badge styles
	publicBadge = lipgloss.NewStyle().
			Foreground(colors.success).
			SetString("public")

	privateBadge = lipgloss.NewStyle().
			Foreground(colors.warning).
			SetString("private")
)

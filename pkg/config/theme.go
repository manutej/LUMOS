package config

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the application
type Theme struct {
	Name       string
	Background string
	Text       string
	Accent     string
	Muted      string
	Warning    string
	Success    string
	Error      string
}

// DarkTheme is the default dark mode theme (VSCode Dark+)
var DarkTheme = Theme{
	Name:       "dark",
	Background: "#1e1e1e",
	Text:       "#e0e0e0",
	Accent:     "#61afef",
	Muted:      "#858585",
	Warning:    "#e06c75",
	Success:    "#98c379",
	Error:      "#f48771",
}

// LightTheme is an alternative light mode theme
var LightTheme = Theme{
	Name:       "light",
	Background: "#ffffff",
	Text:       "#383838",
	Accent:     "#0184bc",
	Muted:      "#a0a0a0",
	Warning:    "#c33e1a",
	Success:    "#0b7c15",
	Error:      "#d1394d",
}

// GetTheme returns a theme by name
func GetTheme(name string) Theme {
	switch name {
	case "light":
		return LightTheme
	default:
		return DarkTheme
	}
}

// Styles contains all UI styles for a theme
type Styles struct {
	Theme           Theme
	Background      lipgloss.Style
	Text            lipgloss.Style
	Accent          lipgloss.Style
	Muted           lipgloss.Style
	Warning         lipgloss.Style
	Success         lipgloss.Style
	Error           lipgloss.Style
	PaneBorder      lipgloss.Style
	PaneTitle       lipgloss.Style
	SelectedItem    lipgloss.Style
	UnselectedItem  lipgloss.Style
	SearchMatch     lipgloss.Style
	LineNumber      lipgloss.Style
	StatusBar       lipgloss.Style
	HelpText        lipgloss.Style
}

// NewStyles creates styles from a theme
func NewStyles(theme Theme) Styles {
	return Styles{
		Theme: theme,
		Background: lipgloss.NewStyle().
			Background(theme.Background).
			Foreground(theme.Text),

		Text: lipgloss.NewStyle().
			Foreground(theme.Text).
			Background(theme.Background),

		Accent: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Background(theme.Background).
			Bold(true),

		Muted: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Background(theme.Background),

		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Background(theme.Background),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success).
			Background(theme.Background),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			Background(theme.Background),

		PaneBorder: lipgloss.NewStyle().
			Foreground(theme.Muted).
			BorderStyle(lipgloss.NormalBorder()),

		PaneTitle: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Background(theme.Background).
			Bold(true).
			Padding(0, 1),

		SelectedItem: lipgloss.NewStyle().
			Background(theme.Accent).
			Foreground(theme.Background).
			Bold(true).
			Padding(0, 1),

		UnselectedItem: lipgloss.NewStyle().
			Foreground(theme.Text).
			Background(theme.Background).
			Padding(0, 1),

		SearchMatch: lipgloss.NewStyle().
			Background(theme.Success).
			Foreground(theme.Background).
			Bold(true),

		LineNumber: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Background(theme.Background),

		StatusBar: lipgloss.NewStyle().
			Foreground(theme.Text).
			Background(theme.Muted).
			Padding(0, 1),

		HelpText: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Background(theme.Background).
			Italic(true),
	}
}

// Package tui implements the yallaops interactive dashboard.
package tui

import "github.com/charmbracelet/lipgloss"

// Theme is a named set of Lip Gloss styles for the dashboard.
type Theme struct {
	Name        string
	Title       lipgloss.Style
	Border      lipgloss.Style
	StatusOK    lipgloss.Style
	StatusWarn  lipgloss.Style
	StatusError lipgloss.Style
	Muted       lipgloss.Style
}

// Themes maps theme name to Theme, in selection order.
var Themes = map[string]Theme{
	"dark":  darkTheme(),
	"light": lightTheme(),
}

// ThemeNames lists available theme names in a stable order.
var ThemeNames = []string{"dark", "light"}

func darkTheme() Theme {
	return Theme{
		Name:        "dark",
		Title:       lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")),
		Border:      lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4B5563")),
		StatusOK:    lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E")),
		StatusWarn:  lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308")),
		StatusError: lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")),
		Muted:       lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")),
	}
}

func lightTheme() Theme {
	return Theme{
		Name:        "light",
		Title:       lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5B21B6")),
		Border:      lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#D1D5DB")),
		StatusOK:    lipgloss.NewStyle().Foreground(lipgloss.Color("#15803D")),
		StatusWarn:  lipgloss.NewStyle().Foreground(lipgloss.Color("#A16207")),
		StatusError: lipgloss.NewStyle().Foreground(lipgloss.Color("#B91C1C")),
		Muted:       lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")),
	}
}

// StatusStyle returns the style matching a release/environment status string.
func (t Theme) StatusStyle(status string) lipgloss.Style {
	switch status {
	case "RELEASE_STATUS_DEPLOYED", "ENVIRONMENT_STATUS_DEPLOYED":
		return t.StatusOK
	case "RELEASE_STATUS_FAILED", "ENVIRONMENT_STATUS_FAILED":
		return t.StatusError
	case "RELEASE_STATUS_RUNNING", "ENVIRONMENT_STATUS_DEPLOYING", "ENVIRONMENT_STATUS_PENDING":
		return t.StatusWarn
	default:
		return t.Muted
	}
}

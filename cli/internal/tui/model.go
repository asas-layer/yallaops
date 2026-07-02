package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Row is a single release's dashboard row.
type Row struct {
	ID       string
	Service  string
	Version  string
	Status   string
	EnvNames []string
	EnvGood  []string // parallel slice: per-env status
}

// FetchFunc loads the current set of dashboard rows, e.g. from the control
// plane over gRPC. Returning an error surfaces it in the dashboard footer
// rather than crashing the TUI.
type FetchFunc func() ([]Row, error)

type refreshMsg struct {
	rows []Row
	err  error
}

// Model is the Bubble Tea model backing `yallaops dashboard`.
type Model struct {
	theme  Theme
	fetch  FetchFunc
	rows   []Row
	err    error
	cursor int
}

// New builds a dashboard Model using themeName (falls back to "dark" if
// unknown) and fetch to load rows.
func New(themeName string, fetch FetchFunc) Model {
	theme, ok := Themes[themeName]
	if !ok {
		theme = Themes["dark"]
	}
	return Model{theme: theme, fetch: fetch}
}

func (m Model) Init() tea.Cmd {
	return m.refresh()
}

func (m Model) refresh() tea.Cmd {
	return func() tea.Msg {
		rows, err := m.fetch()
		return refreshMsg{rows: rows, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		m.rows, m.err = msg.rows, msg.err
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			return m, m.refresh()
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(m.theme.Title.Render("YallaOps — Releases"))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(m.theme.StatusError.Render(fmt.Sprintf("error loading releases: %v", m.err)))
		b.WriteString("\n")
	} else if len(m.rows) == 0 {
		b.WriteString(m.theme.Muted.Render("no releases yet"))
		b.WriteString("\n")
	}

	for i, row := range m.rows {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		b.WriteString(cursor)
		b.WriteString(fmt.Sprintf("%-24s %-10s ", row.ID, row.Version))
		b.WriteString(m.theme.StatusStyle(row.Status).Render(row.Status))
		b.WriteString("\n")

		for j, env := range row.EnvNames {
			status := ""
			if j < len(row.EnvGood) {
				status = row.EnvGood[j]
			}
			b.WriteString(fmt.Sprintf("      %-10s ", env))
			b.WriteString(m.theme.StatusStyle(status).Render(status))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(m.theme.Muted.Render("↑/↓ select · r refresh · q quit"))
	return b.String()
}

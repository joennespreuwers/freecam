package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joennespreuwers/freecam/internal/watcher"
)

// ── Styles ────────────────────────────────────────────────────────────────────

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	greenDot  = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render("●")
	yellowDot = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render("●")

	logBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	timeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	pidStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
)

// ── Messages ──────────────────────────────────────────────────────────────────

// killMsg is sent from the watcher goroutine when processes are killed.
type killMsg struct {
	results []watcher.KillResult
}

// tickMsg drives the watch loop.
type tickMsg struct{}

func watchCmd(processName string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		results, _ := watcher.FindAndKill(processName)
		return killMsg{results: results}
	}
}

func tickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return tickMsg{}
	}
}

// ── Model ─────────────────────────────────────────────────────────────────────

// Model is the bubbletea application model.
type Model struct {
	version     string
	processName string
	paused      bool
	killCount   int
	logLines    []string
	viewport    viewport.Model
	spinner     spinner.Model
	width       int
	ready       bool
}

// New creates a new Model.
func New(version, processName string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))

	return Model{
		version:     version,
		processName: processName,
		spinner:     s,
	}
}

// ── Init ──────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, watchCmd(m.processName))
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "p":
			m.paused = !m.paused
			if !m.paused {
				cmds = append(cmds, watchCmd(m.processName))
			}
		case "c":
			m.logLines = nil
			m.viewport.SetContent("")
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		headerH := 5
		statusH := 3
		footerH := 2
		logH := msg.Height - headerH - statusH - footerH - 4
		if logH < 3 {
			logH = 3
		}
		if !m.ready {
			m.viewport = viewport.New(msg.Width-4, logH)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - 4
			m.viewport.Height = logH
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case killMsg:
		for _, r := range msg.results {
			entry := fmt.Sprintf("%s  Killed %s %s",
				timeStyle.Render(r.KilledAt.Format("15:04:05")),
				r.ProcessName,
				pidStyle.Render(fmt.Sprintf("(PID %d)", r.PID)),
			)
			m.logLines = append([]string{entry}, m.logLines...)
			m.killCount++
		}
		if m.ready {
			m.viewport.SetContent(strings.Join(m.logLines, "\n"))
			m.viewport.GotoTop()
		}
		if !m.paused {
			cmds = append(cmds, watchCmd(m.processName))
		}

	case tickMsg:
		if !m.paused {
			cmds = append(cmds, watchCmd(m.processName))
		}
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing…"
	}

	// Header
	header := headerStyle.Render(fmt.Sprintf("📷  freecam %s", m.version))

	// Status row
	var statusDot, statusText string
	if m.paused {
		statusDot = yellowDot
		statusText = "PAUSED"
	} else {
		statusDot = greenDot
		statusText = "WATCHING"
	}
	statusRow := fmt.Sprintf("\n  Status:  %s %s   %s %s\n  %s %d times\n",
		statusDot,
		lipgloss.NewStyle().Bold(true).Render(statusText),
		labelStyle.Render("Process:"),
		m.processName,
		labelStyle.Render("Killed:"),
		m.killCount,
	)

	// Spinner (shown when active)
	spinnerStr := ""
	if !m.paused {
		spinnerStr = "  " + m.spinner.View() + " monitoring\n"
	}

	// Event log
	logTitle := "─ Event Log "
	logBox := logBorderStyle.
		Width(m.width - 2).
		Render(logTitle + "\n" + m.viewport.View())

	// Footer
	footer := "\n" + footerStyle.Render("  [Q] Quit   [P] Pause/Resume   [C] Clear log")

	return header + "\n" + statusRow + spinnerStr + "\n" + logBox + footer
}

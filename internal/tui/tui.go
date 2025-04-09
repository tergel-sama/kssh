package tui

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tergel-sama/kssh/internal/models"
	"github.com/tergel-sama/kssh/internal/ssh"
)

type model struct {
	hosts     []models.HostConfig
	selected  int
	shouldRun bool
	width     int
	height    int
}

func initialModel(hosts []models.HostConfig) model {
	return model{
		hosts:    hosts,
		selected: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			m.selected = (m.selected + 1) % len(m.hosts)
		case "k", "up":
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.hosts) - 1
			}
		case "enter":
			m.shouldRun = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Background(lipgloss.Color("236")).
		Padding(0, 1).
		Bold(true).
		Render(" Select SSH Host ")
	b.WriteString(title + "\n\n")
	for i, h := range m.hosts {
		line := fmt.Sprintf("%s → %s@%s", h.Name, h.User, h.Hostname)
		if i == m.selected {
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Padding(0, 1).
				Render(line)
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n↑/↓ or j/k to move · Enter to connect · q to quit\n")
	return b.String()
}

func StartTui(hosts []models.HostConfig) {
	p := tea.NewProgram(initialModel(hosts))
	m, err := p.Run()
	if err != nil {
		log.Fatalf("TUI error: %v", err)
	}

	if mm, ok := m.(model); ok && mm.shouldRun {
		ssh.RunSSH(mm.hosts[mm.selected])
	}
}

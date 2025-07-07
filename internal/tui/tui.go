package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tergel-sama/kssh/internal/configfile"
	"github.com/tergel-sama/kssh/internal/models"
	"github.com/tergel-sama/kssh/internal/ssh"
)

const (
	viewList int = iota
	viewAdd
	viewEdit
	viewConfirmDelete
)

type model struct {
	cfgPath         string
	cfg             models.SSHConfig
	selected        int
	view            int
	form            formModel
	confirmDelete   bool
	width           int
	height          int
	shouldConnect   bool
}

func initialModel(cfgPath string, cfg models.SSHConfig) model {
	return model{
		cfgPath:  cfgPath,
		cfg:      cfg,
		selected: 0,
		view:     viewList,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch m.view {
		case viewList:
			return updateList(msg, m)
		case viewAdd, viewEdit:
			currentView := m.view
			newForm, formCmd := m.form.Update(msg)
			m.form = newForm.(formModel)
			cmd = formCmd
			if m.form.done {
				m.view = viewList
				if m.form.save {
					if currentView == viewAdd {
						m.cfg.Hosts = append(m.cfg.Hosts, m.form.getHost())
					} else {
						m.cfg.Hosts[m.selected] = m.form.getHost()
					}
					configfile.Write(m.cfgPath, m.cfg)
				}
			}
			return m, cmd
		case viewConfirmDelete:
			return updateConfirmDelete(msg, m)
		}
	}
	return m, cmd
}

func (m model) View() string {
	switch m.view {
	case viewList:
		return viewListRender(m)
	case viewAdd, viewEdit:
		return m.form.View()
	case viewConfirmDelete:
		return viewConfirmDeleteRender(m)
	default:
		return ""
	}
}

func Start(cfgPath string, cfg models.SSHConfig) error {
	p := tea.NewProgram(initialModel(cfgPath, cfg))
	m, err := p.Run()
	if err != nil {
		return err
	}

	if m.(model).shouldConnect {
		ssh.RunSSH(m.(model).cfg.Hosts[m.(model).selected])
	}

	return nil
}

func updateList(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "j", "down":
		if len(m.cfg.Hosts) > 0 {
			m.selected = (m.selected + 1) % len(m.cfg.Hosts)
		}
	case "k", "up":
		if len(m.cfg.Hosts) > 0 {
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.cfg.Hosts) - 1
			}
		}
	case "enter":
		if len(m.cfg.Hosts) > 0 {
			m.shouldConnect = true
			return m, tea.Quit
		}
	case "a":
		m.view = viewAdd
		m.form = newFormModel("Add New Host", models.HostConfig{Port: 22})
		return m, m.form.Init()
	case "e":
		if len(m.cfg.Hosts) > 0 {
			m.view = viewEdit
			m.form = newFormModel("Edit Host", m.cfg.Hosts[m.selected])
			return m, m.form.Init()
		}
	case "d":
		if len(m.cfg.Hosts) > 0 {
			m.view = viewConfirmDelete
		}
	}
	return m, nil
}

func viewListRender(m model) string {
	var b strings.Builder

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E7F6F2")).
		Background(lipgloss.Color("#395B64")).
		Padding(0, 1).
		Bold(true).
		Render(" KSSH - Select Host ")
	b.WriteString(title + "\n\n")

	if len(m.cfg.Hosts) == 0 {
		b.WriteString("No hosts found. Press 'a' to add a new host.")
	} else {
		for i, h := range m.cfg.Hosts {
			line := fmt.Sprintf("%s → %s@%s", h.Name, h.User, h.Hostname)
			if i == m.selected {
				line = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#395B64")).
					Background(lipgloss.Color("#E7F6F2")).
					Padding(0, 1).
					Render(line)
			}
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\n↑/↓ or j/k to move · Enter to connect · a to add · e to edit · d to delete · q to quit\n")

	return b.String()
}

func updateConfirmDelete(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.cfg.Hosts = append(m.cfg.Hosts[:m.selected], m.cfg.Hosts[m.selected+1:]...)
		configfile.Write(m.cfgPath, m.cfg)
		m.view = viewList
		if m.selected >= len(m.cfg.Hosts) && len(m.cfg.Hosts) > 0 {
			m.selected = len(m.cfg.Hosts) - 1
		}
	default:
		m.view = viewList
	}
	return m, nil
}

func viewConfirmDeleteRender(m model) string {
	host := m.cfg.Hosts[m.selected]
	return fmt.Sprintf("\nAre you sure you want to delete %s (%s@%s)? (y/n)", host.Name, host.User, host.Hostname)
}

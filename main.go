package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

type HostConfig struct {
	Name     string `yaml:"name"`
	Hostname string `yaml:"hostname"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Key      string `yaml:"key"`
}

type SSHConfig struct {
	Hosts []HostConfig `yaml:"hosts"`
}

type model struct {
	hosts     []HostConfig
	selected  int
	shouldRun bool
	width     int
	height    int
}

func initialModel(hosts []HostConfig) model {
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

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}

func runSSH(host HostConfig) {
	args := []string{}
	if host.Key != "" {
		args = append(args, "-i", expandPath(host.Key))
	}
	if host.Port != 0 {
		args = append(args, "-p", fmt.Sprint(host.Port))
	}

	target := fmt.Sprintf("%s@%s", host.User, host.Hostname)
	args = append(args, target)

	fmt.Printf("🔐 Connecting to %s...\n", target)
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("SSH failed: %v", err)
	}
}

func main() {
	// Load config
	cfgPath := filepath.Join(os.Getenv("HOME"), ".ssh_hosts.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var cfg SSHConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("YAML error: %v", err)
	}

	if len(cfg.Hosts) == 0 {
		fmt.Println("No hosts found.")
		os.Exit(1)
	}

	// Start Bubbletea TUI
	p := tea.NewProgram(model{hosts: cfg.Hosts})
	m, err := p.Run() // ✅ use this instead of Start()
	if err != nil {
		log.Fatalf("TUI error: %v", err)
	}

	if mm, ok := m.(model); ok && mm.shouldRun {
		runSSH(mm.hosts[mm.selected])
	}
}

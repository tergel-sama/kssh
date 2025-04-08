package main

import (
	"bufio"
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

func createDefaultConfig(cfgPath string) error {
	fmt.Println("No configuration file found at:", cfgPath)
	fmt.Print("Would you like to create a default configuration? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		return fmt.Errorf("config creation cancelled by user")
	}

	// Get current user for default config
	currentUser, err := user.Current()
	if err != nil {
		currentUser = &user.User{Username: "user"}
	}

	// Create example configuration
	defaultConfig := SSHConfig{
		Hosts: []HostConfig{
			{
				Name:     "example-server",
				Hostname: "example.com",
				User:     currentUser.Username,
				Port:     22,
				Key:      "~/.ssh/id_rsa",
			},
			{
				Name:     "local-vm",
				Hostname: "192.168.1.100",
				User:     currentUser.Username,
				Port:     22,
			},
		},
	}

	// Create the config file
	yamlData, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %v", err)
	}

	// Ensure directory exists
	configDir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	if err := os.WriteFile(cfgPath, yamlData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Println("✅ Default configuration created at:", cfgPath)
	fmt.Println("Please edit this file to add your actual SSH hosts.")

	return nil
}

func main() {
	// Get config path
	cfgPath := filepath.Join(os.Getenv("HOME"), ".ssh_hosts.yaml")

	// Check if config exists
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config doesn't exist, offer to create it
			if err := createDefaultConfig(cfgPath); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			// Try reading the config again
			data, err = os.ReadFile(cfgPath)
			if err != nil {
				log.Fatalf("Failed to read newly created config: %v", err)
			}
		} else {
			log.Fatalf("Failed to read config: %v", err)
		}
	}

	var cfg SSHConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("YAML error: %v", err)
	}

	if len(cfg.Hosts) == 0 {
		fmt.Println("No hosts found in config. Please add hosts to your config file.")
		os.Exit(1)
	}

	// Start Bubbletea TUI
	p := tea.NewProgram(initialModel(cfg.Hosts))
	m, err := p.Run()
	if err != nil {
		log.Fatalf("TUI error: %v", err)
	}

	if mm, ok := m.(model); ok && mm.shouldRun {
		runSSH(mm.hosts[mm.selected])
	}
}

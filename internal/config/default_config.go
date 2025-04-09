package config

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/tergel-sama/kssh/internal/models"
	"gopkg.in/yaml.v2"
)

func CreateDefaultConfig(cfgPath string) error {
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
	defaultConfig := models.SSHConfig{
		Hosts: []models.HostConfig{
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

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tergel-sama/kssh/internal/config"
	"github.com/tergel-sama/kssh/internal/models"
	"github.com/tergel-sama/kssh/internal/tui"
	"gopkg.in/yaml.v2"
)

func main() {
	// Get config path
	cfgPath := filepath.Join(os.Getenv("HOME"), ".ssh_hosts.yaml")

	// Check if config exists
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Config doesn't exist, offer to create it
			if err := config.CreateDefaultConfig(cfgPath); err != nil {
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

	var cfg models.SSHConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("YAML error: %v", err)
	}

	if len(cfg.Hosts) == 0 {
		fmt.Println("No hosts found in config. Please add hosts to your config file.")
		os.Exit(1)
	}

	// Start Bubbletea TUI
	tui.StartTui(cfg.Hosts)
}

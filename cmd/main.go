package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tergel-sama/kssh/internal/configfile"
	"github.com/tergel-sama/kssh/internal/tui"
)

func main() {
	// Get config path
	cfgPath := filepath.Join(os.Getenv("HOME"), ".kssh_hosts.yaml")

	// Read initial config
	cfg, err := configfile.Read(cfgPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Start the TUI, which now handles all logic
	if err := tui.Start(cfgPath, cfg); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

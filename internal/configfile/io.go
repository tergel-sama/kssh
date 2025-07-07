package configfile

import (
	"fmt"
	"os"

	"github.com/tergel-sama/kssh/internal/models"
	"gopkg.in/yaml.v2"
)

// Read reads the SSH configuration from the given path.
func Read(cfgPath string) (models.SSHConfig, error) {
	var cfg models.SSHConfig

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty config if the file doesn't exist
			return models.SSHConfig{}, nil
		}
		return models.SSHConfig{}, fmt.Errorf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return models.SSHConfig{}, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	return cfg, nil
}

// Write writes the SSH configuration to the given path.
func Write(cfgPath string, cfg models.SSHConfig) error {
	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %v", err)
	}

	if err := os.WriteFile(cfgPath, yamlData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

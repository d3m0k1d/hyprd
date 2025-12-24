package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/d3m0k1d/hyprd/pkg/logger"
)

type Config struct {
	Rules []Rule `toml:"rule"`
}

type Rule struct {
	Name    string   `toml:"name"`
	Trigger string   `toml:"trigger"`
	Actions []string `toml:"actions"`
}

func LoadConfig() (*Config, error) {
	logger := logger.New(false)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Error getting user home directory", "err", err)
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".config", "hyprd")
	configPath := filepath.Join(configDir, "config.toml")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Error("Failed to create config directory", "err", err)
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.Info("Config file not found, creating default", "path", configPath)
		if err := createDefaultConfig(configPath); err != nil {
			logger.Error("Failed to create default config", "err", err)
			return nil, err
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("Failed to read config file", "err", err)
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		logger.Error("Failed to parse config", "err", err)
		return nil, err
	}

	return &cfg, nil
}

func createDefaultConfig(path string) error {
	defaultConfig := `# Hyprd Configuration

[[rule]]
name = "example-rule"
trigger = ""
actions = ["reboot"]

[[rule]]
name = "terminal-rule"
trigger = ""
actions = ["kitty"]
`
	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

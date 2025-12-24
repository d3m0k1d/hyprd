package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	Rules []Rule `toml:"rules"`
}

type Rule struct {
	Name       string   `toml:"name"`
	Trigger    string   `toml:"trigger"`
	Conditions []string `toml:"conditions"`
	Actions    []string `toml:"actions"`
}

func LoadConfig(path string) (*Config, error) {
	data, _ := os.ReadFile(path)
	var cfg Config
	toml.Unmarshal(data, &cfg)
	return &cfg, nil
}

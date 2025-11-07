package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TargetURL          string `yaml:"target_url"`
	CheckInterval      int    `yaml:"check_interval_seconds"`
	TimeoutSeconds     int    `yaml:"timeout_seconds"`
	ExpectedStatusCode int    `yaml:"expected_status_code"`
}

// reads and parses the YAML config file
func Load(path string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

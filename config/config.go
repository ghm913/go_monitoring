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

// validate checks if the config has valid values
func (c *Config) validate() error {
	if c.TargetURL == "" {
		return fmt.Errorf("target URL cannot be empty")
	}
	if c.CheckInterval == 0 {
		return fmt.Errorf("check interval cannot be zero")
	}
	if c.TimeoutSeconds == 0 {
		return fmt.Errorf("timeout cannot be zero")
	}
	if c.ExpectedStatusCode == 0 {
		return fmt.Errorf("expected status code cannot be zero")
	}
	return nil
}

// String returns a string representation of the config
func (c *Config) String() string {
	return fmt.Sprintf("Config:\n  Target URL: %s\n  Check Interval: %ds\n  Timeout: %ds\n  Expected Status: %d",
		c.TargetURL, c.CheckInterval, c.TimeoutSeconds, c.ExpectedStatusCode)
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

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &cfg, nil
}

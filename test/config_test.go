package test

import (
	"ghm913/go_monitoring/config"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := config.Load("../config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Print the loaded configuration
	t.Logf("Loaded configuration:\n"+
		"  Target URL: %s\n"+
		"  Check Interval: %d seconds\n"+
		"  Timeout: %d seconds\n"+
		"  Expected Status: %d",
		cfg.TargetURL,
		cfg.CheckInterval,
		cfg.TimeoutSeconds,
		cfg.ExpectedStatusCode)

	// Verify config
	if cfg.TargetURL == "" ||
		cfg.CheckInterval == 0 ||
		cfg.TimeoutSeconds == 0 ||
		cfg.ExpectedStatusCode == 0 {
		t.Error("Configuration values should not be empty or zero")
	}
}

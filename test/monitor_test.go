package test

import (
	"os"
	"testing"

	"ghm913/go_monitoring/config"
	"ghm913/go_monitoring/services"
)

func TestMonitor(t *testing.T) {
	// Clean up log file after test
	defer os.RemoveAll("logs")

	cfg := &config.Config{
		TargetURL:          "http://localhost:1", // This will fail
		CheckInterval:      1,
		TimeoutSeconds:     1,
		ExpectedStatusCode: 200,
	}

	mon := services.NewMonitor(cfg)
	if mon == nil {
		t.Fatal("NewMonitor returned nil")
	}
	defer mon.Close()

	// Initial logs should be empty
	initialLogs := mon.GetRecentLogs(5)
	t.Logf("Initial logs count: %d", len(initialLogs))

	mon.CheckEndpoint() // generate a log

	// Get and verify logs
	logs := mon.GetRecentLogs(5)
	t.Logf("Found %d logs after checking endpoint", len(logs))

	for i, log := range logs {
		t.Logf("Log %d: URL=%s, Expected=%d, Status=%d, Error=%s, Time=%v",
			i+1, log.URL, log.ExpectedStatusCode, log.StatusCode, log.Error, log.Timestamp)
	}

	if len(logs) == 0 {
		t.Error("Expected at least one log after failed check")
	}
}

package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ghm913/go_monitoring/config"
)

// log structure
type FailureLog struct {
	Timestamp          time.Time `json:"timestamp"`
	URL                string    `json:"url"`
	ExpectedStatusCode int       `json:"expected_status_code"`
	StatusCode         int       `json:"status_code"`
	Error              string    `json:"error,omitempty"`
}

// webpage monitoring info strcuture
type Monitor struct {
	cfg     *config.Config
	mu      sync.RWMutex
	logs    []FailureLog
	logFile *os.File
}

func NewMonitor(cfg *config.Config) *Monitor {
	logDir := "logs"
	os.MkdirAll(logDir, 0755)

	logPath := filepath.Join(logDir, "monitoring.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		file = nil
	}

	return &Monitor{
		cfg:     cfg,
		logs:    make([]FailureLog, 0, 100),
		logFile: file,
	}
}

// Start monitoring
func (m *Monitor) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(m.cfg.CheckInterval) * time.Second)
	defer ticker.Stop()
	defer m.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.CheckEndpoint()
		}
	}
}

// Close flushes remaining logs and closes the file
func (m *Monitor) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.logFile != nil {
		for _, entry := range m.logs {
			if data, err := json.Marshal(entry); err == nil {
				m.logFile.Write(append(data, '\n'))
			}
		}
		m.logFile.Sync()
		m.logFile.Close()
	}
}

// the most recent failure logs
func (m *Monitor) GetRecentLogs(limit int) []FailureLog {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.logs) <= limit {
		result := make([]FailureLog, len(m.logs))
		copy(result, m.logs)
		return result
	}

	n := min(len(m.logs), limit)
	start := len(m.logs) - n
	result := make([]FailureLog, limit)
	copy(result, m.logs[start:])
	return result
}

// a single request
func (m *Monitor) CheckEndpoint() {
	client := &http.Client{
		Timeout: time.Duration(m.cfg.TimeoutSeconds) * time.Second,
	}

	resp, err := client.Get(m.cfg.TargetURL)
	if err != nil {
		m.logFailure(0, err.Error())
		RecordRequest(false)
		return
	}
	defer resp.Body.Close()

	isSuccess := resp.StatusCode == m.cfg.ExpectedStatusCode

	RecordRequest(isSuccess)

	if !isSuccess {
		m.logFailure(resp.StatusCode, resp.Status)
	}
}

// stores a new failure log
func (m *Monitor) logFailure(status int, errMsg string) {
	entry := FailureLog{
		Timestamp:          time.Now(),
		URL:                m.cfg.TargetURL,
		ExpectedStatusCode: m.cfg.ExpectedStatusCode,
		StatusCode:         status,
		Error:              errMsg,
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If we have 10 or more logs, write oldest to file
	if len(m.logs) >= 10 {
		oldEntry := m.logs[0]
		m.logs = m.logs[1:]

		if m.logFile != nil {
			if data, err := json.Marshal(oldEntry); err == nil {
				m.logFile.Write(append(data, '\n'))
				m.logFile.Sync()
			}
		}
	}

	m.logs = append(m.logs, entry)
}

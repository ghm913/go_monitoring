package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"ghm913/go_monitoring/config"
)

// log structure
type FailureLog struct {
	Timestamp  time.Time `json:"timestamp"`
	StatusCode int       `json:"status_code"`
	Error      string    `json:"error,omitempty"`
}

// webpage monitoring info strcuture
type Monitor struct {
	cfg  *config.Config
	mu   sync.RWMutex
	logs []FailureLog
}

func NewMonitor(cfg *config.Config) *Monitor {
	return &Monitor{
		cfg:  cfg,
		logs: make([]FailureLog, 0, 100),
	}
}

// Start monitoring
func (m *Monitor) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(m.cfg.CheckInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.CheckEndpoint()
		}
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
	log := FailureLog{
		Timestamp:  time.Now(),
		StatusCode: status,
		Error:      errMsg,
	}

	m.mu.Lock()
	m.logs = append(m.logs, log)
	m.mu.Unlock()
}

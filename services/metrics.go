package services

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// total requests with timestamp
	RequestsTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "webpage",
		Name:      "requests_total",
		Help:      "Total number of requests made to the monitored webpage",
	}, []string{"timestamp"})

	// successful requests with timestamp
	RequestsSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "webpage",
		Name:      "requests_success",
		Help:      "Number of successful requests (matching expected status code)",
	}, []string{"timestamp"})

	// total availability percentage with timestamp
	AvailabilityPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "webpage",
		Name:      "availability_percent",
		Help:      "Total availability percentage (0-100) at each timestamp",
	}, []string{"timestamp"})
)

var (
	totalRequests   int = 0
	successRequests int = 0
)

// result of a single request
func RecordRequest(isSuccess bool) {
	totalRequests++
	if isSuccess {
		successRequests++
	}

	// Record with current timestamp
	timestamp := time.Now().Format(time.RFC3339)
	RequestsTotal.WithLabelValues(timestamp).Set(float64(totalRequests))
	RequestsSuccess.WithLabelValues(timestamp).Set(float64(successRequests))

	// Calculate total availability
	availability := float64(successRequests) / float64(totalRequests) * 100
	AvailabilityPercent.WithLabelValues(timestamp).Set(availability)
}

// CalculateAvailability returns the current total availability percentage
func CalculateAvailability() float64 {
	if totalRequests == 0 {
		return 0
	}
	return float64(successRequests) / float64(totalRequests) * 100
}

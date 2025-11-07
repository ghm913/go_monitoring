package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
)

var (
	// total requests
	RequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "webpage",
		Name:      "requests_total",
		Help:      "Total number of requests made to the monitored webpage",
	})

	// successful requests
	RequestsSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "webpage",
		Name:      "requests_success",
		Help:      "Number of successful requests (matching expected status code)",
	})

	// current availability percentage
	AvailabilityPercent = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: "webpage",
			Name:      "availability_percent",
			Help:      "Current availability percentage (0-100)",
		},
		calculateAvailability,
	)
)

// result of a single request
func RecordRequest(isSuccess bool) {
	RequestsTotal.Inc()
	if isSuccess {
		RequestsSuccess.Inc()
	}
}

// calculate current availability
func calculateAvailability() float64 {
	total := &dto.Metric{}
	success := &dto.Metric{}

	if err := RequestsTotal.Write(total); err != nil || total.Counter == nil {
		return 0
	}

	if err := RequestsSuccess.Write(success); err != nil || success.Counter == nil {
		return 0
	}

	totalVal := total.Counter.GetValue()
	if totalVal == 0 {
		return 0
	}

	return (success.Counter.GetValue() / totalVal) * 100
}

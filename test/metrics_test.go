package test

import (
	"ghm913/go_monitoring/services"
	"testing"
)

func TestMetrics(t *testing.T) {
	services.RecordRequest(true)  // success
	services.RecordRequest(false) // failure
	services.RecordRequest(true)  // success

	// Get the availability value
	availability := services.CalculateAvailability()
	// about 66.67% (2/3)
	t.Logf("Calculated availability: %.2f%%", availability)

	if availability <= 0 || availability > 100 {
		t.Errorf("Availability should be between 0 and 100, got %.2f", availability)
	}
}

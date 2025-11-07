package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"

	"ghm913/go_monitoring/services"
)

type Handler struct {
	monitor *services.Monitor
}

func NewHandler(mon *services.Monitor) *Handler {
	return &Handler{monitor: mon}
}

// configures all HTTP routes
func (h *Handler) SetupRoutes(router *gin.Engine) {
	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Availability info
	router.GET("/availability", h.getAvailability)

	// Latest failure logs endpoint
	router.GET("/logs", h.getLogs)
}

// get current availability
func (h *Handler) getAvailability(c *gin.Context) {
	metric := &dto.Metric{}
	if err := services.AvailabilityPercent.Write(metric); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"availability_percent": metric.GetGauge().GetValue(),
	})
}

// most recent failure logs
func (h *Handler) getLogs(c *gin.Context) {
	logs := h.monitor.GetRecentLogs(20)
	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

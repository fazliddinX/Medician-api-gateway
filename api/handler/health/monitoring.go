package health

import (
	"api-gateway/generated/healthAnalytics"
	"api-gateway/pkg/models"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Monitoring struct {
	Logger     *slog.Logger
	Monitoring healthAnalytics.HealthRecommendationsServiceClient
}

// GetMonitoringRealTime
// @Summary Get Real-Time Health Monitoring
// @Description Retrieves real-time health monitoring data
// @Tags Monitoring
// @Accept json
// @Produce json
// @Success 200 {object} healthAnalytics.MonitoringRealTime
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/monitoring/realtime [get]
func (m *Monitoring) GetMonitoringRealTime(c *gin.Context) {
	var req *healthAnalytics.Void

	res, err := m.Monitoring.GetRealtimeHealthMonitoring(context.Background(), req)
	if err != nil {
		m.Logger.Error("Error while calling GetRealtimeHealthMonitoring:", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetMonitoringDailySummary
// @Summary Get Daily Health Summary
// @Description Retrieves a summary of daily health monitoring data
// @Tags Monitoring
// @Accept json
// @Produce json
// @Success 200 {object} healthAnalytics.Monitoring
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/monitoring/daily-summary [get]
func (m *Monitoring) GetMonitoringDailySummary(c *gin.Context) {
	var req *healthAnalytics.Void

	res, err := m.Monitoring.GetDailyHealthSummary(context.Background(), req)
	if err != nil {
		m.Logger.Error("Error while calling GetRealtimeHealthMonitoring:", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetMonitoringWeeklySummary
// @Summary Get Weekly Health Summary
// @Description Retrieves a summary of weekly health monitoring data
// @Tags Monitoring
// @Accept json
// @Produce json
// @Success 200 {object} healthAnalytics.Monitoring
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/monitoring/weekly-summary [get]
func (m *Monitoring) GetMonitoringWeeklySummary(c *gin.Context) {
	var req *healthAnalytics.Void

	res, err := m.Monitoring.GetWeeklyHealthSummary(context.Background(), req)
	if err != nil {
		m.Logger.Error("Error while calling GetRealtimeHealthMonitoring:", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

package health

import (
	"api-gateway/generated/healthAnalytics"
	"api-gateway/pkg/models"
	"api-gateway/pkg/token"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type HealthRecommendImps struct {
	Logger *slog.Logger
	Health healthAnalytics.HealthRecommendationsServiceClient
}

// @Summary Add Health Recommendation
// @Description Generates and adds health recommendations for the user
// @Tags HealthRecommendations
// @Accept json
// @Produce json
// @Param data body healthAnalytics.HealthRecommendationReq true "Health Recommendation Request"
// @Success 200 {object} healthAnalytics.HealthRecommendation
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 401 {object} models.Error "Unauthorized"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/health-recommendations [post]
func (h *HealthRecommendImps) AddHealthRecommendation(c *gin.Context) {
	value, exists := c.Get("claims")
	if !exists {
		h.Logger.Warn("claims missing from context")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		h.Logger.Warn("claims type assertion error")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims type assertion error"})
	}

	var req *healthAnalytics.HealthRecommendationReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("Error in ShouldBindJSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.UserId = claims.ID

	res, err := h.Health.GenerateHealthRecommendations(context.Background(), req)
	if err != nil {
		h.Logger.Error("Error in GenerateHealthRecommendations", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetHealthRecommendation
// @Summary Get Health Recommendation
// @Description Retrieves health recommendations by ID
// @Tags HealthRecommendations
// @Accept json
// @Produce json
// @Param id path string true "Health Recommendation ID"
// @Success 200 {object} healthAnalytics.HealthRecommendation
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/health-recommendations/{id} [get]
func (h *HealthRecommendImps) GetHealthRecommendation(c *gin.Context) {
	id := c.Param("id")

	req := &healthAnalytics.HealthRecommendationID{Id: id}

	res, err := h.Health.GetHealthRecommendations(context.Background(), req)
	if err != nil {
		h.Logger.Error("Error in HealthRecommendations", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetAllHealthRecommendations
// @Summary Get All Health Recommendations
// @Description Retrieves all health recommendations for the authenticated user
// @Tags HealthRecommendations
// @Accept json
// @Produce json
// @Success 200 {object} healthAnalytics.UserHealthRecommendation
// @Failure 401 {object} models.Error "Unauthorized"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/health-recommendations [get]
func (h *HealthRecommendImps) GetAllHealthRecommendations(c *gin.Context) {
	value, exists := c.Get("claims")
	if !exists {
		h.Logger.Warn("claims missing from context")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		h.Logger.Warn("claims type assertion error")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims type assertion error"})
	}

	req := &healthAnalytics.UserID{Id: claims.ID}

	res, err := h.Health.GetAllHealthRecommendations(context.Background(), req)
	if err != nil {
		h.Logger.Error("Error in GetAllHealthRecommendations", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

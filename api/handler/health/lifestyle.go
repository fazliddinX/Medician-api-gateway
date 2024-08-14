package health

import (
	"api-gateway/generated/healthAnalytics"
	"api-gateway/pkg/models"
	"api-gateway/pkg/token"
	"api-gateway/queue"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Lifestyle struct {
	Logger *slog.Logger
	Life   healthAnalytics.LifestyleServiceClient
	Q      queue.Queue
}

// @Summary Add Lifestyle Data
// @Description Adds new lifestyle data for the user
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Param data body healthAnalytics.Lifestyle true "Lifestyle Data"
// @Success 200 {object} healthAnalytics.LifestyleResponse
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 401 {object} models.Error "Unauthorized"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/lifestyle [post]
func (l *Lifestyle) AddLifestyleData(c *gin.Context) {
	value, exists := c.Get("claims")
	if !exists {
		l.Logger.Warn("claims missing from context")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		l.Logger.Warn("claims type assertion error")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims type assertion error"})
	}

	var req *healthAnalytics.Lifestyle

	if err := c.ShouldBindJSON(&req); err != nil {
		l.Logger.Error("error in ShouldBindJSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.UserId = claims.ID

	breq, err := json.Marshal(req)
	if err != nil {
		l.Logger.Error("error in json.Marshal", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	message, err := l.Q.Create(breq)
	if err != nil {
		l.Logger.Error("error in Q.Create", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})

		time.Sleep(time.Second * 3)

		res, err := l.Life.AddLifestyleData(context.Background(), req)
		if err != nil {
			l.Logger.Error("error in AddLifestyleData", "error", err.Error())
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, message)
}

// @Summary Get Lifestyle Data
// @Description Retrieves lifestyle data by ID
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Param id path string true "Lifestyle Data ID"
// @Success 200 {object} healthAnalytics.LifestyleResponse
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/lifestyle/{id} [get]
func (l *Lifestyle) GetLifestyleData(c *gin.Context) {
	id := c.Param("id")

	req := &healthAnalytics.LifestyleID{Id: id}

	res, err := l.Life.GetLifestyleData(context.Background(), req)
	if err != nil {
		l.Logger.Error("error in GetLifestyleData", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Update Lifestyle Data
// @Description Updates existing lifestyle data by ID
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Param id path string true "Lifestyle Data ID"
// @Param data body healthAnalytics.UpdateLifestyle true "Updated Lifestyle Data"
// @Success 200 {object} healthAnalytics.LifestyleResponse
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/lifestyle/{id} [put]
func (l *Lifestyle) UpdateLifestyleData(c *gin.Context) {
	id := c.Param("id")

	var req *healthAnalytics.UpdateLifestyle

	if err := c.ShouldBindJSON(&req); err != nil {
		l.Logger.Error("error in ShouldBindJSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.Id = id

	res, err := l.Life.UpdateLifestyleData(context.Background(), req)
	if err != nil {
		l.Logger.Error("error in UpdateLifestyleData", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Delete Lifestyle Data
// @Description Deletes lifestyle data by ID
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Param id path string true "Lifestyle Data ID"
// @Success 200 {object} healthAnalytics.Message
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/lifestyle/{id} [delete]
func (l *Lifestyle) DeleteLifestyleData(c *gin.Context) {
	id := c.Param("id")

	req := &healthAnalytics.LifestyleID{Id: id}

	breq, err := json.Marshal(req)
	if err != nil {
		l.Logger.Error("error in json.Marshal", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	message, err := l.Q.Delete(breq)
	if err != nil {
		l.Logger.Error("error in Q.Delete", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})

		time.Sleep(time.Second * 3)

		res, err := l.Life.DeleteLifestyleData(context.Background(), req)
		if err != nil {
			l.Logger.Error("error in DeleteLifestyleData", "error", err.Error())
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, message)
}

// ListLifestyles
// @Summary List Lifestyle Data
// @Description Lists all lifestyle data with optional filters
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Param limit query string false "Pagination limit"
// @Param offset query string false "Pagination offset"
// @Success 200 {object} healthAnalytics.AllLifestyles
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/lifestyle [get]
func (l *Lifestyle) ListLifestyles(c *gin.Context) {
	var filter models.FilterLifestyle

	if err := c.ShouldBindQuery(&filter); err != nil {
		l.Logger.Error("error in ShouldBindQuery", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	var limit int
	var offset int
	var err error
	if filter.Offset != "" {
		limit, err = strconv.Atoi(filter.Limit)
		if err != nil {
			l.Logger.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}
	}

	if filter.Limit != "" {
		offset, err = strconv.Atoi(filter.Offset)
		if err != nil {
			l.Logger.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}
	}

	req := &healthAnalytics.LifestyleFilter{
		Limit:  int64(limit),
		Offset: int64(offset),
	}

	res, err := l.Life.GetAllLifestyleData(context.Background(), req)
	if err != nil {
		l.Logger.Error("error in GetAllLifestyleData", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

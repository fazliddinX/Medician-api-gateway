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

type WearableDataImpl struct {
	Log      *slog.Logger
	Wearable healthAnalytics.WearableDataClient
	Q        queue.Queue
}

// @Summary Add Wearable Data
// @Description Adds new wearable data
// @Tags WearableData
// @Accept json
// @Produce json
// @Param data body healthAnalytics.WearableDate true "Wearable Data"
// @Success 200 {object} healthAnalytics.WearableDataResponse
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/wearable-data [post]
func (w *WearableDataImpl) AddWearableData(c *gin.Context) {
	value, exists := c.Get("claims")
	if !exists {
		w.Log.Warn("claims missing from context")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		w.Log.Warn("claims type assertion error")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims type assertion error"})
	}

	var req *healthAnalytics.WearableDate

	if err := c.ShouldBindJSON(&req); err != nil {
		w.Log.Error("Error in ShouldBindJSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.UserId = claims.ID

	breq, err := json.Marshal(req)
	if err != nil {
		w.Log.Error("Error in json.Marshal", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	message, err := w.Q.Create(breq)
	if err != nil {
		w.Log.Error("Error in Q.Create", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})

		time.Sleep(time.Second * 3)

		res, err := w.Wearable.AddWearableData(context.Background(), req)
		if err != nil {
			w.Log.Error("Error in AddWearableData", "error", err.Error())
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, message)
}

// @Summary Get Wearable Data
// @Description Retrieves wearable data by ID
// @Tags WearableData
// @Accept json
// @Produce json
// @Param id path string true "Wearable Data ID"
// @Success 200 {object} healthAnalytics.WearableDataResponse
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/wearable-data/{id} [get]
func (w *WearableDataImpl) GetWearableData(c *gin.Context) {
	id := c.Param("id")

	req := &healthAnalytics.WearableDataID{Id: id}

	res, err := w.Wearable.GetWearableData(context.Background(), req)
	if err != nil {
		w.Log.Error("Error in GetWearableData", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Update Wearable Data
// @Description Updates existing wearable data by ID
// @Tags WearableData
// @Accept json
// @Produce json
// @Param id path string true "Wearable Data ID"
// @Param data body healthAnalytics.UpdateWearableDate true "Updated Wearable Data"
// @Success 200 {object} healthAnalytics.WearableDataResponse
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/wearable-data/{id} [put]
func (w *WearableDataImpl) UpdateWearableData(c *gin.Context) {
	id := c.Param("id")

	var req *healthAnalytics.UpdateWearableDate

	if err := c.ShouldBindJSON(&req); err != nil {
		w.Log.Error("Error in ShouldBindJSON", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.Id = id

	res, err := w.Wearable.UpdateWearableData(context.Background(), req)
	if err != nil {
		w.Log.Error("Error in UpdateWearableData", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Delete Wearable Data
// @Description Deletes wearable data by ID
// @Tags WearableData
// @Accept json
// @Produce json
// @Param id path string true "Wearable Data ID"
// @Success 200 {object} healthAnalytics.Message
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/wearable-data/{id} [delete]
func (w *WearableDataImpl) DeleteWearableData(c *gin.Context) {
	id := c.Param("id")

	req := &healthAnalytics.WearableDataID{Id: id}

	breq, err := json.Marshal(req)
	if err != nil {
		w.Log.Error("Error in json.Marshal", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	message, err := w.Q.Delete(breq)
	if err != nil {
		w.Log.Error("Error in Q.Delete", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})

		time.Sleep(time.Second * 3)

		res, err := w.Wearable.DeleteWearableData(context.Background(), req)
		if err != nil {
			w.Log.Error("Error in DeleteWearableData", "error", err.Error())
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}

	c.JSON(http.StatusOK, message)
}

// @Summary List Wearable Data
// @Description Lists all wearable data with optional filters
// @Tags WearableData
// @Accept json
// @Produce json
// @Param limit query string false "Pagination limit"
// @Param offset query string false "Pagination offset"
// @Param device_type query string false "Device Type filter"
// @Success 200 {object} healthAnalytics.AllWearableData
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/wearable-data [get]
func (w *WearableDataImpl) ListWearableData(c *gin.Context) {
	var filter models.Filter

	if err := c.ShouldBindQuery(&filter); err != nil {
		w.Log.Error("Error in ShouldBindQuery", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	var limit int
	var offset int
	var err error

	if filter.Limit != "" && filter.Offset != "" {
		limit, err = strconv.Atoi(filter.Limit)
		if err != nil {
			w.Log.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}

		offset, err = strconv.Atoi(filter.Offset)
		if err != nil {
			w.Log.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}
	}

	req := &healthAnalytics.WearableDataFilter{
		Limit:      int64(limit),
		Offset:     int64(offset),
		DeviceType: filter.FirstName,
	}

	res, err := w.Wearable.GetAllWearableData(context.Background(), req)
	if err != nil {
		w.Log.Error("Error in GetAllWearableData", "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

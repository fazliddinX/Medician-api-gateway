package health

import (
	"api-gateway/generated/healthAnalytics"
	"api-gateway/pkg/models"
	"api-gateway/pkg/token"
	"api-gateway/queue"
	"api-gateway/storage"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type MedicalRecordsImpl struct {
	Q   queue.Queue
	Log *slog.Logger
	Med healthAnalytics.MedicalRecordsServiceClient
	St  storage.Storage
}

// @Summary Add Medical Record
// @Description Adds a new medical record
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Param record body healthAnalytics.AddMedicalRecordRequest true "Medical Record"
// @Success 200 {object} healthAnalytics.MedicalRecord
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/medical-records [post]
func (m *MedicalRecordsImpl) AddMedicalRecord(c *gin.Context) {
	var req *healthAnalytics.AddMedicalRecordRequest

	value, exists := c.Get("claims")
	if !exists {
		m.Log.Warn("claims missing from context")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims not found"})
		return
	}

	claims, ok := value.(*token.Claims)
	if !ok {
		m.Log.Warn("claims type assertion error")
		c.JSON(http.StatusUnauthorized, models.Error{Error: "claims type assertion error"})
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.UserId = claims.ID
	req.Id = uuid.NewString()

	breq, err := json.Marshal(req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	message, err := m.Q.Create(breq)
	if err != nil {

		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})

		time.Sleep(time.Second * 5)

		res, err := m.Med.AddMedicalRecord(context.Background(), req)
		if err != nil {
			m.Log.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
		return
	}
	redisReq := healthAnalytics.MedicalRecord{
		Id:          req.Id,
		UserId:      req.UserId,
		RecordType:  req.RecordType,
		RecordDate:  req.RecordDate,
		Description: req.Description,
		DoctorId:    req.DoctorId,
		Attachments: req.Attachments,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	err = m.St.AddMedicalRecord(context.Background(), &redisReq)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		time.Sleep(time.Second * 4)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message.Message,
		"id":      req.Id,
	})
}

// @Summary Get Medical Record
// @Description Retrieves a medical record by ID
// @Tags MedicalRecords
// @Produce json
// @Param id path string true "Medical Record ID"
// @Success 200 {object} healthAnalytics.MedicalRecord
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/medical-records/{id} [get]
func (m *MedicalRecordsImpl) GetMedicalRecord(c *gin.Context) {
	id := c.Param("id")

	req := &healthAnalytics.MedicalRecordID{Id: id}

	res, err := m.St.GetMedicalRecord(context.Background(), req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
	}
	if res.Id != "" {
		c.JSON(http.StatusOK, res)
		return
	}

	res, err = m.Med.GetMedicalRecord(context.Background(), req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Update Medical Record
// @Description Updates an existing medical record
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Param id path string true "Medical Record ID"
// @Param record body healthAnalytics.UpdateMedicalRecordReq true "Updated Medical Record"
// @Success 200 {object} healthAnalytics.MedicalRecord
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/medical-records/{id} [put]
func (m *MedicalRecordsImpl) UpdateMedicalRecord(c *gin.Context) {
	id := c.Param("id")
	var req *healthAnalytics.UpdateMedicalRecordReq

	if err := c.ShouldBindJSON(&req); err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	req.Id = id

	res, err := m.Med.UpdateMedicalRecord(context.Background(), req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	err = m.St.UpdateMedicalRecord(context.Background(), req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Delete Medical Record
// @Description Deletes a medical record by ID
// @Tags MedicalRecords
// @Produce json
// @Param id path string true "Medical Record ID"
// @Success 200 {object} healthAnalytics.MedicalRecordID
// @Failure 400 {object} models.Error "Internal Server Error"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/medical-records/{id} [delete]
func (m *MedicalRecordsImpl) DeleteMedicalRecord(c *gin.Context) {
	id := c.Param("id")

	req := &healthAnalytics.MedicalRecordID{Id: id}

	breq, err := json.Marshal(req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	message, err := m.Q.Delete(breq)
	if err != nil {

		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})

		time.Sleep(time.Second * 3)

		res, err := m.Med.DeleteMedicalRecord(context.Background(), req)
		if err != nil {
			m.Log.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}

	err = m.St.DeleteMedicalRecord(context.Background(), req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
	}

	c.JSON(http.StatusOK, message)
}

// ListMedicalRecords
// @Summary List Medical Records
// @Description Lists medical records based on filters
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Param offset query string false "Pagination offset"
// @Param limit query string false "Pagination limit"
// @Param description query string false "Description filter"
// @Param doctor_id query string false "Doctor ID filter"
// @Success 200 {object} healthAnalytics.ListMedicalRecord
// @Failure 400 {object} models.Error "Bad Request"
// @Failure 500 {object} models.Error "Internal Server Error"
// @Router /api/medical-records [get]
func (m *MedicalRecordsImpl) ListMedicalRecords(c *gin.Context) {
	var filter models.FilterMedicalRecords

	if err := c.ShouldBindQuery(&filter); err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}

	var limit int
	var offset int
	var err error
	if filter.Offset != "" && filter.Limit != "" {

		limit, err = strconv.Atoi(filter.Limit)
		if err != nil {
			m.Log.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}

		offset, err = strconv.Atoi(filter.Offset)
		if err != nil {
			m.Log.Error("Error binding json: ", "error", err)
			c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
			return
		}
	}

	req := &healthAnalytics.MedicalRecordFilter{
		Limit:       int64(limit),
		Offset:      int64(offset),
		Description: filter.Description,
		DoctorId:    filter.DoctorID,
	}

	res, err := m.Med.ListMedicalRecords(context.Background(), req)
	if err != nil {
		m.Log.Error("Error binding json: ", "error", err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

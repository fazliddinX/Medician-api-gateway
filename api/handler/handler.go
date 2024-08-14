package handler

import (
	"api-gateway/api/handler/health"
	user2 "api-gateway/api/handler/user"
	"api-gateway/pkg/models"
	"api-gateway/queue"
	"api-gateway/storage"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUserProfile(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	CreateUser(c *gin.Context)
}

type MedicalRecords interface {
	AddMedicalRecord(c *gin.Context)
	GetMedicalRecord(c *gin.Context)
	UpdateMedicalRecord(c *gin.Context)
	DeleteMedicalRecord(c *gin.Context)
	ListMedicalRecords(c *gin.Context)
}

type WearableData interface {
	AddWearableData(c *gin.Context)
	GetWearableData(c *gin.Context)
	UpdateWearableData(c *gin.Context)
	DeleteWearableData(c *gin.Context)
	ListWearableData(c *gin.Context)
}

type Lifestyle interface {
	AddLifestyleData(c *gin.Context)
	GetLifestyleData(c *gin.Context)
	UpdateLifestyleData(c *gin.Context)
	DeleteLifestyleData(c *gin.Context)
	ListLifestyles(c *gin.Context)
}

type HealthRecommend interface {
	AddHealthRecommendation(c *gin.Context)
	GetHealthRecommendation(c *gin.Context)
	GetAllHealthRecommendations(c *gin.Context)
}

type Monitoring interface {
	GetMonitoringRealTime(c *gin.Context)
	GetMonitoringDailySummary(c *gin.Context)
	GetMonitoringWeeklySummary(c *gin.Context)
}

type MainHandler interface {
	NewUserHandler() UserHandler
	NewMedicalRecords(q queue.Queue, st storage.Storage) MedicalRecords
	NewWearableData(q queue.Queue) WearableData
	NewLifestyle(q queue.Queue) Lifestyle
	NewMedicalRecommend() HealthRecommend
	NewMonitoring() Monitoring
}

type handlerImpl struct {
	models.ServiceManager
}

func NewHandler(manager models.ServiceManager) MainHandler {
	return &handlerImpl{ServiceManager: manager}
}

func (h *handlerImpl) NewUserHandler() UserHandler {
	return &user2.UserImpl{User: h.User, Log: h.Logger}
}

func (h *handlerImpl) NewMedicalRecords(q queue.Queue, st storage.Storage) MedicalRecords {
	return &health.MedicalRecordsImpl{Q: q, Log: h.Logger, Med: h.MedicalRecords, St: st}
}

func (h *handlerImpl) NewWearableData(q queue.Queue) WearableData {
	return &health.WearableDataImpl{h.Logger, h.WearableDate, q}
}

func (h *handlerImpl) NewLifestyle(q queue.Queue) Lifestyle {
	return &health.Lifestyle{Logger: h.Logger, Life: h.Lifestyle, Q: q}
}

func (h *handlerImpl) NewMedicalRecommend() HealthRecommend {
	return &health.HealthRecommendImps{Logger: h.Logger, Health: h.HealthRecommend}
}

func (h *handlerImpl) NewMonitoring() Monitoring {
	return &health.Monitoring{Logger: h.Logger, Monitoring: h.HealthRecommend}
}

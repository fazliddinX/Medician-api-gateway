package models

import (
	pb "api-gateway/generated/healthAnalytics"
	"api-gateway/generated/users"
	"api-gateway/pkg/config"
	"log/slog"
)

type ServiceManager struct {
	Logger          *slog.Logger
	Cfg             config.Config
	User            users.UserServiceClient
	HealthRecommend pb.HealthRecommendationsServiceClient
	Lifestyle       pb.LifestyleServiceClient
	WearableDate    pb.WearableDataClient
	MedicalRecords  pb.MedicalRecordsServiceClient
}

type Filter struct {
	Limit     string `json:"limit"`
	Offset    string `json:"offset"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
}

type Error struct {
	Error string `json:"error"`
}

type FilterMedicalRecords struct {
	Limit       string `json:"limit"`
	Offset      string `json:"offset"`
	Description string `json:"description"`
	DoctorID    string `json:"doctor_id"`
}

type FilterWearable struct {
	Limit      string `json:"limit"`
	Offset     string `json:"offset"`
	DeviceType string `json:"device_type"`
}

type FilterLifestyle struct {
	Limit  string `json:"limit"`
	Offset string `json:"offset"`
}

type Message struct {
	Message string `json:"message"`
}

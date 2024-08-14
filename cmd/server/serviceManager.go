package server

import (
	pb "api-gateway/generated/healthAnalytics"
	"api-gateway/generated/users"
	"api-gateway/pkg/config"
	"api-gateway/pkg/models"
	"api-gateway/queue"
	"api-gateway/queue/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"log/slog"
)

func NewServiceManager(Cfg config.Config, Logger *slog.Logger) models.ServiceManager {
	userPort := Cfg.USER_HOST + Cfg.USER_SERVER_PORT
	healthPort := Cfg.HEALTH_HOST + Cfg.HEALTH_SERVER_PORT

	log.Println(userPort)
	log.Println(userPort)
	log.Println(userPort)

	userConn, err := grpc.NewClient(userPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		Logger.Error("error in connecting to user service", "error", err)
		log.Fatalf("failed to create user grpc connection: %v", err)
	}

	healthConn, err := grpc.NewClient(healthPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		Logger.Error("error in connecting to health service", "error", err)
		log.Fatalf("failed to create health grpc connection: %v", err)
	}

	user := users.NewUserServiceClient(userConn)
	health := pb.NewHealthRecommendationsServiceClient(healthConn)
	lifestyle := pb.NewLifestyleServiceClient(healthConn)
	wearable := pb.NewWearableDataClient(healthConn)
	medicalRecords := pb.NewMedicalRecordsServiceClient(healthConn)

	manager := models.ServiceManager{
		Logger:          Logger,
		Cfg:             Cfg,
		User:            user,
		HealthRecommend: health,
		Lifestyle:       lifestyle,
		WearableDate:    wearable,
		MedicalRecords:  medicalRecords,
	}

	return manager
}

func NewQueues(manager models.ServiceManager) queue.Queues {
	rabbitConn, err := rabbitmq.ConnectToRabbit(manager.Cfg)
	if err != nil {
		manager.Logger.Error("Error to connect rabbit", "error", err)
		log.Fatal(err)
	}

	medicalRecord := queue.NewMedicalRecordQueue(manager.Logger, rabbitConn)
	lifestyle := queue.NewLifestyleQueue(manager.Logger, rabbitConn)
	wearableData := queue.NewWearableDataQueue(manager.Logger, rabbitConn)

	queues := queue.Queues{
		MedicalRecords: medicalRecord,
		Lifestyle:      lifestyle,
		Wearable:       wearableData,
	}

	return queues
}

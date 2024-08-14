package storage

import (
	pb "api-gateway/generated/healthAnalytics"
	red "api-gateway/storage/redis"
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Storage interface {
	AddMedicalRecord(ctx context.Context, in *pb.MedicalRecord) error
	GetMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.MedicalRecord, error)
	UpdateMedicalRecord(ctx context.Context, in *pb.UpdateMedicalRecordReq) error
	DeleteMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) error
}

func NewStorage(log *slog.Logger, rd *redis.Client) Storage {
	return &red.RedisRepo{Log: log, Redis: rd}
}

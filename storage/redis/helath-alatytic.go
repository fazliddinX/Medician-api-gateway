package redis

import (
	pb "api-gateway/generated/healthAnalytics"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type RedisRepo struct {
	Redis *redis.Client
	Log   *slog.Logger
}

func (r *RedisRepo) AddMedicalRecord(ctx context.Context, in *pb.MedicalRecord) error {
	b, err := json.Marshal(in)
	if err != nil {
		r.Log.Error("error in Decoding", "error", err)
		return err
	}

	status := r.Redis.Set(ctx, in.Id, string(b), 0)

	return status.Err()
}

func (r *RedisRepo) GetMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) (*pb.MedicalRecord, error) {
	med := &pb.MedicalRecord{}

	result, err := r.Redis.Get(ctx, in.Id).Result()
	if err != nil {
		r.Log.Error("error in GetMedicalRecord", "error", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(result), med)
	if err != nil {
		r.Log.Error("error in Unmarshal", "error", err)
		return nil, err
	}

	return med, nil
}

func (r *RedisRepo) UpdateMedicalRecord(ctx context.Context, in *pb.UpdateMedicalRecordReq) error {
	med := &pb.MedicalRecord{}

	result, err := r.Redis.Get(ctx, in.Id).Result()
	if err != nil {
		r.Log.Error("error in GetMedicalRecord", "error", err)
		return err
	}

	err = json.Unmarshal([]byte(result), med)
	if err != nil {
		r.Log.Error("error in Unmarshal", "error", err)
		return err
	}

	r.Redis.Del(ctx, in.Id)

	med.RecordType = in.RecordType
	med.RecordDate = in.RecordDate

	b, err := json.Marshal(med)
	if err != nil {
		r.Log.Error("error in Decoding", "error", err)
		return err
	}

	status := r.Redis.Set(ctx, in.Id, string(b), 0)
	return status.Err()
}

func (r *RedisRepo) DeleteMedicalRecord(ctx context.Context, in *pb.MedicalRecordID) error {
	err := r.Redis.Del(ctx, in.Id).Err()
	if err != nil {
		r.Log.Error("error in DeleteMedicalRecord", "error", err)
		return err
	}

	return nil
}

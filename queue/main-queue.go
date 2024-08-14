package queue

import (
	"api-gateway/pkg/models"
	"api-gateway/queue/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type Queues struct {
	MedicalRecords Queue
	Lifestyle      Queue
	Wearable       Queue
}

type Queue interface {
	Create(event []byte) (models.Message, error)
	Delete(event []byte) (models.Message, error)
}

func NewMedicalRecordQueue(log *slog.Logger, rabbitConn *amqp.Connection) Queue {
	return &rabbitmq.MedicalRecordsQueue{Log: log, RabbitConn: rabbitConn}
}
func NewLifestyleQueue(log *slog.Logger, rabbitConn *amqp.Connection) Queue {
	return &rabbitmq.LifestyleQueue{Log: log, RabbitConn: rabbitConn}
}
func NewWearableDataQueue(log *slog.Logger, rabbitConn *amqp.Connection) Queue {
	return &rabbitmq.WearableDataQueue{Log: log, RabbitConn: rabbitConn}
}

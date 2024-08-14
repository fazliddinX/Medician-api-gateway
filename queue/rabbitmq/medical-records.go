package rabbitmq

import (
	"api-gateway/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type MedicalRecordsQueue struct {
	RabbitConn *amqp.Connection
	Log        *slog.Logger
}

func (r *MedicalRecordsQueue) Create(event []byte) (models.Message, error) {
	channel, err := r.RabbitConn.Channel()
	if err != nil {
		r.Log.Error("error creating channel", "error", err)
		return models.Message{}, err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"CreateMedicalRecords",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.Log.Error("error creating queue", "error", err)
		return models.Message{}, err
	}

	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        event,
		})
	if err != nil {
		r.Log.Error("error publishing request", "error", err)
		return models.Message{}, err
	}

	return models.Message{Message: "Medical Record Created"}, nil
}

func (r *MedicalRecordsQueue) Delete(event []byte) (models.Message, error) {
	channel, err := r.RabbitConn.Channel()
	if err != nil {
		r.Log.Error("error creating channel", "error", err)
		return models.Message{}, err
	}

	queue, err := channel.QueueDeclare(
		"DeleteMedicalRecords",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.Log.Error("error creating queue", "error", err)
		return models.Message{}, err
	}

	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        event,
		})
	if err != nil {
		r.Log.Error("error publishing request", "error", err)
		return models.Message{}, err
	}

	return models.Message{Message: "Health Record Deleted"}, nil
}

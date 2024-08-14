package rabbitmq

import (
	"api-gateway/pkg/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type WearableDataQueue struct {
	RabbitConn *amqp.Connection
	Log        *slog.Logger
}

func (r *WearableDataQueue) Create(event []byte) (models.Message, error) {
	channel, err := r.RabbitConn.Channel()
	if err != nil {
		r.Log.Error("error creating channel", "error", err)
		return models.Message{}, err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"CreateWearableData",
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

	return models.Message{Message: "User Created"}, nil
}

func (r *WearableDataQueue) Delete(event []byte) (models.Message, error) {
	channel, err := r.RabbitConn.Channel()
	if err != nil {
		r.Log.Error("error creating channel", "error", err)
		return models.Message{}, err
	}

	queue, err := channel.QueueDeclare(
		"DeleteWearableData",
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

	return models.Message{Message: "Health Record Created"}, nil
}

package rabbitmq

import (
	"api-gateway/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbit(cfg config.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.RABBIT_URL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

package rabbitmq

import (
	"api/pkg/config"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(cfg *config.Config) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s/",
		cfg.RabbitMQConfig.User,
		cfg.RabbitMQConfig.Password,
		cfg.RabbitMQConfig.URL,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQ) Channel() *amqp.Channel {
	return r.ch
}

func (r *RabbitMQ) Close() error {
	if r.ch != nil {
		if err := r.ch.Close(); err != nil {
			log.Printf("Error closing RabbitMQ channel: %v", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("Error closing RabbitMQ connection: %v", err)
			return err
		}
	}
	return nil
}

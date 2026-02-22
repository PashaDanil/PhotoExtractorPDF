package rabbitmq

// TODO: переделать клиент rabbitmq

import (
	"api/internal/config"
	"errors"
	"fmt"

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
		// обработать ошибку
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		// обработать ошибку
		return nil, err
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
	var err error

	if r.ch != nil {
		if e := r.ch.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}

	if r.conn != nil {
		if e := r.conn.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}

	// обработать ошибку

	return err
}

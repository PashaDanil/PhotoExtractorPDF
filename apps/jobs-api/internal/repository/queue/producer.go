package queue

// TODO: переделать клиент rabbitmq

import (
	"api/internal/model/domain"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch     *amqp.Channel
	logger *slog.Logger
}

func NewPublisher(ch *amqp.Channel, log *slog.Logger) *Publisher {
	return &Publisher{ch: ch, logger: log}
}

func (p *Publisher) PublishJob(ctx context.Context, jb domain.Job) error {
	const op = "queue.Publisher.PublishJob"

	msg := jb.ToTask()

	p.logger.Debug("marshalling job message", "op", op, "job_id", msg.JobID)

	body, err := json.Marshal(msg)
	if err != nil {
		p.logger.Debug("failed to marshal job message", "op", op, "job_id", msg.JobID, "error", err)
		return fmt.Errorf("%s: marshal: %w", op, err)
	}

	p.logger.Debug("publishing job to queue", "op", op, "job_id", msg.JobID)

	err = p.ch.PublishWithContext(
		ctx,
		"jobs_exchange",
		"jobs",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		p.logger.Debug("failed to publish job", "op", op, "job_id", msg.JobID, "error", err)
		return fmt.Errorf("%s: publish: %w", op, err)
	}

	p.logger.Debug("job published successfully", "op", op, "job_id", msg.JobID)
	return nil
}

package rabbitmq

import (
	"api/internal/domain/task"
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) *Publisher {
	return &Publisher{ch: ch}
}

func (p *Publisher) PublishJob(ctx context.Context, jobID string, pdfKey string) error {
	msg := task.JobTask{
		JobID:  jobID,
		PDFKey: pdfKey,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

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
		return err
	}

	return nil
}

func (p *Publisher) Close() error {
	return p.ch.Close()
}

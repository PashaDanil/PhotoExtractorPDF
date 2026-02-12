package rabbitmq

func Setup(rmq *RabbitMQ) error {
	ch := rmq.Channel()
	err := ch.ExchangeDeclare(
		"jobs_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"jobs_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return ch.QueueBind(
		q.Name,
		"jobs",
		"jobs_exchange",
		false,
		nil,
	)
}

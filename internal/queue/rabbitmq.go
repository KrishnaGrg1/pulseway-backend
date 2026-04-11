package queue

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "pulseway.checks"

type CheckJob struct {
	MonitorID    int64  `json:"monitor_id"`
	URL          string `json:"url"`
	AttemptCount int    `json:"attempt_count"`
}

type Queue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(url string) (*Queue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable — survives RabbitMQ restart
		false, // not auto-deleted
		false, // not exclusive
		false, // no wait
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &Queue{conn: conn, channel: ch}, nil
}

func (q *Queue) Publish(ctx context.Context, job CheckJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.channel.PublishWithContext(ctx,
		"",        // default exchange
		queueName, // routing key = queue name
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent, // survive restart
		},
	)
}

func (q *Queue) Consume() (<-chan amqp.Delivery, error) {
	return q.channel.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack = false, we manually ack
		false,
		false,
		false,
		nil,
	)
}

func (q *Queue) Close() {
	if err := q.channel.Close(); err != nil {
		log.Println("Error closing channel:", err)
	}
	if err := q.conn.Close(); err != nil {
		log.Println("Error closing connection:", err)
	}
}

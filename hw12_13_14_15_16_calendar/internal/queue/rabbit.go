package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrNotConnected = errors.New("connection isn't exists")

type RabbitMq struct {
	Host      string
	Port      int
	User      string
	Password  string
	QueueName string
	Timeout   time.Duration
	ctx       context.Context
	conn      *amqp.Connection
	ch        *amqp.Channel
	q         amqp.Queue
}

func New(host string, port int, user, pass, queueName string, timeout time.Duration) *RabbitMq {
	return &RabbitMq{
		Host:      host,
		Port:      port,
		User:      user,
		Password:  pass,
		QueueName: queueName,
		Timeout:   timeout,
	}
}

func (r *RabbitMq) Start(ctx context.Context) error {
	r.ctx = ctx
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port),
	)
	if err != nil {
		return err
	}
	r.conn = conn

	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	r.ch = ch

	q, err := ch.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.q = q

	return nil
}

func (r *RabbitMq) Stop(_ context.Context) error {
	err := r.conn.Close()
	if err != nil {
		return err
	}

	err = r.ch.Close()
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMq) Produce(msg []byte, contentType string) error {
	if r.ch == nil {
		return ErrNotConnected
	}

	if contentType == "" {
		contentType = "application/json"
	}

	return r.ch.PublishWithContext(
		r.ctx,
		"",
		r.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        msg,
		},
	)
}

func (r *RabbitMq) Consume() (<-chan amqp.Delivery, error) {
	if r.ch == nil {
		return nil, ErrNotConnected
	}

	messages, consumeErr := r.ch.Consume(
		r.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if consumeErr != nil {
		return nil, consumeErr
	}

	return messages, nil
}

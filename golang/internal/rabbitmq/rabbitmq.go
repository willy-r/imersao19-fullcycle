package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewRabbitClient(connUrl string) (*RabbitClient, error) {
	conn, channel, err := newConnection(connUrl)
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		conn:    conn,
		channel: channel,
		url:     connUrl,
	}, nil
}

func (client *RabbitClient) ConsumeMessages(exchange, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	err := client.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	err = client.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	messages, err := client.channel.Consume(
		queueName,
		"go_app",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume message: %v", err)
	}

	return messages, nil
}

func (client *RabbitClient) PublishMessage(exchange, routingKey, queueName string, message []byte) error {
	err := client.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName,
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = client.channel.QueueBind(queue.Name, routingKey, exchange, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	err = client.channel.Publish(
		exchange, routingKey, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (client *RabbitClient) Close() {
	client.conn.Close()
	client.channel.Close()
}

func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return conn, channel, nil
}

package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TRIP_EXCHANGE = "trip"
)

type MessageHandler func(context.Context, amqp.Delivery) error

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQManager(context context.Context, amqpUrlString string) (*RabbitMQ, error) {
	rabbitmqConfig := NewRabbitMQConfig()
	rmq := &RabbitMQ{}

	if err := rabbitmqConfig.CreateConnection(amqpUrlString); err != nil {
		return nil, fmt.Errorf("failed_to_create_rabbitmq_connection: %v", err)
	}

	rmq.conn = rabbitmqConfig.GetConnection()
	rmq.Channel = rabbitmqConfig.GetChannel()

	go rabbitmqConfig.ReconnectConnection(context, amqpUrlString,
		func(conn *amqp.Connection, channel *amqp.Channel) {
			rmq.conn = conn
			rmq.Channel = channel
			log.Println("rabbitmq_reconnected_successfully")
		},
		func(err error) {
			log.Printf("rabbitmq_reconnection_error:%v", err)
		})

	return rmq, nil
}

func (t *RabbitMQ) PublishingMessage(ctx context.Context, routingKeys string, message any) error {
	dataToSent, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return t.Channel.PublishWithContext(ctx,
		TRIP_EXCHANGE, //Exchange name
		routingKeys,   // routing key
		false,         // mandatory,
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        dataToSent,
		},
	)
}

func (t *RabbitMQ) ConsumeMessage(queueName string, handler MessageHandler) error {
	msgs, err := t.Channel.Consume(
		queueName, // Queue name
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed_to_consume_queue_%s_with_err:%v", queueName, err)
	}

	ctx := context.Background()
	go func() {
		for m := range msgs {
			if err := handler(ctx, m); err != nil {
				log.Printf("get_error_while_listen_to_queue_%s:%v", queueName, err)

				// Nack the message set requeue to false to avoid immediate redelivery loops.
				// Consider a dead-letter exchange (DQl) or a more sophisticated retry mechanism for production

				if neckErr := m.Nack(false, false); neckErr != nil {
					log.Printf("error_failed_to_nack_message:%v", err)
				}
				continue
			}

			// Only ack if the handler success
			if ackErr := m.Ack(false); ackErr != nil {
				log.Printf("failed_to_ack_the_message:%v", err)
			}
		}
	}()

	return nil

}

func (t *RabbitMQ) Close() {
	if t.conn != nil {
		t.conn.Close()
	}

	if t.Channel != nil {
		t.Channel.Close()
	}
}

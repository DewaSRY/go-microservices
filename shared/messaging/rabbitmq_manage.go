package messaging

import (
	"DewaSRY/go-microservices/shared/contracts"
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

func NewRabbitMQManager(amqpUrlString string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(amqpUrlString)

	if err != nil {
		return nil, fmt.Errorf("failed_to_create_rabbitmq_connection: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed_to_create_channel :%v", err)
	}
	rmq := &RabbitMQ{conn: conn, Channel: ch}

	if err := rmq.setupExchangesAndQueues(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed_to_setup_exchanges_and_queues:%v", err)
	}

	return rmq, nil
}

func (t *RabbitMQ) setupExchangesAndQueues() error {
	err := t.Channel.ExchangeDeclare(
		TRIP_EXCHANGE,
		"topic", // type
		true,    //durable
		false,   // auto-deleted
		false,   //internal
		false,   //no wait
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed_to_create_exchange_%s:%v", TRIP_EXCHANGE, err)
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		FindAvailableDriversQueue,
		[]string{
			contracts.TripEventCreated,
			contracts.TripEventDriverNotInterested,
		},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		DriverCmdTripRequestQueue,
		[]string{contracts.DriverCmdTripRequest},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		DriverTripResponseQueue,
		[]string{contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		NotifyDriverNoDriversFoundQueue,
		[]string{contracts.TripEventNoDriversFound},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		NotifyDriverAssignQueue,
		[]string{contracts.TripEventDriverAssigned},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		PaymentTripResponseQueue,
		[]string{contracts.PaymentCmdCreateSession},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		NotifyPaymentSessionCreatedQueue,
		[]string{contracts.PaymentEventSessionCreated},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		TRIP_EXCHANGE,
		NotifyPaymentSuccessQueue,
		[]string{contracts.PaymentEventSuccess},
	); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) declareAndBindingQueue(exchange, queueName string, messageTypes []string) error {
	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments with DLX config
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range messageTypes {
		if err := r.Channel.QueueBind(
			q.Name,   // queue name
			msg,      // routing key
			exchange, // exchange
			false,
			nil,
		); err != nil {
			return fmt.Errorf("failed to bind queue to %s: %v", queueName, err)
		}
	}

	return nil
}

func (t *RabbitMQ) PublishingMessage(ctx context.Context, routingKeys string, message any) error {
	dataToSent, err := json.Marshal(message)

	if err != nil {
		return err
	}

	return t.Channel.PublishWithContext(ctx,
		TRIP_EXCHANGE, //Exchange name
		routingKeys,   // routing key //hallo
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

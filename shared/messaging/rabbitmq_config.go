package messaging

import (
	"DewaSRY/go-microservices/shared/contracts"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HandleReconnectOnSuccess func(conn *amqp.Connection, channel *amqp.Channel)
type HandleReconnectOnError func(err error)

type RabbitMQConfig interface {
	CreateConnection(url string) error
	ReconnectConnection(context context.Context, url string, onSuccess HandleReconnectOnSuccess, onError HandleReconnectOnError)
	GetChannel() *amqp.Channel
	GetConnection() *amqp.Connection
}

type rabbitMQConfig struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	mu              sync.RWMutex
}

// GetChannel implements RabbitMQConfig.
func (t *rabbitMQConfig) GetChannel() *amqp.Channel {
	return t.channel
}

// GetConnection implements RabbitMQConfig.
func (t *rabbitMQConfig) GetConnection() *amqp.Connection {
	return t.conn
}

// Reconnect implements RabbitMQConfig.
func (t *rabbitMQConfig) ReconnectConnection(context context.Context, url string, onSuccess HandleReconnectOnSuccess, onError HandleReconnectOnError) {
	for {

		select {
		case <-context.Done():
			log.Println("stopping_rabbitmq_reconnection_due_to_context_cancellation")
			return
		case err := <-t.notifyChanClose:
			log.Printf("rabbitmq channel closed: %v", err)
		case err := <-t.notifyConnClose:
			log.Printf("rabbitmq connection closed: %v", err)

		}

		amount := 0

		for {
			select {
			case <-context.Done():
				log.Println("stopping_rabbitmq_reconnection_due_to_context_cancellation")
				return
			default:
			}

			log.Println("attempting_to_reconnect_rabbitmq_connection")
			if amount > 0 {
				time.Sleep(5 * time.Second)
			}

			t.mu.Lock()
			err := t.CreateConnection(url)
			t.mu.Unlock()

			if err != nil {
				amount++
				go onError(fmt.Errorf("failed_to_reconnect_rabbitmq_connection: %v", err))
				continue
			} else {
				amount = 0
				log.Println("rabbitmq_reconnected_successfully")
				go onSuccess(t.conn, t.channel)
				break
			}
		}

	}
}

// CreateConnection implements RabbitMQConfig.
func (t *rabbitMQConfig) CreateConnection(url string) error {
	conn, err := amqp.Dial(url)

	if err != nil {
		return fmt.Errorf("failed_to_create_rabbitmq_connection: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed_to_create_channel :%v", err)
	}

	t.notifyChanClose = make(chan *amqp.Error, 1)
	t.notifyConnClose = make(chan *amqp.Error, 1)

	t.conn = conn
	t.channel = ch

	if err := t.setupExchangesAndQueues(ch); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed_to_setup_exchanges_and_queues:%v", err)
	}

	conn.NotifyClose(t.notifyConnClose)
	ch.NotifyClose(t.notifyChanClose)

	return nil
}

func (t *rabbitMQConfig) setupExchangesAndQueues(channel *amqp.Channel) error {
	err := channel.ExchangeDeclare(
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

	//	User
	if err := t.declareAndBindingQueue(
		channel,
		TRIP_EXCHANGE,
		UserEstablishConnectionQueue,
		[]string{
			contracts.UserInitEventProcess,
			contracts.UserCloseConnectiondataEvent,
			contracts.UserDisconnectedProcess,
		},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		channel,
		TRIP_EXCHANGE,
		UserEstablishConnectionNotificationQueue,
		[]string{
			contracts.UserInitSuccessResponse,
			contracts.RouteFoundEvent,
			contracts.DriverActiveResponse,
			contracts.RiderCreateTransactionResponse,
			contracts.TransactionAcceptedResponse,
		},
	); err != nil {
		return err
	}

	// trip flow
	if err := t.declareAndBindingQueue(
		channel,
		TRIP_EXCHANGE,
		TripFlowQueue,
		[]string{
			contracts.TripCreateInitProcess,
			contracts.RiderCreateTripProcess,
			contracts.DriverInitEventProcess,
			contracts.RiderCreateTransactionProcess,
			contracts.DriverAcceptTransactionProcess,
		},
	); err != nil {
		return err
	}

	if err := t.declareAndBindingQueue(
		channel,
		TRIP_EXCHANGE,
		TripFlowNotificationQueue,
		[]string{},
	); err != nil {
		return err
	}

	return nil
}

func (r *rabbitMQConfig) declareAndBindingQueue(channel *amqp.Channel, exchange, queueName string, messageTypes []string) error {
	q, err := channel.QueueDeclare(
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
		if err := channel.QueueBind(
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

func NewRabbitMQConfig() RabbitMQConfig {
	return &rabbitMQConfig{}
}

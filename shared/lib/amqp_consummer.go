package lib

import (
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"encoding/json"
	"log"
)

type QueueConsumer interface {
	Start() error
}

type queueConsumer struct {
	rb        messaging.RabbitMQ
	connMgr   ConnectionManager
	queueName string
}

func NewQueueConsumer(rb messaging.RabbitMQ, connMgr ConnectionManager, queueName string) QueueConsumer {
	return &queueConsumer{
		rb:        rb,
		connMgr:   connMgr,
		queueName: queueName,
	}
}

func (qc *queueConsumer) Start() error {
	msgs, err := qc.rb.Channel.Consume(
		qc.queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var msgBody contracts.AmqpMessage
			if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
				log.Println("Failed to unmarshal message:", err)
				continue
			}

			log.Print(msgBody.OwnerID)
			if len(msgBody.OwnerID) == 0 {

				continue
			}

			var payload any
			if msgBody.Data != nil {
				if err := json.Unmarshal(msgBody.Data, &payload); err != nil {
					log.Println("Failed to unmarshal payload:", err)
					continue
				}
			}

			clientMsg := contracts.WSMessage{
				Type: msg.RoutingKey,
				Data: payload,
			}

			if err := qc.connMgr.Emit(msgBody.OwnerID, clientMsg); err != nil {
				log.Printf("Failed_to_send_message_to_user_%s: %v", msgBody.OwnerID, err)
			}
		}
	}()

	return nil
}

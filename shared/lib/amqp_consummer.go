package lib

import (
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"encoding/json"
)

type OnSuccess func(msg contracts.WSMessage)
type OnError func(err error)

type QueueConsumer interface {
	Start(onsuccess OnSuccess, onError OnError) error
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

func (qc *queueConsumer) Start(onsuccess OnSuccess, onError OnError) error {
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
			var msgBody contracts.MessageData
			if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
				onError(err)
				continue
			}

			if len(msgBody.ConnectionId) == 0 {
				continue
			}

			var payload any
			if msgBody.Data != nil {
				if err := json.Unmarshal(msgBody.Data, &payload); err != nil {
					onError(err)
					continue
				}
			}

			clientMsg := contracts.WSMessage{
				Type: msg.RoutingKey,
				Data: payload,
			}

			if err := qc.connMgr.Emit(msgBody.ConnectionId, clientMsg); err != nil {
				onError(err)
				continue
			}

			onsuccess(clientMsg)
		}
	}()

	return nil
}

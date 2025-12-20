package handler

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/domain"
	"DewaSRY/go-microservices/services/api-gateway/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/lib"
	"DewaSRY/go-microservices/shared/logger"
	"DewaSRY/go-microservices/shared/messaging"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type WsHandler struct {
	connManager lib.ConnectionManager
	rabbitMq    *messaging.RabbitMQ
	tripService domain.RideShareServices
	logger      logger.Logger
}

func NewWsHandler(connManager lib.ConnectionManager, rabbitMq *messaging.RabbitMQ, tripService domain.RideShareServices, logger logger.Logger) *WsHandler {
	return &WsHandler{connManager: connManager, rabbitMq: rabbitMq, tripService: tripService, logger: logger}
}

func (t *WsHandler) WsHandleStartConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := t.connManager.InitUpgrade(w, r)
	if err != nil {
		t.logger.Error("failed_to_start_connection", err, map[string]any{
			"event": "failed_to_start_connection",
		})
		return
	}

	ctx := r.Context()
	connectionId := uuid.New().String()
	t.connManager.Add(connectionId, conn)

	t.logger.Info("new_connection_establish", map[string]any{
		"connectionId": connectionId,
	})

	defer func() {
		t.logger.Info("connection_down", map[string]any{
			"connectionId": connectionId,
		})

		conn.Close()
		t.tripService.UserDisconnected(ctx, connectionId)
		t.connManager.Remove(connectionId)
	}()

	t.makeQueueConsumer(
		[]string{
			messaging.UserEstablishConnectionNotificationQueue,
			messaging.TripFlowNotificationQueue,
		},
	)

	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			t.logger.Error("failed_to_read_rider_connection", err, map[string]any{
				"connection_id": connectionId,
			})
			break
		}

		var messageData types.WsMessageData
		if err := json.Unmarshal(p, &messageData); err != nil {
			t.logger.Error("error_unmarshaling_rider_message", err, map[string]any{
				"data": p,
			})
			continue
		}

		t.logger.Info("received_message_from_connection", map[string]any{
			"connection_id": connectionId,
			"message_type":  messageData.Type,
		})

		switch messageData.Type {
		case contracts.TripCreateInitEvent:
			if err := t.tripService.CreateTripEvent(ctx, connectionId, messageData.Data); err != nil {
				t.logger.Error("failed_to_process_trip_create_init_event", err, map[string]any{
					"connection_id": connectionId,
					"route_type":    messageData.Type,
				})
			}
		case contracts.RiderCreateTripRequest:
			if err := t.tripService.RiderCreateTripRequest(ctx, connectionId, messageData.Data); err != nil {
				t.logger.Error("failed_to_process_rider_create_trip_request", err, map[string]any{
					"connection_id": connectionId,
					"route_type":    messageData.Type,
				})
			}
		case contracts.DriverInitEvent:
			if err := t.tripService.DriverInitRequest(ctx, connectionId, messageData.Data); err != nil {
				t.logger.Error("failed_to_process_driver_init_request", err, map[string]any{
					"connection_id": connectionId,
					"route_type":    messageData.Type,
				})
			}
		case contracts.RiderCreateTransactionRequest:
			if err := t.tripService.RiderCreateTransaction(ctx, connectionId, messageData.Data); err != nil {
				t.logger.Error("failed_to_process_rider_create_transaction", err, map[string]any{
					"connection_id": connectionId,
					"route_type":    messageData.Type,
				})
			}
		case contracts.DriverAcceptTransactionRequest:
			if err := t.tripService.DriverAcceptedTransaction(ctx, connectionId, messageData.Data); err != nil {
				t.logger.Error("failed_to_process_driver_accepted_transaction", err, map[string]any{
					"connection_id": connectionId,
					"route_type":    messageData.Type,
				})
			}
		default:
			t.logger.Warn("trip_received_unknown_messages", map[string]any{
				"data":       messageData.Type,
				"connection": connectionId,
			})
		}
	}
}

func (t *WsHandler) makeQueueConsumer(queueList []string) {
	for _, q := range queueList {
		consumer := lib.NewQueueConsumer(*t.rabbitMq, t.connManager, q)
		if err := consumer.Start(
			t._queueConsumerOnSuccess,
			t._queueConsumerOnError(q),
		); err != nil {
			t.logger.Error("failed_to_start_consumer_for_queue", err, map[string]any{
				"queue": q,
			})
		}
	}
}

func (t *WsHandler) _queueConsumerOnSuccess(msg contracts.WSMessage) {
	t.logger.Info("message_sent_to_user_successfully", map[string]any{
		"message_type": msg.Type,
	})
}

func (t *WsHandler) _queueConsumerOnError(q string) func(err error) {
	return func(err error) {
		t.logger.Error("error_in_queue_consumer", err, map[string]any{
			"queue": q,
		})
	}
}

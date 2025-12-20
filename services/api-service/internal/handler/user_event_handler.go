package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/logger"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
)

type userEventHandler struct {
	rabbitmq    *messaging.RabbitMQ
	userService domain.UserService
	logger      logger.Logger
}

// HandlerUserDisconnect implements domain.UserEventHandler.
func (t *userEventHandler) HandlerUserDisconnect(ctx context.Context, data []byte) {
	var payload messaging.UserDisconnectionRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		t.logger.Error("failed_to_unmarshal_user_disconnection_request", err, map[string]any{
			"data": data,
		})
		return
	}

	if err := t.userService.UserCleanUpData(ctx, payload.ConnectionId); err != nil {
		t.logger.Error("failed_to_cleanup_user_data_on_disconnect", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	if err := t.userService.NotifyDriverActive(ctx); err != nil {
		t.logger.Error("failed_to_notify_driver_active_on_user_disconnect", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	t.logger.Info("user_disconnect_processed_successfully", map[string]any{
		"connection_id": payload.ConnectionId,
	})

}

func NewUserEventHandler(userService domain.UserService, rabbitmq *messaging.RabbitMQ, logger logger.Logger) domain.UserEventHandler {
	return &userEventHandler{userService: userService, rabbitmq: rabbitmq, logger: logger}
}

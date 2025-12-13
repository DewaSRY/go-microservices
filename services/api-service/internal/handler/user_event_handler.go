package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
)

type userEventHandler struct {
	rabbitmq    *messaging.RabbitMQ
	userService domain.UserService
}

// HandlerUserDisconnect implements domain.UserEventHandler.
func (t *userEventHandler) HandlerUserDisconnect(ctx context.Context, data []byte) {
	var payload messaging.UserDisconnectionRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		return
	}

	if err := t.userService.UserCleanUpData(ctx, payload.ConnectionId); err != nil {
		return
	}

	if err := t.userService.NotifyDriverActive(ctx); err != nil {
		return
	}

}

func NewUserEventHandler(userService domain.UserService, rabbitmq *messaging.RabbitMQ) domain.UserEventHandler {
	return &userEventHandler{userService: userService, rabbitmq: rabbitmq}
}

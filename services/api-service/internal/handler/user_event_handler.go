package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"
)

type userEventHandler struct {
	rabbitmq    *messaging.RabbitMQ
	userService domain.UserService
}

// HandlerUserDisconnect implements domain.UserEventHandler.
func (u *userEventHandler) HandlerUserDisconnect(ctx context.Context, data []byte) {
	panic("unimplemented")
}

// HandlerUserInitConnection implements domain.UserEventHandler.
func (t *userEventHandler) HandlerUserInitConnection(ctx context.Context, data []byte) {
	var payload messaging.InitConnectionRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("error_failed_to_init_user_connection:%v", err)
		return
	}

	if err := t.userService.UserInit(ctx, payload); err != nil {
		log.Printf("error_failed_to_init_user_connection:%v", err)
		return
	}

	successResponse, err := json.Marshal(map[string]any{
		"message": "success_init_data",
	})

	if err != nil {
		log.Printf("error_failed_to_init_user_connection:%v", err)
		return
	}

	log.Print("send data to ", payload.ConnectionId)
	if err := t.rabbitmq.PublishingMessage(ctx, contracts.UserInitSuccessResponse, contracts.MessageData{
		ConnectionId: payload.ConnectionId,
		Data:         successResponse,
	}); err != nil {
		log.Printf("error_failed_to_publish_%v", err)
		return
	}

}

func NewUserEventHandler(userService domain.UserService, rabbitmq *messaging.RabbitMQ) domain.UserEventHandler {
	return &userEventHandler{userService: userService, rabbitmq: rabbitmq}
}

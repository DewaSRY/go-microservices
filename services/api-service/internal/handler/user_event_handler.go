package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
)

type userEventHandler struct {
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
		return
	}

	t.userService.UserInit(ctx, payload)

}

func NewUserEventHandler(userService domain.UserService) domain.UserEventHandler {
	return &userEventHandler{userService: userService}
}

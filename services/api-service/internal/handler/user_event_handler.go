package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"context"
)

type userEventHandler struct {
	userService domain.UserService
}

// HandlerUserDisconnect implements domain.UserEventHandler.
func (u *userEventHandler) HandlerUserDisconnect(ctx context.Context, data []byte) {
	panic("unimplemented")
}

// HandlerUserInitConnection implements domain.UserEventHandler.
func (u *userEventHandler) HandlerUserInitConnection(ctx context.Context, data []byte) {
	panic("unimplemented")
}

func NewUserEventHandler(userService domain.UserService) domain.UserEventHandler {
	return &userEventHandler{userService: userService}
}

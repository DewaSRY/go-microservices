package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"
)

type tripFlowEventHandler struct {
	rabbitmq        *messaging.RabbitMQ
	tripFlowService domain.TripFlowService
	userService     domain.UserService
	osrmIntegration domain.OsrmIntegrationService
}

// HandlerTripCreate implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandlerTripCreate(ctx context.Context, data []byte) {
	var payload messaging.CreateTripRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		return
	}

	if err := t.tripFlowService.TripCreate(ctx, payload.RiderId); err != nil {
		return
	}

	if err := t.userService.UpdateRiderLocation(ctx, types.UpdateRiderLocationParam{
		RiderId:     payload.ConnectionId,
		Location:    payload.Pickup,
		Destination: payload.Destination,
	}); err != nil {
		return
	}

	routeValue, err := t.osrmIntegration.GetRoutes(ctx, payload.Pickup, payload.Destination)
	if err != nil {
		return
	}

	routeResponse, err := json.Marshal(messaging.RoutesResponse{
		Coordinate: routeValue.Coordinate,
		Distance:   routeValue.Distance,
		Duration:   routeValue.Duration,
	})

	if err != nil {
		return
	}

	log.Print("send route data to ", payload.ConnectionId)
	if err := t.rabbitmq.PublishingMessage(ctx, contracts.RouteFoundEvent, contracts.MessageData{
		ConnectionId: payload.ConnectionId,
		Data:         routeResponse,
	}); err != nil {
		return
	}

}

func NewTripFlowEventHandler(
	rabbitmq *messaging.RabbitMQ,
	tripFlowService domain.TripFlowService,
	userService domain.UserService,
	osrmIntegration domain.OsrmIntegrationService,
) domain.TripFlowHandler {
	return &tripFlowEventHandler{
		rabbitmq:        rabbitmq,
		tripFlowService: tripFlowService,
		userService:     userService,
		osrmIntegration: osrmIntegration,
	}
}

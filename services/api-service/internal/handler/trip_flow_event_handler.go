package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"time"
)

type tripFlowEventHandler struct {
	rabbitmq        *messaging.RabbitMQ
	tripFlowService domain.TripFlowService
	userService     domain.UserService
	osrmIntegration domain.OsrmIntegrationService
}

// HandlerDriverInit implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandlerDriverInit(ctx context.Context, data []byte) {
	var payload messaging.DriverInitRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		return
	}

	//create driver
	if err := t.userService.CreateDriver(ctx, types.CreateDriverParam{
		ConnectionId: payload.ConnectionId,
		Location:     payload.Location,
		PackageSlug:  payload.PackageSlug,
	}); err != nil {
		return
	}

	if err := t.userService.NotifyDriverActive(ctx); err != nil {
		return
	}
}

// HandlerRiderCreateTrip implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandlerRiderCreateTrip(ctx context.Context, data []byte) {
	var payload messaging.RiderCreateTripRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		return
	}

	if err := t.userService.CreateRider(ctx, types.CreateRiderParam{
		ConnectionId: payload.ConnectionId,
		Location:     payload.Pickup,
		Destination:  payload.Destination,
	}); err != nil {
		return
	}

	if err := t.tripFlowService.CreateRiderTrip(ctx, types.CreateTripParam{
		RiderId: payload.ConnectionId,
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

	if err := t.rabbitmq.PublishingMessage(ctx, contracts.RouteFoundEvent, contracts.MessageData{
		ConnectionId: payload.ConnectionId,
		Data:         routeResponse,
	}); err != nil {
		return
	}

}

// HandlerTripCreate implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandlerTripCreate(ctx context.Context, data []byte) {
	var payload messaging.CreateTripRequest
	if err := json.Unmarshal(data, &payload); err != nil {
		// t.logger.Error("failed to unmarshal CreateTripRequest", "error", err)
		return
	}

	// ---- Goroutine 1: Create Trip & Update Rider Location ----
	go func(p messaging.CreateTripRequest) {
		// each goroutine gets its own context (prevent leaks)
		gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// trip create
		if err := t.tripFlowService.TripCreate(gCtx, p.RiderId); err != nil {
			// t.logger.Error("TripCreate failed", "error", err)
			return

		}

		// update rider location
		if err := t.userService.UpdateRiderLocation(gCtx, types.UpdateRiderLocationParam{
			RiderId:     p.ConnectionId,
			Location:    p.Pickup,
			Destination: p.Destination,
		}); err != nil {
			// t.logger.Error("UpdateRiderLocation failed", "error", err)
			return
		}
	}(payload)

	// ---- Goroutine 2: Routing via OSRM & Publish Event ----
	go func(p messaging.CreateTripRequest) {
		gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		routeValue, err := t.osrmIntegration.GetRoutes(gCtx, p.Pickup, p.Destination)
		if err != nil {
			// t.logger.Error("GetRoutes failed", "error", err)
			return
		}

		routeResponse, err := json.Marshal(messaging.RoutesResponse{
			Coordinate: routeValue.Coordinate,
			Distance:   routeValue.Distance,
			Duration:   routeValue.Duration,
		})
		if err != nil {
			// t.logger.Error("failed to marshal RoutesResponse", "error", err)
			return
		}

		if err := t.rabbitmq.PublishingMessage(
			gCtx,
			contracts.RouteFoundEvent,
			contracts.MessageData{
				ConnectionId: p.ConnectionId,
				Data:         routeResponse,
			},
		); err != nil {
			// t.logger.Error("PublishingMessage failed", "error", err)
			return
		}
	}(payload)
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

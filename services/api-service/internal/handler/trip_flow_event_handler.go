package handler

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/logger"
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
	logger          logger.Logger
}

// HandleDriverAcceptedTransaction implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandleDriverAcceptedTransaction(ctx context.Context, data []byte) {
	var payload messaging.DriverAcceptedTransactionRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		t.logger.Error("failed_to_unmarshal_driver_accepted_transaction_request", err, map[string]any{
			"data": data,
		})
		return
	}

	transactionModel, err := t.tripFlowService.DriverAcceptedTransaction(ctx, payload.TransactionId)
	if err != nil {
		t.logger.Error("failed_to_process_driver_accepted_transaction", err, map[string]any{
			"transaction_id": payload.TransactionId,
		})
		return
	}

	if err := t.userService.NotifyDriverAcceptedTransaction(
		ctx, payload.TransactionId, transactionModel.DriverId, transactionModel.RiderId); err != nil {
		t.logger.Error("failed_to_notify_driver_accepted_transaction", err, map[string]any{
			"transaction_id": payload.TransactionId,
		})
		return
	}

	t.logger.Info("driver_accepted_transaction_processed_successfully", map[string]interface{}{
		"transaction_id": payload.TransactionId,
	})
}

// HandleRiderCreateTransaction implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandleRiderCreateTransaction(ctx context.Context, data []byte) {
	var payload messaging.RiderCreateTransactionRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		t.logger.Error("failed_to_unmarshal_rider_create_transaction_request", err, map[string]any{
			"data": data,
		})
		return
	}

	transactionId, err := t.userService.RiderStartTransaction(ctx, payload.ConnectionId, payload.DriverId)
	if err != nil {
		t.logger.Error("failed_to_process_rider_create_transaction", err, map[string]any{
			"connection_id": payload.ConnectionId,
			"driver_id":     payload.DriverId,
		})
		return
	}

	if err := t.userService.DriverNotifyTransaction(ctx, payload.DriverId, transactionId); err != nil {
		t.logger.Error("failed_to_notify_driver_of_transaction", err, map[string]any{
			"driver_id":      payload.DriverId,
			"transaction_id": transactionId,
		})
		return
	}

	if err := t.userService.RiderNotifyTransaction(ctx, payload.ConnectionId, transactionId); err != nil {
		t.logger.Error("failed_to_notify_rider_of_transaction", err, map[string]any{
			"connection_id":  payload.ConnectionId,
			"transaction_id": transactionId,
		})
		return
	}

	t.logger.Info("rider_create_transaction_processed_successfully", map[string]interface{}{
		"transaction_id": transactionId,
	})

}

// HandlerDriverInit implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandlerDriverInit(ctx context.Context, data []byte) {
	var payload messaging.DriverInitRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		t.logger.Error("failed_to_unmarshal_driver_init_request", err, map[string]any{
			"data": data,
		})
		return
	}

	//create driver
	if err := t.userService.CreateDriver(ctx, types.CreateDriverParam{
		ConnectionId: payload.ConnectionId,
		Location:     payload.Location,
		PackageSlug:  payload.PackageSlug,
	}); err != nil {
		t.logger.Error("failed_to_create_driver", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	if err := t.userService.NotifyDriverActive(ctx); err != nil {
		t.logger.Error("failed_to_notify_driver_active", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	t.logger.Info("driver_init_processed_successfully", map[string]interface{}{
		"connection_id": payload.ConnectionId,
	})

}

// HandlerRiderCreateTrip implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandlerRiderCreateTrip(ctx context.Context, data []byte) {
	var payload messaging.RiderCreateTripRequest

	if err := json.Unmarshal(data, &payload); err != nil {
		t.logger.Error("failed_to_unmarshal_rider_create_trip_request", err, map[string]any{
			"data": data,
		})
		return
	}

	if err := t.userService.CreateRider(ctx, types.CreateRiderParam{
		ConnectionId: payload.ConnectionId,
		Location:     payload.Pickup,
		Destination:  payload.Destination,
	}); err != nil {
		t.logger.Error("failed_to_create_rider", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	if err := t.tripFlowService.CreateRiderTrip(ctx, types.CreateTripParam{
		RiderId: payload.ConnectionId,
	}); err != nil {
		t.logger.Error("failed_to_create_rider_trip", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	routeValue, err := t.osrmIntegration.GetRoutes(ctx, payload.Pickup, payload.Destination)
	if err != nil {
		t.logger.Error("failed_to_get_routes_from_osrm", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	routeResponse, err := json.Marshal(messaging.RoutesResponse{
		Coordinate: routeValue.Coordinate,
		Distance:   routeValue.Distance,
		Duration:   routeValue.Duration,
	})

	if err != nil {
		t.logger.Error("failed_to_marshal_routes_response", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	if err := t.rabbitmq.PublishingMessage(ctx, contracts.RouteFoundEvent, contracts.MessageData{
		ConnectionId: payload.ConnectionId,
		Data:         routeResponse,
	}); err != nil {
		t.logger.Error("failed_to_publish_route_found_event", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	if err := t.userService.NotifyDriverActive(ctx); err != nil {
		t.logger.Error("failed_to_notify_driver_active", err, map[string]any{
			"connection_id": payload.ConnectionId,
		})
		return
	}

	t.logger.Info("rider_create_trip_processed_successfully", map[string]interface{}{
		"connection_id": payload.ConnectionId,
	})

}

// HandlerTripCreate implements domain.TripFlowHandler.
func (t *tripFlowEventHandler) HandlerTripCreate(ctx context.Context, data []byte) {
	var payload messaging.CreateTripRequest
	if err := json.Unmarshal(data, &payload); err != nil {
		t.logger.Error("failed_to_unmarshal_create_trip_request", err, map[string]any{
			"data": data,
		})
		return
	}

	// ---- Goroutine 1: Create Trip & Update Rider Location ----
	go func(p messaging.CreateTripRequest) {
		// each goroutine gets its own context (prevent leaks)
		gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// trip create
		if err := t.tripFlowService.TripCreate(gCtx, p.RiderId); err != nil {
			t.logger.Error("TripCreate failed", err, map[string]any{
				"rider_id": p.RiderId,
			})
			return

		}

		// update rider location
		if err := t.userService.UpdateRiderLocation(gCtx, types.UpdateRiderLocationParam{
			RiderId:     p.ConnectionId,
			Location:    p.Pickup,
			Destination: p.Destination,
		}); err != nil {
			t.logger.Error("UpdateRiderLocation failed", err, map[string]any{
				"rider_id": p.RiderId,
			})
			return
		}
		t.logger.Info("trip_create_and_update_rider_location_successfully", map[string]any{
			"rider_id": p.RiderId,
		})
	}(payload)

	// ---- Goroutine 2: Routing via OSRM & Publish Event ----
	go func(p messaging.CreateTripRequest) {
		gCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		routeValue, err := t.osrmIntegration.GetRoutes(gCtx, p.Pickup, p.Destination)
		if err != nil {
			t.logger.Error("GetRoutes failed", err, map[string]any{
				"rider_id": p.RiderId,
			})
			return
		}

		routeResponse, err := json.Marshal(messaging.RoutesResponse{
			Coordinate: routeValue.Coordinate,
			Distance:   routeValue.Distance,
			Duration:   routeValue.Duration,
		})
		if err != nil {
			t.logger.Error("failed to marshal RoutesResponse", err, map[string]any{
				"rider_id": p.RiderId,
			})
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
			t.logger.Error("PublishingMessage failed", err, map[string]any{
				"rider_id": p.RiderId,
			})
			return
		}

		t.logger.Info("routing_via_osrm_and_publish_event_successfully", map[string]any{
			"rider_id": p.RiderId,
		})
	}(payload)
}

func NewTripFlowEventHandler(
	rabbitmq *messaging.RabbitMQ,
	tripFlowService domain.TripFlowService,
	userService domain.UserService,
	osrmIntegration domain.OsrmIntegrationService,
	logger logger.Logger,
) domain.TripFlowHandler {
	return &tripFlowEventHandler{
		rabbitmq:        rabbitmq,
		tripFlowService: tripFlowService,
		userService:     userService,
		osrmIntegration: osrmIntegration,
		logger:          logger,
	}
}

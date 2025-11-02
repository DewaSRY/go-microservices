package service

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/models"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tripService struct {
	Repo domain.TripRepository
}

func (t *tripService) GetUserTrip(ctx context.Context, userId string, fareId string) (*models.TripModel, error) {
	return t.Repo.FindTripWithFilter(func(d *gorm.DB) *gorm.DB {
		return d.Where(map[string]interface{}{
			"user_id":      userId,
			"ride_fare_id": fareId,
		})
	})
}

// GetTripProto implements domain.TripService.
func (t *tripService) GetTripProto(ctx context.Context, tripId string) (*tripgrpc.Trip, error) {
	currentTripModel, err := t.GetTripByID(ctx, tripId)

	trip := &tripgrpc.Trip{
		Id:     currentTripModel.ID,
		UserID: currentTripModel.UserID,
		Status: currentTripModel.Status,
	}

	if err != nil {
		return nil, fmt.Errorf("model_not_found:%v", err)
	}

	tripFareModel, _ := t.GetFareById(ctx, currentTripModel.RideFareID)

	if tripFareModel != nil {
		trip.SelectedFare = &tripgrpc.RideFare{
			Id:                tripFareModel.ID,
			PackageSlug:       tripFareModel.PackageSlug,
			UserID:            trip.UserID,
			TotalPriceInCents: tripFareModel.TotalPriceInCents,
		}
	}

	if tripFareModel != nil && tripFareModel.Routes != nil {
		var route types.Routes

		if err := json.Unmarshal(tripFareModel.Routes, &route); err != nil {
			return nil, err
		}
		trip.Route = route.ToRouteProto()
	}

	return trip, nil
}

// GetRoute implements domain.TripService.
func (t *tripService) GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	res, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("failed_to_parse:%v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("failed_to_read:%v", err)
	}

	var routeResponse types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResponse); err != nil {
		log.Print(err)
		return nil, fmt.Errorf("failed_to_unmarshal:%v", err)
	}

	return &routeResponse, nil
}

// GetTripByID implements domain.TripService.
func (t *tripService) GetTripByID(ctx context.Context, id string) (*models.TripModel, error) {
	return t.Repo.GetTripByID(ctx, id)
}

// GetUserRideFare implements domain.TripService.
func (t *tripService) GetUserRideFare(ctx context.Context, userID string, rideFareId string) (*models.FareModel, error) {
	currentFare, err := t.Repo.GetFareById(ctx, rideFareId)

	if err != nil {
		return nil, fmt.Errorf("fare_with_%s_not_found", rideFareId)
	}

	if currentFare.UserID != userID {
		return nil, fmt.Errorf("fare_and_user_id_not_match")
	}

	return currentFare, nil
}

// GetFareById implements domain.TripService.
func (t *tripService) GetFareById(ctx context.Context, fareId string) (*models.FareModel, error) {
	return t.Repo.GetFareById(ctx, fareId)
}

// CreateTrip implements domain.TripService.
func (t *tripService) CreateTrip(ctx context.Context, fare *models.FareModel) (*models.TripModel, error) {
	newTrip := &models.TripModel{
		ID:         uuid.New().String(),
		UserID:     fare.UserID,
		Status:     "pending",
		RideFareID: fare.ID,
	}

	if err := t.Repo.CreateTrip(ctx, newTrip); err != nil {
		return nil, fmt.Errorf("failed_to_create_ride_fare")
	}

	return newTrip, nil
}

// UpdateTrip implements domain.TripService.
func (t *tripService) UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error {
	return t.Repo.UpdateTrip(ctx, tripID, status, driver)
}

func NewTripService(repo domain.TripRepository) domain.TripService {
	return &tripService{Repo: repo}
}

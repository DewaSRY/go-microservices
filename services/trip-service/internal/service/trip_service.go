package service

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrFailedToGetRoute = errors.New("failed_to_get_route")
	ErrFailedToRead     = errors.New("failed_to_read")
	ErrFailedToParse    = errors.New("failed_to_parse")
)

type tripService struct {
	Repo domain.TripRepository
}

// GetUserRideFare implements domain.TripService.
func (t *tripService) GetUserRideFare(ctx context.Context, userID string, fareId string) (*domain.RideFareModel, error) {
	fare, err := t.Repo.GetRideFareById(ctx, fareId)
	if err != nil {
		return nil, err
	}

	if fare.UserID != userID {
		return nil, errors.New("ride_fare_not_found")
	}

	return fare, nil
}

// GetRoute implements domain.TripService.
func (t *tripService) GetRoute(ctx context.Context, pickup *types.Coordinate, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson", pickup.Longitude, pickup.Latitude, destination.Longitude, destination.Latitude)

	res, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return nil, ErrFailedToGetRoute
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
		return nil, ErrFailedToRead
	}

	var routeResponse types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResponse); err != nil {
		log.Print(err)
		return nil, ErrFailedToParse
	}

	return &routeResponse, nil
}

func (t *tripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	newTrip := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: *fare,
	}
	return t.Repo.CreateTrip(ctx, newTrip)
}

func NewTripService(repo domain.TripRepository) domain.TripService {
	return &tripService{Repo: repo}
}

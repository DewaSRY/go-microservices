package service

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"

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

// GenerateTripFares implements domain.TripService.
func (t *tripService) GenerateTripFares(ctx context.Context, fares []*triptype.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*triptype.RideFareModel, error) {
	fareList := make([]*triptype.RideFareModel, len(fares))

	for i, f := range fares {
		Id := primitive.NewObjectID()
		fare := triptype.RideFareModel{
			UserID:            userId,
			ID:                Id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
			Route:             &route.Routes[0],
		}

		fareList[i] = &fare

	}

	if err := t.Repo.SaveRIdeFareList(ctx, fareList); err != nil {
		return nil, fmt.Errorf("failed_to_save_trip_fare:%v", err)
	}

	return fareList, nil
}

// GetFareById implements domain.TripService.
func (t *tripService) GetFareById(ctx context.Context, fareId string) (*triptype.RideFareModel, error) {
	return t.Repo.GetFareById(ctx, fareId)
}

// GetTripByID implements domain.TripService.
func (t *tripService) GetTripByID(ctx context.Context, id string) (*triptype.TripModel, error) {
	return t.Repo.GetTripByID(ctx, id)
}

// UpdateTrip implements domain.TripService.
func (t *tripService) UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error {
	return t.Repo.UpdateTrip(ctx, tripID, status, driver)
}

// GetUserRideFare implements domain.TripService.
func (t *tripService) GetUserRideFare(ctx context.Context, userID string, fareId string) (*triptype.RideFareModel, error) {
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

func (t *tripService) CreateTrip(ctx context.Context, fare *triptype.RideFareModel) (*triptype.TripModel, error) {
	newTrip := &triptype.TripModel{
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

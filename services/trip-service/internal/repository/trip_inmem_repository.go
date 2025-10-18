package repository

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"context"
	"errors"
	"fmt"
)

var ErrTripIdAlreadyUser = errors.New("trip_with_id_already_use")

type inMemoryTripRepository struct {
	tripsMap  map[string]*triptype.TripModel
	rideFares map[string]*triptype.RideFareModel
}

// GetFareById implements domain.TripRepository.
func (i *inMemoryTripRepository) GetFareById(ctx context.Context, fareId string) (*triptype.RideFareModel, error) {
	fare, exits := i.rideFares[fareId]
	if !exits {
		return nil, fmt.Errorf("fare_with_id_%s_not_found", fareId)
	}

	return fare, nil
}

// GetTripByID implements domain.TripRepository.
func (i *inMemoryTripRepository) GetTripByID(ctx context.Context, id string) (*triptype.TripModel, error) {
	trip, exits := i.tripsMap[id]
	if !exits {
		return nil, fmt.Errorf("fare_with_id_%s_not_found", id)
	}

	return trip, nil
}

// UpdateTrip implements domain.TripRepository.
func (i *inMemoryTripRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error {
	trip, exits := i.tripsMap[tripID]
	if !exits {
		return fmt.Errorf("trip_with_id_%s_not_found", tripID)
	}
	trip.Status = status

	if driver != nil {
		trip.Driver = &triptype.TripDriver{
			ID:             driver.Id,
			Name:           driver.Name,
			ProfilePicture: driver.ProfilePicture,
			CarPlate:       driver.CarPlate,
		}
	}

	i.tripsMap[tripID] = trip
	return nil
}

// GetRideFareById implements domain.TripRepository.
func (i *inMemoryTripRepository) GetRideFareById(ctx context.Context, rideFareID string) (*triptype.RideFareModel, error) {
	fare, exist := i.rideFares[rideFareID]
	if exist {
		return fare, nil
	}

	return nil, errors.New("fare_not_found")
}

// SaveRIdeFareList implements domain.TripRepository.
func (i *inMemoryTripRepository) SaveRIdeFareList(ctx context.Context, fares []*triptype.RideFareModel) error {

	for _, fare := range fares {
		i.rideFares[fare.ID.Hex()] = fare
	}

	return nil
}

// CreateTrip implements domain.TripRepository.
func (i *inMemoryTripRepository) CreateTrip(ctx context.Context, trip *triptype.TripModel) (*triptype.TripModel, error) {
	_, exist := i.tripsMap[trip.ID.Hex()]

	if !exist {
		i.tripsMap[trip.ID.Hex()] = trip
		return trip, nil
	} else {
		return nil, ErrTripIdAlreadyUser
	}

}

func NewInMemoryTripRepository() domain.TripRepository {
	return &inMemoryTripRepository{
		tripsMap:  make(map[string]*triptype.TripModel),
		rideFares: make(map[string]*triptype.RideFareModel),
	}
}

package repository

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"
	"context"
	"errors"
)

var ErrTripIdAlreadyUser = errors.New("trip_with_id_already_use")

type inMemoryTripRepository struct {
	tripsMap  map[string]*triptype.TripModel
	rideFares map[string]*triptype.RideFareModel
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

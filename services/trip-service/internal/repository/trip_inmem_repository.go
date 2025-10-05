package repository

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"context"
	"errors"
)

var ErrTripIdAlreadyUser = errors.New("trip_with_id_already_use")

type inMemoryTripRepository struct {
	tripsMap  map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

// SaveRIdeFareList implements domain.TripRepository.
func (i *inMemoryTripRepository) SaveRIdeFareList(ctx context.Context, fares []*domain.RideFareModel) error {

	for _, fare := range fares {
		i.rideFares[fare.ID.Hex()] = fare
	}

	return nil
}

// CreateTrip implements domain.TripRepository.
func (i *inMemoryTripRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
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
		tripsMap:  make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

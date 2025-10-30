package repository

// import (
// 	"DewaSRY/go-microservices/services/trip-service/internal/domain"
// 	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"

// 	"DewaSRY/go-microservices/shared/models"
// 	"context"
// 	"errors"
// 	"fmt"
// )

// var ErrTripIdAlreadyUser = errors.New("trip_with_id_already_use")

// type inMemoryTripRepository struct {
// 	tripsMap  map[string]*models.TripModel
// 	rideFares map[string]*models.FareModel
// }

// // GetFareById implements domain.TripRepository.
// func (i *inMemoryTripRepository) GetFareById(ctx context.Context, fareId uint64) (*models.FareModel, error) {
// 	fare, exits := i.rideFares[fmt.Sprintf("{%d}", fareId)]
// 	if !exits {
// 		return nil, fmt.Errorf("fare_with_id_%s_not_found", fareId)
// 	}

// 	return fare, nil
// }

// // GetTripByID implements domain.TripRepository.
// func (i *inMemoryTripRepository) GetTripByID(ctx context.Context, id uint64) (*models.TripModel, error) {
// 	trip, exits := i.tripsMap[fmt.Sprintf("{%d}", id)]
// 	if !exits {
// 		return nil, fmt.Errorf("fare_with_id_%s_not_found", id)
// 	}

// 	return trip, nil
// }

// // UpdateTrip implements domain.TripRepository.
// func (i *inMemoryTripRepository) UpdateTrip(ctx context.Context, tripID uint64, status string, driver *drivergrpc.Driver) error {
// 	trip, exits := i.tripsMap[fmt.Sprintf("{%d}", tripID)]
// 	if !exits {
// 		return fmt.Errorf("trip_with_id_%s_not_found", tripID)
// 	}
// 	trip.Status = status

// 	// if driver != nil {
// 	// 	trip.Driver = &models.DriverModel{
// 	// 		ID:             0,
// 	// 		Name:           driver.Name,
// 	// 		ProfilePicture: driver.ProfilePicture,
// 	// 		CarPlate:       driver.CarPlate,
// 	// 	}
// 	// }

// 	i.tripsMap[fmt.Sprintf("{%d}", fareId)] = trip
// 	return nil
// }

// // GetRideFareById implements domain.TripRepository.
// func (i *inMemoryTripRepository) GetRideFareById(ctx context.Context, rideFareID uint64) (*models.FareModel, error) {
// 	fare, exist := i.rideFares[fmt.Sprintf("{%d}", fareId)]
// 	if exist {
// 		return fare, nil
// 	}

// 	return nil, errors.New("fare_not_found")
// }

// // SaveRIdeFareList implements domain.TripRepository.
// func (i *inMemoryTripRepository) SaveRIdeFareList(ctx context.Context, fares []*models.FareModel) error {

// 	for _, fare := range fares {
// 		hashId := fmt.Sprintf("{%d}", fare.ID)
// 		i.rideFares[hashId] = fare
// 	}

// 	return nil
// }

// // CreateTrip implements domain.TripRepository.
// func (i *inMemoryTripRepository) CreateTrip(ctx context.Context, trip *models.TripModel) (*models.TripModel, error) {
// 	hashTripID := fmt.Sprintf("{%d}", trip.ID)
// 	_, exist := i.tripsMap[hashTripID]

// 	if !exist {
// 		i.tripsMap[hashTripID] = trip
// 		return trip, nil
// 	} else {
// 		return nil, ErrTripIdAlreadyUser
// 	}

// }

// func NewInMemoryTripRepository() domain.TripRepository {
// 	return &inMemoryTripRepository{
// 		tripsMap:  make(map[string]*models.TripModel),
// 		rideFares: make(map[string]*models.TripModel),
// 	}
// }

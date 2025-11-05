package service

import (
	"DewaSRY/go-microservices/services/driver-service/internal/domain"
	driverUtil "DewaSRY/go-microservices/services/driver-service/internal/util"
	"DewaSRY/go-microservices/shared/models"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"DewaSRY/go-microservices/shared/types"
	"DewaSRY/go-microservices/shared/util"
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

type driverService struct {
	driverRepo domain.DriverRepository
	mu         sync.RWMutex
}

// GetDriverProto implements domain.DriverService.
func (t *driverService) GetDriverProto(ctx context.Context, driverId string) (*drivergrpc.Driver, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	currentModel, err := t.driverRepo.GetDriverById(ctx, driverId)

	if err != nil {
		return nil, fmt.Errorf("driver_service_failed_to_get_driver:%w", err)
	}

	var driverLocation drivergrpc.Location

	if err := json.Unmarshal(currentModel.Location, &driverLocation); err != nil {
		return nil, fmt.Errorf("driver_service_failed_to_get_driver:%w", err)
	}

	driverProto := &drivergrpc.Driver{
		Id:             currentModel.ID,
		Name:           currentModel.Name,
		ProfilePicture: currentModel.ProfilePicture,
		CarPlate:       currentModel.CarPlate,
		Geohash:        currentModel.Geohash,
		PackageSlug:    currentModel.PackageSlug,
		Location:       &driverLocation,
	}
	return driverProto, nil
}

// FindAvailableDrivers implements domain.DriverService.
func (t *driverService) FindAvailableDrivers(ctx context.Context, packageTypes string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	list, err := t.driverRepo.GetActiveDriverIdList(ctx, func(d *gorm.DB) *gorm.DB {
		return d.Where("package_slug = ?", packageTypes).Where("is_active", true)
	})

	if err != nil {
		return make([]string, 0)
	}

	return list
}

// RegisterDriver implements domain.DriverService.
func (t *driverService) RegisterDriver(ctx context.Context, driverId string, packageSlug string) (*models.DriverModel, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	randomIndex := rand.IntN(len(driverUtil.PredefinedRoutes))
	randomRoute := driverUtil.PredefinedRoutes[randomIndex]

	randomPlat := driverUtil.GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)
	geoHash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	jsonCoordinate, err := json.Marshal(types.Coordinate{
		Latitude:  randomRoute[0][0],
		Longitude: randomRoute[0][1],
	})

	if err != nil {
		return nil, fmt.Errorf("failed_to_create_driver:%w", err)
	}

	currentDriver := models.DriverModel{
		ID:             driverId,
		Name:           "temp",
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlat,
		Geohash:        geoHash,
		IsActive:       true,
		Location:       jsonCoordinate,
	}

	if err := t.driverRepo.CreateDriver(ctx, &currentDriver); err != nil {
		return nil, fmt.Errorf("failed_to_create_driver:%w", err)
	}

	return &currentDriver, nil
}

// UnregisterDriver implements domain.DriverService.
func (t *driverService) UnregisterDriver(ctx context.Context, driverId string) error {

	return t.driverRepo.UpdateDriverData(ctx, driverId,
		map[string]interface{}{
			"is_active": false,
		},
	)
}

func NewDriverService(driverRepo domain.DriverRepository) domain.DriverService {
	return &driverService{
		driverRepo: driverRepo,
	}
}

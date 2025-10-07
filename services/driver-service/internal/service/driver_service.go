package service

import (
	"DewaSRY/go-microservices/services/driver-service/internal/domain"
	driverUtil "DewaSRY/go-microservices/services/driver-service/internal/util"
	"DewaSRY/go-microservices/services/driver-service/pkg"
	"DewaSRY/go-microservices/shared/util"
	"math/rand/v2"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type driverService struct {
	driverList []*pkg.ActiveDriver
	mu         sync.RWMutex
}

func NewDriverService() domain.DriverService {
	return &driverService{
		driverList: make([]*pkg.ActiveDriver, 0),
	}
}

func (s *driverService) FindAvailableDrivers(packageTypes string) []string {
	var matchDrivers []string

	for _, driver := range s.driverList {
		if driver.Driver.PackageSlug == packageTypes {
			matchDrivers = append(matchDrivers, driver.Driver.ID)
		}
	}

	if len(matchDrivers) == 0 {
		return []string{}
	}

	return matchDrivers
}

func (s *driverService) RegisterDriver(driverId string, packageSlug string) (pkg.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := rand.IntN(len(driverUtil.PredefinedRoutes))
	randomRoute := driverUtil.PredefinedRoutes[randomIndex]

	randomPlat := driverUtil.GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)

	geoHash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := pkg.Driver{
		ID:             driverId,
		Geohash:        geoHash,
		Name:           "Land Nories",
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlage:       randomPlat,
		Location: &pkg.Location{
			Latitude:  randomRoute[0][0],
			Longitude: randomRoute[0][1],
		},
	}

	s.driverList = append(s.driverList, &pkg.ActiveDriver{
		Driver: driver,
	})

	return driver, nil
}

func (s *driverService) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, driver := range s.driverList {
		if driver.Driver.ID == driverId {
			s.driverList = append(s.driverList[:i], s.driverList[i+1:]...)
		}
	}
}

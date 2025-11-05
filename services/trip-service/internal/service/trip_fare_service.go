package service

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/models"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type tripFareService struct {
	repo domain.TripRepository
}

// EstimatePackagesPriceWithRoute implements domain.TripFareService.
func (t *tripFareService) EstimatePackagesPriceWithRoute(route *types.OsrmApiResponse) []*models.FareModel {
	baseFares := getBaseFares()
	fareList := make([]*models.FareModel, len(baseFares))

	for i, f := range baseFares {
		fareList[i] = estimationFareRoute(f, route)
	}

	return baseFares
}

// EstimatePackagesPrice implements domain.TripFareService.
func (t *tripFareService) EstimatePackagesPrice(distanceInKm float64, duration float64) []*models.FareModel {
	baseFareList := getBaseFares()
	priceConfig := domain.DefaultPricingConfig()
	estimateFareList := make([]*models.FareModel, len(baseFareList))

	for idx, fare := range baseFareList {
		estimateFareList[idx] = t.estimatePackagePice(fare, priceConfig, distanceInKm, duration)
	}

	return estimateFareList
}

// GenerateTripFares implements domain.TripFareService.
func (t *tripFareService) GenerateTripFares(ctx context.Context, fares []*models.FareModel, userId string, route *types.OsrmApiResponse) ([]*models.FareModel, error) {
	fareList := make([]*models.FareModel, len(fares))

	routesJSON, err := json.Marshal(&route.Routes[0])
	if err != nil {
		return nil, fmt.Errorf("failed_to_parse_route:%v", err)
	}

	for i, fare := range fares {
		createFare := &models.FareModel{
			ID:                uuid.New().String(),
			UserID:            userId,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
			Routes:            routesJSON,
		}
		fareList[i] = createFare
	}

	if err := t.repo.CreateFareList(ctx, fareList); err != nil {
		return nil, err
	}

	return fareList, nil
}

func (t *tripFareService) estimatePackagePice(fare *models.FareModel, priceConfig *domain.PricingConfig, distanceInKm float64, duration float64) *models.FareModel {
	carPackagePrice := fare.TotalPriceInCents

	distanceFare := distanceInKm * priceConfig.PricePerUnitOfDistance
	timeFare := duration * priceConfig.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare*timeFare

	return &models.FareModel{
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: totalPrice,
	}
}

func estimationFareRoute(f *models.FareModel, route *types.OsrmApiResponse) *models.FareModel {
	pricingCfg := domain.DefaultPricingConfig()

	carPackagePrice := f.TotalPriceInCents

	distanceKM := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceKM * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &models.FareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func getBaseFares() []*models.FareModel {
	return []*models.FareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}
func NewTripFareService(repo domain.TripRepository) domain.TripFareService {
	return &tripFareService{repo: repo}
}

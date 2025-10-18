package service

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"time"

	triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tripFareService struct {
	repo domain.TripRepository
}

// EstimatePackagesPriceWithRoute implements domain.TripFareService.
func (t *tripFareService) EstimatePackagesPriceWithRoute(route *types.OsrmApiResponse) []*triptype.RideFareModel {
	baseFares := getBaseFares()
	fareList := make([]*triptype.RideFareModel, len(baseFares))

	for i, f := range baseFares {
		fareList[i] = estimationFareRoute(f, route)
	}

	return baseFares
}

func (t *tripFareService) EstimatePackagesPrice(distanceInKm float64, duration float64) []*triptype.RideFareModel {
	baseFareList := getBaseFares()
	priceConfig := domain.DefaultPricingConfig()
	estimateFareList := make([]*triptype.RideFareModel, len(baseFareList))

	for idx, fare := range baseFareList {
		estimateFareList[idx] = t.estimatePackagePice(fare, priceConfig, distanceInKm, duration)
	}

	return estimateFareList
}

func (t *tripFareService) GenerateTripFares(ctx context.Context, fares []*triptype.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*triptype.RideFareModel, error) {
	fareList := make([]*triptype.RideFareModel, len(fares))

	for i, fare := range fares {
		id := primitive.NewObjectID()

		createFare := &triptype.RideFareModel{
			UserID:            userId,
			ID:                id,
			PackageSlug:       fare.PackageSlug,
			Route:             &route.Routes[0],
			TotalPriceInCents: fare.TotalPriceInCents,
			ExpiresAt:         time.Now(),
		}
		fareList[i] = createFare
	}

	if err := t.repo.SaveRIdeFareList(ctx, fareList); err != nil {
		return nil, err
	}

	return fareList, nil
}

func (t *tripFareService) estimatePackagePice(fare *triptype.RideFareModel, priceConfig *domain.PricingConfig, distanceInKm float64, duration float64) *triptype.RideFareModel {
	carPackagePrice := fare.TotalPriceInCents

	distanceFare := distanceInKm * priceConfig.PricePerUnitOfDistance
	timeFare := duration * priceConfig.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare*timeFare

	return &triptype.RideFareModel{
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: totalPrice,
	}
}

func estimationFareRoute(f *triptype.RideFareModel, route *types.OsrmApiResponse) *triptype.RideFareModel {
	pricingCfg := domain.DefaultPricingConfig()

	carPackagePrice := f.TotalPriceInCents

	distanceKM := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceKM * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &triptype.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func getBaseFares() []*triptype.RideFareModel {
	return []*triptype.RideFareModel{
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

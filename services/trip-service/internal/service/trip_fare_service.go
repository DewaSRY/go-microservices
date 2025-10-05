package service

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tripFareService struct {
	repo domain.TripRepository
}

func (t *tripFareService) EstimatePackagesPrice(distanceInKm float64, duration float64) []*domain.RideFareModel {
	baseFareList := getBaseFares()
	priceConfig := domain.DefaultPricingConfig()
	estimateFareList := make([]*domain.RideFareModel, len(baseFareList))

	for idx, fare := range baseFareList {
		estimateFareList[idx] = t.estimatePackagePice(fare, priceConfig, distanceInKm, duration)
	}

	return estimateFareList
}

func (t *tripFareService) GenerateTripFares(ctx context.Context, fares []*domain.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*domain.RideFareModel, error) {
	fareList := make([]*domain.RideFareModel, len(fares))

	for i, fare := range fares {
		id := primitive.NewObjectID()

		createFare := &domain.RideFareModel{
			UserID:            userId,
			ID:                id,
			PackageSlug:       fare.PackageSlug,
			Route:             route,
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

func (t *tripFareService) estimatePackagePice(fare *domain.RideFareModel, priceConfig *domain.PricingConfig, distanceInKm float64, duration float64) *domain.RideFareModel {
	carPackagePrice := fare.TotalPriceInCents

	distanceFare := distanceInKm * priceConfig.PricePerUnitOfDistance
	timeFare := duration * priceConfig.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare*timeFare

	return &domain.RideFareModel{
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: totalPrice,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
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

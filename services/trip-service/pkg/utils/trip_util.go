package utils

import triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"

func getBaseFares() []*triptype.FareType {
	return []*triptype.FareType{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200.0,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350.0,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400.0,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000.0,
		},
	}
}

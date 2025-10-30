package tests

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/services/trip-service/internal/repository"
	"DewaSRY/go-microservices/services/trip-service/internal/service"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"log"
	"testing"
)

func createService() domain.TripService {
	db, err := makeDb()
	if err != nil {
		log.Panicf("failed_to_make_db_connection:%v", err)
		return nil
	}
	repository := repository.NewTripRepository(db)

	return service.NewTripService(repository)
}

func TestCRUDTripService(t *testing.T) {
	db, err := makeDb()
	if err != nil {
		log.Panicf("failed_to_make_db_connection:%v", err)

	}
	ctx := context.Background()

	service := createService()
	if service == nil {
		log.Panic("failed_to_create_service")
	}

	tempUser := &models.UserModel{
		Name: "tempUser",
	}

	if result := db.DB.WithContext(context.Background()).Create(tempUser); result.Error != nil {
		t.Fatalf("failed_to_create_user: %v", err)
	}

	tempFare := &models.FareModel{
		UserID:            tempUser.ID,
		PackageSlug:       "car",
		TotalPriceInCents: 1300,
	}

	newTrip, err := service.CreateTrip(ctx, tempFare)

	if err != nil {
		t.Fatalf("failed_to_create_trip: %v", err)
	}

	log.Printf("result_is:%v", newTrip)

}

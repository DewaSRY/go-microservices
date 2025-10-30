package tests

import (
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"log"
)

// var postgres_uri = env.GetString("POSTGRES_URI", "localhost://postgres:postgres@postgres:5433/appdb?sslmode=disable")
// var postgres_uri = "postgres://postgres:postgres@localhost:5433/appdb?sslmode=disable"
var postgres_uri = "postgres://postgres:postgres@localhost:5434/mydb?sslmode=disable"

func makeDb() (*db.PostgresManager, error) {
	db, err := db.NewPostgresManager(context.Background(), postgres_uri)
	if err != nil {
		log.Panicf("failed_to_make_db_connection:%v", err)
		return nil, err
	}

	if err := db.DB.AutoMigrate(&models.TripModel{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	if err := db.DB.AutoMigrate(&models.DriverModel{}); err != nil {
		log.Fatal("migration failed:", err)
	}
	if err := db.DB.AutoMigrate(&models.FareModel{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	if err := db.DB.AutoMigrate(&models.TransactionModel{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	if err := db.DB.AutoMigrate(&models.UserModel{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	return db, err
}

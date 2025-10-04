package main

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/services/trip-service/internal/infrastructure/repository"
	"DewaSRY/go-microservices/services/trip-service/internal/repository"
	"DewaSRY/go-microservices/services/trip-service/internal/service"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/util"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	grpcserver "google.golang.org/grpc"
)

var (
	serverName = "trip_service"
	PORT       = env.GetString("PORT", "8083")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Printf("start_service:%s\n", serverName)

	tripRepo := repository.NewInMemoryTripRepository()
	tripService := service.NewTripService(tripRepo)
	// tripRepo := repository.NewInMemoryTripRepository()
	// tripService := service.NewTripService(tripRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("failed_to_listen:%v", err)

	}
	grpcService := grpcserver.NewServer()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /preview", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		urlQuery := r.URL.Query()

		UserID := urlQuery.Get("userId")
		PackageSlug := urlQuery.Get("packageSlug")

		if UserID == "" || PackageSlug == "" {
			util.WriteJSONResponse(w, http.StatusUnprocessableEntity, map[string]any{
				"message": "userId and packageSlug are required",
			})
			return
		}

		createdTrip, err := tripService.CreateTrip(ctx, &domain.RideFareModel{
			ID:                primitive.NewObjectID(),
			UserID:            UserID,
			PackageSlug:       PackageSlug,
			TotalPriceInCents: 18,
			ExpiresAt:         time.Now(),
		})

		if err != nil {
			errorResponse := make(map[string]any)
			errorResponse["messages"] = err.Error()
			util.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
			return
		}

		util.WriteJSONResponse(w, http.StatusCreated, createdTrip)
	})

	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", PORT),
	// 	Handler: mux,
	// }

	go func() {
		log.Printf("success_run_service:%s\n", serverName)
		// if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// 	log.Fatalf("Listen : %s\n", err)
		// }
		if err := grpcService.Serve(lis); err != nil {
			log.Fatalf("Listen : %s\n", err)
		}
		cancel()
	}()

	<-ctx.Done()
	{
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM)
		<-quit
		log.Println("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("server_forced_to_shutdown:%v", err)
		}

		log.Println("server_exiting_gracefully")
	}
}

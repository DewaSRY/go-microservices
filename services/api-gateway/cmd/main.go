package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "DewaSRY/go-microservices/services/api-gateway/internal/infrastructure/http"
	"DewaSRY/go-microservices/shared/env"
)

var (
	serviceName = "API_GATEWAY"
	PORT        = env.GetString("PORT", "8081")
)

func main() {
	log.Printf("starting_app:%s\n", serviceName)
	// INIT
	mux := http.NewServeMux()
	handler := httpHandler.NewHttpHandler()

	//REGISTER HANDLER
	mux.HandleFunc("GET /", handler.GetHealthCheck)
	mux.HandleFunc("POST /trip/preview", handler.PostTripPreview)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed_to_run:%s", serviceName)
		}
	}()

	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGINT, syscall.SIGALRM)
	<-quite
	log.Println("shout down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("failed_to_shout_down")
		server.Close()
	}

	log.Println("gracefully_Shout_down")

}

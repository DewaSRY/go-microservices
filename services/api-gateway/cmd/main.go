package main

import (
	"fmt"
	"log"
	"net/http"

	httpHandler "DewaSRY/go-microservices/services/api-gateway/internal/infrastructure/http"
	"DewaSRY/go-microservices/shared/env"
)

var (
	serviceName = "API_GATEWAY"
	httpAddr    = env.GetString("HTTP_ADDR", "8081")
)

func main() {
	log.Printf("starting_app:%s\n", serviceName)
	// INIT
	mux := http.NewServeMux()
	handler := httpHandler.NewHttpHandler()

	//REGISTER HANDLER
	mux.HandleFunc("GET /", handler.GetHealthCheck)
	mux.HandleFunc("POST /trip/preview", handler.PostTripPreview)

	//RUN SERVER
	log.Printf("success_run:%s", serviceName)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", httpAddr), mux); err != nil {
		log.Printf("%s", fmt.Sprintf("failed_to_run:%s\n", serviceName))
	}

}

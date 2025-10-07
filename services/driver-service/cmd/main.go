package main

import (
	"DewaSRY/go-microservices/services/driver-service/internal/handler"
	"DewaSRY/go-microservices/services/driver-service/internal/service"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

var GrpcAddr = ":9092"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCn := make(chan os.Signal, 1)
		signal.Notify(sigCn, os.Interrupt, syscall.SIGTERM)
		<-sigCn
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)

	if err != nil {
		log.Fatalf("failed_to_listen: %v", err)
	}

	svc := service.NewDriverService()

	grpcServer := grpc.NewServer()
	handler.NewGrpcHandler(grpcServer, svc)

	log.Printf("starting_grpc_server_driver_service_on_port:%s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed_to_serve:%s", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutting_down_the_server..")
	grpcServer.GracefulStop()
}

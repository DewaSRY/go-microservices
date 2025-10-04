package grpcclient

import (
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"os"

	"google.golang.org/grpc"
)

type TripServiceClient struct {
	Client tripgrpc.TripServiceClient
	Conn   *grpc.ClientConn
}

func NewTripServiceClient() (*TripServiceClient, error) {
	tripServiceUrl := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceUrl == "" {
		tripServiceUrl = "trip-service:9093"
	}
	coon, err := grpc.NewClient(tripServiceUrl)

	if err != nil {
		return nil, err
	}

	Client := tripgrpc.NewTripServiceClient(coon)
	return &TripServiceClient{
		Conn:   coon,
		Client: Client,
	}, nil
}

func (c *TripServiceClient) Close() {
	if c.Conn != nil {
		if err := c.Conn.Close(); err != nil {
			return
		}
	}
}

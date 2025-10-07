package handler

import (
	"DewaSRY/go-microservices/services/driver-service/internal/domain"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type driverGrpcHandler struct {
	drivergrpc.UnimplementedDriverServiceServer
	service domain.DriverService
}

func NewGrpcHandler(server *grpc.Server, service domain.DriverService) *driverGrpcHandler {
	handler := &driverGrpcHandler{
		service: service,
	}

	drivergrpc.RegisterDriverServiceServer(server, handler)
	return handler
}

func (h *driverGrpcHandler) RegisterDriver(ctx context.Context, req *drivergrpc.RegisterDriverRequest) (*drivergrpc.RegisterDriverResponse, error) {
	driver, err := h.service.RegisterDriver(req.GetDriverId(), req.GetPackageSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed_to_register_driver:%v", err)
	}

	return &drivergrpc.RegisterDriverResponse{
		Driver: driver.ToDriverProto(),
	}, nil
}

func (h *driverGrpcHandler) UnregisterDriver(ctx context.Context, req *drivergrpc.RegisterDriverRequest) (*drivergrpc.RegisterDriverResponse, error) {
	h.service.UnregisterDriver(req.GetDriverId())
	return &drivergrpc.RegisterDriverResponse{
		Driver: &drivergrpc.Driver{
			Id: req.GetDriverId(),
		},
	}, nil
}

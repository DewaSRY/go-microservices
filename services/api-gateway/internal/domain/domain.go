package domain

import (
	"context"
	"net/http"
)

type HttpHandler interface {
	GetHealthCheck(w http.ResponseWriter, r *http.Request)
	PostTripPreview(w http.ResponseWriter, r *http.Request)
	PostStartTrip(w http.ResponseWriter, r *http.Request)
}

type RideShareServices interface {
	UserInitEventRequest(ctx context.Context, connectionId string, data []byte) error
	CreateTripEvent(ctx context.Context, connectionId string, data []byte) error
	RiderCreateTripRequest(ctx context.Context, connectionId string, data []byte) error
	DriverInitRequest(ctx context.Context, connectionId string, data []byte) error
	RiderCreateTransaction(ctx context.Context, connectionId string, data []byte) error
}

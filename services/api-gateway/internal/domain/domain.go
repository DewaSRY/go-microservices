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
	UserInitEvent(ctx context.Context, connectionId string, data []byte)
}

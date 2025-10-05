package domain

import "net/http"

type HttpHandler interface {
	GetHealthCheck(w http.ResponseWriter, r *http.Request)
	PostTripPreview(w http.ResponseWriter, r *http.Request)
	PostStartTrip(w http.ResponseWriter, r *http.Request)
}

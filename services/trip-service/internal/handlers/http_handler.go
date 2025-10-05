package handlers

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/dto"
	"DewaSRY/go-microservices/shared/util"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HttpHandler interface {
	GetPreview(w http.ResponseWriter, r *http.Request)
	GetTripPreview(w http.ResponseWriter, r *http.Request)
}

type httpHandler struct {
	service domain.TripService
}

func (h *httpHandler) GetTripPreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqBody dto.PreviewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "failed_to_parse_json",
			"data":    err,
		})
		return
	}

	route, err := h.service.GetRoute(ctx, &reqBody.Pickup, &reqBody.Destination)

	if err != nil {
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "failed_to_get_the_route",
			"data":    err,
		})
	}

	util.WriteJSONResponse(w, http.StatusOK, route)

}

// TODO: this is i am not really sure about this
func (h *httpHandler) GetPreview(w http.ResponseWriter, r *http.Request) {
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

	createdTrip, err := h.service.CreateTrip(ctx, &domain.RideFareModel{
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
}

func NewHttpHandler(service domain.TripService) HttpHandler {
	return &httpHandler{service: service}
}

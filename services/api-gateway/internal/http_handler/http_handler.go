package http

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/domain"
	"DewaSRY/go-microservices/services/api-gateway/internal/dto"
	grpcclient "DewaSRY/go-microservices/services/api-gateway/internal/grpc_client"
	"DewaSRY/go-microservices/shared/contracts"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/util"
	"encoding/json"
	"net/http"
)

type httpHandler struct {
}

// PostStartTrip implements domain.HttpHandler.
func (h *httpHandler) PostStartTrip(w http.ResponseWriter, r *http.Request) {
	var reqBody dto.StartTripRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_parse_JSON_request"
		errorResponse["data"] = err
		util.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := util.ValidateStruct(reqBody); err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "validation_error"
		errorResponse["data"] = err.Error()
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, errorResponse)
		return
	}

	tripService, err := grpcclient.NewTripServiceClient()
	if err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_make_connection"
		errorResponse["data"] = err
		util.WriteJSONResponse(w, http.StatusInternalServerError, errorResponse)
		return
	}
	defer tripService.Close()

	tripCreated, err := tripService.Client.CreateTrip(ctx, &tripgrpc.CreateTripRequest{
		UserID:     reqBody.UserID,
		RideFareID: reqBody.RideFareID,
	})
	if err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_create_the_trip"
		errorResponse["data"] = err
		util.WriteJSONResponse(w, http.StatusInternalServerError, errorResponse)
		return
	}

	response := contracts.APIResponse{
		Data: tripCreated,
	}
	util.WriteJSONResponse(w, http.StatusCreated, response)

}

// PostTripPreview implements domain.HttpHandler.
func (h *httpHandler) PostTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody dto.PreviewTripRequest
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_parse_JSON_request"
		errorResponse["data"] = err
		util.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := util.ValidateStruct(reqBody); err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "validation_error"
		errorResponse["data"] = err.Error()
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, errorResponse)
		return
	}

	tripService, err := grpcclient.NewTripServiceClient()
	if err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_make_connection"
		errorResponse["data"] = err
		util.WriteJSONResponse(w, http.StatusInternalServerError, errorResponse)
		return
	}
	defer tripService.Close()

	tripResponse, err := tripService.Client.PreviewTrip(ctx, &tripgrpc.PreviewTripRequest{
		UserID: reqBody.UserID,
		StartLocation: &tripgrpc.Coordinate{
			Latitude:  reqBody.Pickup.Latitude,
			Longitude: reqBody.Pickup.Longitude,
		},
		EndLocation: &tripgrpc.Coordinate{
			Latitude:  reqBody.Destination.Latitude,
			Longitude: reqBody.Destination.Longitude,
		},
	})

	if err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_fetch_preview_trip"
		errorResponse["data"] = err.Error()
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, errorResponse)
		return
	}

	response := contracts.APIResponse{
		Data: tripResponse,
	}
	util.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *httpHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]any)
	response["message"] = "server_healthy"
	util.WriteJSONResponse(w, http.StatusOK, response)
}

func NewHttpHandler() domain.HttpHandler {
	return &httpHandler{}
}

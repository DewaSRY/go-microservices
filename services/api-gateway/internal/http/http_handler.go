package http

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/domain"
	"DewaSRY/go-microservices/services/api-gateway/internal/dto"
	"DewaSRY/go-microservices/shared/types"
	"DewaSRY/go-microservices/shared/util"
	"bytes"
	"encoding/json"
	"net/http"
)

type httpHandler struct {
}

// PostTripPreview implements domain.HttpHandler.
func (h *httpHandler) PostTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody dto.PreviewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_parse_JSON_request"
		util.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := util.ValidateStruct(reqBody); err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = err.Error()
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, errorResponse)
		return
	}

	jsonData, err := json.Marshal(reqBody)

	if err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = err.Error()
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, errorResponse)
		return
	}

	url := "http://localhost:8083/preview"
	// url := "http://trip-service:8083/preview"
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = err.Error()
		util.WriteJSONResponse(w, http.StatusUnprocessableEntity, errorResponse)
		return
	}

	var tripResponse types.OsrmApiResponse
	if err := json.NewDecoder(res.Body).Decode(&tripResponse); err != nil {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "failed_to_parse_JSON_request"
		util.WriteJSONResponse(w, http.StatusBadRequest, errorResponse)
		return
	}

	// response := make(map[string]any)
	// response["message"] = "success"
	util.WriteJSONResponse(w, http.StatusConflict, tripResponse)

}

func (h *httpHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]any)
	response["message"] = "server_healthy"
	util.WriteJSONResponse(w, http.StatusOK, response)
}

func NewHttpHandler() domain.HttpHandler {
	return &httpHandler{}
}

package messaging

type CommonSuccessResponse struct {
	Message string `json:"message"`
	Code    uint32 `json:"code"`
}

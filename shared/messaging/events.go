package messaging

import (
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
)

const (
	FindAvailableDriversQueue        = "find_available_drivers"
	DriverCmdTripRequestQueue        = "driver_cmd_trip_request"
	DriverTripResponseQueue          = "driver_trip_response"
	NotifyMatchingTripQueue          = "notify_match_trip"
	NotifyDriverNoDriversFoundQueue  = "notify_driver_no_drivers_found"
	NotifyDriverAssignQueue          = "notify_driver_assign"
	PaymentTripResponseQueue         = "payment_trip_response"
	NotifyPaymentSessionCreatedQueue = "notify_payment_session_created"
	NotifyPaymentSuccessQueue        = "payment_success"
	DeadLetterQueue                  = "dead_letter_queue"
)

type TripEventData struct {
	Trip *tripgrpc.Trip `json:"trip"`
}

type DriverTripResponseData struct {
	Driver  *drivergrpc.Driver `json:"driver"`
	TripID  string             `json:"tripID"`
	RiderId string             `json:"riderID"`
}

type AcceptedPayment struct {
	TripID  string `json:"tripID"`
	RiderID string `json:"riderID"`
}

//  tripID: trip?.tripID ?? "",
//   riderID: userID ?? "",

type PaymentEventSessionCreatedData struct {
	TripID    string  `json:"tripID"`
	SessionID string  `json:"sessionID"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type PaymentTripResponseData struct {
	TripID   string  `json:"tripID"`
	UserID   string  `json:"userID"`
	DriverID string  `json:"driverID"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type PaymentStatusUpdateData struct {
	TripID   string `json:"tripID"`
	UserID   string `json:"userID"`
	DriverID string `json:"driverID"`
}

// FindAvailableDriversQueue
type TripDriverFindData struct {
	AmountDriver int `json:"amountDriver"`
}

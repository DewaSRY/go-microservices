package contracts

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

type MessageData struct {
	ConnectionId string `json:"ConnectionId"`
	Data         []byte `json:"data"`
}

// Routing keys - using consistent event/command patterns
const (
	// Trip events (trip.event.*)
	TripEventCreated             = "trip.event.created"
	TripEventDriverAssigned      = "trip.event.driver_assigned"
	TripEventNoDriversFound      = "trip.event.no_drivers_found"
	TripEventDriversFound        = "trip.event.drivers_found"
	TripEventDriverNotInterested = "trip.event.driver_not_interested"

	// Driver commands (driver.cmd.*)
	DriverCmdTripRequest = "driver.cmd.trip_request"
	DriverCmdTripAccept  = "driver.cmd.trip_accept"
	DriverCmdTripDecline = "driver.cmd.trip_decline"
	DriverCmdLocation    = "driver.cmd.location"
	DriverCmdRegister    = "driver.cmd.register"

	// Payment events (payment.event.*)
	PaymentEventSessionCreated = "payment.event.session_created"
	PaymentEventSuccess        = "payment.event.success"
	PaymentEventFailed         = "payment.event.failed"
	PaymentEventCancelled      = "payment.event.cancelled"
	PaymentEventComplete       = "payment.event.complete"

	// Payment commands (payment.cmd.*)
	PaymentCmdCreateSession = "payment.cmd.create_session"

	//#################################
	// user
	UserInitEventProcess         = "user.init.process"
	UserInitSuccessResponse      = "user.init-success.response"
	UserCloseConnectiondataEvent = "user.event.disconnect"

	UserDisconnectedProcess = "user.disconnected.process"

	//
	RiderEventCreateTrip = "rider.event.create-trip"
	RiderUpdateTrip      = "rider.event.update-trip"

	//rider
	RiderCreateTripRequest = "rider.create-trip.request"
	RiderCreateTripProcess = "rider.create-trip.process"

	//trip flow
	TripCreateInitEvent = "trip_flow.create.request"

	TripCreateInitProcess  = "trip_flow.create.process"
	TripCreateSuccessEvent = "trip_flow.create-success.response"

	RouteFoundEvent = "trip.route-found.response"

	//
	DriverInitEvent        = "driver.init.request"
	DriverInitEventProcess = "driver.init.process"

	DriverActiveResponse = "driver.active-notify.response"

	//transaction
	RiderCreateTransactionRequest = "rider.create-transaction.request"
	RiderCreateTransactionProcess = "rider.create-transaction.process"

	RiderCreateTransactionResponse = "rider.create-transaction.response"

	DriverAcceptTransactionRequest = "driver.transaction.accepted.request"
	DriverAcceptTransactionProcess = "driver.transaction.accepted.process"

	TransactionAcceptedResponse = "transaction.accepted.response"
)

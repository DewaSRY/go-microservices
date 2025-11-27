package handler

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/domain"
	grpcclient "DewaSRY/go-microservices/services/api-gateway/internal/grpc_client"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/lib"
	"DewaSRY/go-microservices/shared/messaging"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"DewaSRY/go-microservices/shared/util"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Driver struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	CarPlate       string `json:"carPlage"`
	PackageSlug    string `json:"packageSlug"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type wsMessageData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type WsHandler struct {
	connManager lib.ConnectionManager
	rabbitMq    *messaging.RabbitMQ
	tripService domain.RideShareServices
}

type AcceptedPayment struct {
	TripID  string `json:"tripID"`
	RiderID string `json:"riderID"`
}

func NewWsHandler(connManager lib.ConnectionManager, rabbitMq *messaging.RabbitMQ, tripService domain.RideShareServices) *WsHandler {
	return &WsHandler{connManager: connManager, rabbitMq: rabbitMq, tripService: tripService}
}

func (t *WsHandler) WsHandleStartConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := t.connManager.InitUpgrade(w, r)
	if err != nil {
		log.Printf("failed_to_start_connection")
	}
	defer conn.Close()
	ctx := r.Context()
	connectionId := uuid.New().String()

	t.connManager.Add(connectionId, conn)
	defer t.connManager.Remove(connectionId)

	t.makeQueueConsumer(
		[]string{
			messaging.UserEstablishConnectionNotificationQueue,
			messaging.TripFlowNotificationQueue,
		},
	)

	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			log.Printf("failed_to_read_rider_connection: %v", err)
			break
		}

		var messageData wsMessageData
		if err := json.Unmarshal(p, &messageData); err != nil {
			log.Printf("error_unmarshaling_rider_message: %v", err)
			continue
		}

		log.Printf("connection_if:%s_event:%s", connectionId, messageData.Type)

		switch messageData.Type {
		case contracts.UserInitEventRequest:
			t.tripService.UserInitEventRequest(ctx, connectionId, messageData.Data)
		case contracts.TripCreateInitEvent:
			t.tripService.CreateTripEvent(ctx, connectionId, messageData.Data)
		case contracts.RiderCreateTripRequest:
			t.tripService.RiderCreateTripRequest(ctx, connectionId, messageData.Data)
		default:
			log.Printf("trip_received_unknown_messages: %v", messageData.Type)
		}
	}
}

func (t *WsHandler) WsHandleRider(w http.ResponseWriter, r *http.Request) {
	conn, err := t.connManager.InitUpgrade(w, r)
	if err != nil {
		log.Printf("failed_start_connection:%v", err)
		return
	}
	defer conn.Close()
	ctx := r.Context()
	urQuery := r.URL.Query()
	userID := urQuery.Get("userID")

	if userID == "" {
		log.Println("userId_is_required")
		return
	}

	t.connManager.Add(userID, conn)
	defer t.connManager.Remove(userID)

	t.makeQueueConsumer(
		[]string{
			messaging.NotifyDriverNoDriversFoundQueue,
			messaging.NotifyDriverAssignQueue,
			messaging.NotifyPaymentSessionCreatedQueue, // temp
			messaging.NotifyPaymentSuccessQueue,
			messaging.NotifyMatchingTripQueue,
		},
	)

	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			log.Printf("failed_to_read_rider_connection:%v", err)
			break
		}

		var tripMessage wsMessageData
		if err := json.Unmarshal(p, &tripMessage); err != nil {
			log.Printf("error_unmarshaling_rider_message: %v", err)
			continue
		}

		switch tripMessage.Type {
		case contracts.DriverCmdLocation:
			var coordinateData Coordinate
			if err := json.Unmarshal(p, &coordinateData); err != nil {
				log.Printf("error_unmarshaling_rider_message: %v", err)
				continue
			}
			log.Printf("trip_received_coordinate_data: %v", coordinateData)
		case contracts.PaymentEventSuccess:

			var payload messaging.AcceptedPayment
			if err := json.Unmarshal(tripMessage.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			log.Print("this is the data of PaymentEventSuccess", payload)

			if err := t.rabbitMq.PublishingMessage(
				ctx,
				contracts.PaymentEventSuccess,
				contracts.AmqpMessage{
					OwnerID: userID,
					Data:    tripMessage.Data,
				}); err != nil {
				log.Printf("Error_publishing_message_to_rabbitMQ: %v", err)
			}
		default:
			log.Printf("trip_received_unknown_messages: %v", tripMessage)
		}

	}
}

func (t *WsHandler) WsHandleDriver(w http.ResponseWriter, r *http.Request) {
	conn, err := t.connManager.InitUpgrade(w, r)

	if err != nil {
		log.Printf("failed_start_connection:%v", err)
		return
	}
	defer conn.Close()

	//Get data form query
	urQuery := r.URL.Query()
	userID := urQuery.Get("userID")
	packagesSlug := urQuery.Get("packageSlug")
	if userID == "" || packagesSlug == "" {
		log.Println("userId_and_packageSlug_is_required")
		return
	}
	t.connManager.Add(userID, conn)

	// Call the connection's writeMessage and read message method to send
	driverData := Driver{
		ID:             userID,
		Name:           "Tiago",
		ProfilePicture: util.GetRandomAvatar(1),
		CarPlate:       "hallo",
		PackageSlug:    packagesSlug,
	}
	msg := contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
		Data: driverData,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("failed_to_register_driver;%v", err)
		return
	}

	defer conn.Close()
	ctx := r.Context()

	driverService, err := grpcclient.NewDriverServiceClient()
	if err != nil {
		log.Fatalf("failed_to_make_grpc_connection:%v", err)
	}

	driverPayload := &drivergrpc.RegisterDriverRequest{
		DriverId:    userID,
		PackageSlug: packagesSlug,
	}

	defer func() {
		t.connManager.Remove(userID)
		driverService.Client.UnregisterDriver(ctx, driverPayload)
		driverService.Close()
		log.Println("driver_unregister: ", userID)
	}()

	registerDriver, err := driverService.Client.RegisterDriver(ctx, driverPayload)

	if err != nil {
		log.Printf("error_register_driver: %v", err)
		return
	}

	if err := t.connManager.Emit(userID, contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
		Data: registerDriver.Driver,
	}); err != nil {
		log.Printf("error_sending_message: %v", err)
		return
	}

	t.makeQueueConsumer(
		[]string{
			messaging.DriverCmdTripRequestQueue,
			messaging.NotifyPaymentSessionCreatedQueue, // temp
			messaging.NotifyPaymentSuccessQueue,
		},
	)

	for {
		_, p, err := conn.ReadMessage()

		if err != nil {
			log.Printf("failed_to_read_driver_connection:%v", err)
			return
		}

		var driverMsg wsMessageData
		if err := json.Unmarshal(p, &driverMsg); err != nil {
			log.Printf("error_unmarshaling_driver_message: %v", err)
			continue
		}

		switch driverMsg.Type {
		case contracts.DriverCmdLocation:
			continue
		case contracts.DriverCmdTripAccept:
			if err := t.rabbitMq.PublishingMessage(
				ctx,
				contracts.DriverCmdTripAccept,
				contracts.AmqpMessage{
					OwnerID: userID,
					Data:    driverMsg.Data,
				}); err != nil {
				log.Printf("Error_publishing_message_to_rabbitMQ: %v", err)
			}
		case contracts.DriverCmdTripDecline:
			if err := t.rabbitMq.PublishingMessage(
				ctx,
				contracts.DriverCmdTripDecline,
				contracts.AmqpMessage{
					OwnerID: userID,
					Data:    driverMsg.Data,
				}); err != nil {
				log.Printf("Error_publishing_message_to_rabbitMQ: %v", err)
			}

		default:
			log.Printf("Unknown_message_type: %s", driverMsg.Type)
		}

	}
}

func (t *WsHandler) makeQueueConsumer(queueList []string) {
	for _, q := range queueList {
		consumer := lib.NewQueueConsumer(*t.rabbitMq, t.connManager, q)
		if err := consumer.Start(); err != nil {
			log.Printf("failed_to_start_consumer_for_queue: %s: err: %v", q, err)
		}
	}
}

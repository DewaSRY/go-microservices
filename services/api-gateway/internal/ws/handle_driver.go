package ws

import (
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/util"
	"log"
	"math/rand"
	"net/http"
)

func WsHandleDriver(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
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

	// Call the connection's writeMessage and read message method to send

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Printf("failed_to_read_connection:%v", err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("failed_to_write_message:%v", err)
			return
		}

		driverRegisterData := Driver{
			ID:             "driver-1",
			Name:           "tom",
			ProfilePicture: util.GetRandomAvatar(rand.Intn(100)),
			CarPlate:       "123 ab",
			PackageSlug:    "some-slug",
		}

		msg := contracts.WSMessage{
			Type: "driver.cmd.register",
			Data: driverRegisterData,
		}

		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("error_sending_message:%v", err)
		}
	}
}

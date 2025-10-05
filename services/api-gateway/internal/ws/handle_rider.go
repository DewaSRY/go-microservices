package ws

import (
	"log"
	"net/http"
)

func WsHandleRider(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed_start_connection:%v", err)
		return
	}
	defer conn.Close()

	//Get data form query

	urQuery := r.URL.Query()
	userID := urQuery.Get("userID")

	if userID == "" {
		log.Println("userId_is_required")
		return
	}

	// Call the connection's writeMessage and read message method to send

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Printf("failed_to_read_connection:%v", err)
			break
		}

		log.Println("message", messageType, p)
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Printf("failed_to_write_message:%v", err)
		// 	return
		// }
	}
}

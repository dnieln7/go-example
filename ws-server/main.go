package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/server", serverHanlder)

	log.Println("Starting server...")
	err := http.ListenAndServe(":4444", nil)

	if err != nil {
		log.Fatal("Could not start server", err)
	}
}

func serverHanlder(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		log.Fatal("Upgrade error: ", err)
	}

	defer connection.Close()

	for {
		messageType, message, err := connection.ReadMessage()

		if err != nil {
			log.Println("Error reading message: ", err)
			continue
		}

		log.Println("Message received: ", message, " with type: ", messageType)

		// Reverse
		messageLen := len(message)

		for i := 0; i < messageLen/2; i++ {
			message[i], message[messageLen-1-i] = message[messageLen-1-i], message[i]
		}

		err = connection.WriteMessage(messageType, message)

		if err != nil {
			log.Println("Error writing message:", err)
		} else {
			log.Println("Message sent: ", message)
		}
	}
}

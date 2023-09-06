package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var connections []*websocket.Conn

func main() {
	incomingMessages := make(chan []byte)
	incomingConnections := make(chan *websocket.Conn)

	go func() {
		for {
			select {
			case newMessage := <-incomingMessages:
				broadcastMessage(newMessage, connections)
			case newConnection := <-incomingConnections:
				connections = append(connections, newConnection)
			}

		}
	}()

	http.HandleFunc("/server", func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)

		if err != nil {
			log.Fatal("Upgrade error: ", err)
		}

		incomingConnections <- connection

		for {
			messageType, message, err := connection.ReadMessage()

			if err != nil {
				log.Println("Error reading message: ", err)
				continue
			}

			connection.WriteMessage()

			log.Println("Message received: ", message, " with type: ", messageType)

			incomingMessages <- message
		}
	})

	// http.HandleFunc("/server", serverHanlder)

	log.Println("Starting server...")
	err := http.ListenAndServe(":4444", nil)

	if err != nil {
		log.Fatal("Could not start server", err)
	}
}

func broadcastMessage(message []byte, connections []*websocket.Conn) {
	// Reverse
	messageLen := len(message)

	for i := 0; i < messageLen/2; i++ {
		message[i], message[messageLen-1-i] = message[messageLen-1-i], message[i]
	}

	for _, connection := range connections {
		err := connection.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			log.Println("Error writing message:", err)
		} else {
			log.Println("Message sent: ", message)
		}
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

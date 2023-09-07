package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var connections []*websocket.Conn

func main() {
	incomingMessages := make(chan []byte)
	connectionUpdates := make(chan ConnectionUpdate)

	go func() {
		for {
			select {
			case newMessage := <-incomingMessages:
				broadcastMessage(newMessage, connections)
			case connectionUpdate := <-connectionUpdates:
				if connectionUpdate.append {
					connections = append(connections, connectionUpdate.connection)
				} else {
					removeConnection(connectionUpdate.connection)
				}
			}
		}
	}()

	http.HandleFunc("/server", broadcastServerHanlder(connectionUpdates, incomingMessages))

	// http.HandleFunc("/server", serverHanlder)

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
			closeErr, ok := err.(*websocket.CloseError)

			if ok {
				log.Println("Close frame received, clossing...", closeErr)
				break
			} else {
				log.Println("Error reading message: ", err)
				continue
			}
		}

		messageText := fmt.Sprintf("%s", message)

		log.Println("Message received: ", messageText, " with type: ", messageType)

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

func broadcastServerHanlder(connectionUpdates chan ConnectionUpdate, incomingMessages chan []byte) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request,) {
		connection, err := upgrader.Upgrade(writer, request, nil)

		if err != nil {
			log.Fatal("Upgrade error: ", err)
		}
	
		connectionUpdates <- ConnectionUpdate{
			connection: connection,
			append:     true,
		}
	
		for {
			messageType, message, err := connection.ReadMessage()
	
			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)
	
				if ok {
					log.Println("Close frame received, clossing...", closeErr)
	
					connectionUpdates <- ConnectionUpdate{
						connection: connection,
						append:     false,
					}
	
					break
				} else {
					log.Println("Error reading message: ", err)
					continue
				}
			}
	
			messageText := fmt.Sprintf("%s", message)
	
			log.Println("Message received: ", messageText, " with type: ", messageType)
	
			incomingMessages <- message
		}
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
			messageText := fmt.Sprintf("%s", message)

			log.Println("Message sent: ", messageText)
		}
	}
}

func removeConnection(connection *websocket.Conn) {
	var index = -1
	var last = len(connections) - 1

	if last == 0 {
		log.Println("Cleaning connections... ")

		connections = []*websocket.Conn{}
	} else {
		for i, conn := range connections {
			if conn == connection {
				index = i
				break
			}
		}

		log.Println("Removing connection... ", connection.LocalAddr(), " at index: ", index)

		if index != -1 {
			if index != last {
				connections[index] = connections[last]
			}

			connections = connections[:last]
		}

		log.Println("Connection removed, remaining connections: ", len(connections))
	}

	connection.Close()
}

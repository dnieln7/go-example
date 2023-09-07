package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	url := url.URL{Scheme: "ws", Host: ":4444", Path: "/server"}

	// Establish connection
	connection, response, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	log.Println("Connection establised with response: ", response)

	defer connection.Close()

	// Receive messages
	go func() {
		for {
			messageType, message, err := connection.ReadMessage()

			if err != nil {
				log.Println("Error reading message:", err)
			} else {
				messageText := fmt.Sprintf("%s", message)

				log.Println("Message:", messageText, " with type: ", messageType)
			}

		}
	}()

	// Read from stdin and send through websocket
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		err = connection.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))

		if err != nil {
			log.Println("Error writing message:", err)
		} else {
			log.Println("Message sent: ", scanner.Text())
		}
	}
}

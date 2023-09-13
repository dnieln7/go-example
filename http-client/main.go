package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Request struct {
	ChatID  string `json:"chat_id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func main() {
	const chatID string = "403dabbe-1dbb-4d99-99ae-683200bd53d5"

	requests := []Request{}

	for i := 1; i <= 5; i++ {
		reminder := i % 2

		if reminder == 0 {
			requests = append(requests, Request{
				ChatID:  chatID,
				UserID:  "0fac06d5-bb77-4def-90f7-8c0c1fcffda7",
				Message: fmt.Sprintf("Message number %d", i),
			})
		} else {
			requests = append(requests, Request{
				ChatID:  chatID,
				UserID:  "41359a4e-d066-4f64-a1d8-410c2c690935",
				Message: fmt.Sprintf("Message number %d", i),
			})
		}
	}

	const url string = "http://localhost:4444/messages"

	client := http.Client{
		Timeout: time.Second * 10,
	}

	// body := map[string]any{}

	for _, request := range requests {
		data, err := json.Marshal(request)

		time.Sleep(time.Millisecond * 500)

		if err != nil {
			log.Println("Request: ", request, "\nMarshall error: ", err)
			continue
		}

		buffer := bytes.NewBuffer(data)
		response, err := client.Post(url, "application/json", buffer)

		if err != nil {
			log.Println("Request: ", request, "\nRequest error: ", err)
			continue
		}
	
		defer response.Body.Close()

		// responseData, err := io.ReadAll(response.Body)
	
		// if err != nil {
		// 	log.Println("Request: ", request, "\nReadAll error: ", err)
		// 	continue
		// }

		// err = json.Unmarshal(responseData, &body)
	
		// if err != nil {
		// 	log.Println("Request: ", request, "\nUnmarshal error: ", err)
		// 	continue
		// }
	
		log.Println("Request: ", request, "\nFinished successfully: ", response.Status)
	}
}

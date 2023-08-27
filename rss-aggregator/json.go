package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseJson(writer http.ResponseWriter, code int, payload interface{}) {
	bytes, err := json.Marshal(payload)

	if err != nil {
		log.Println("Failed to marshal payload")

		writer.WriteHeader(500)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(bytes)
}

func responseError(writer http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Code is greather than 499: ", message)
	}

	responseJson(writer, code, errorResponse{Err: message})
}

type errorResponse struct {
	Err string `json:"error"`
}

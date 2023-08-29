package main

import "net/http"

func getTimestamp(writer http.ResponseWriter, request *http.Request) {
	responseJson(writer, 200, struct{}{})
}

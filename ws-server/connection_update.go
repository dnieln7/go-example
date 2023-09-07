package main

import "github.com/gorilla/websocket"

type ConnectionUpdate struct {
	connection *websocket.Conn
	append bool
}

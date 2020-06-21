package main

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func handleQueueSockets(s socketio.Conn) error {
	log.Println(s.RemoteAddr().String() + ": socket connection established.")
	return nil
}
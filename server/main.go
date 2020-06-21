package main

import (
	"log"
	"net"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

var serverActive bool = false

//clientPlayer is a struct that holds all the pointers and information about the player as well as the connection. Some variables can be null depending on the state
//the player is in. Thus, each concurrent operation is handled such that it is either the sub-operation of an operation that it knows will ensure the existence
//of the variables it wants to use or it is said operation itself.
type clientPlayer struct {
	conn *net.Conn
	name string
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:63342")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		next.ServeHTTP(w, r)
	})
}

//main handles initial socket connections and calls the needed functions for each connection. It does so concurrently so that it may keep listening.
func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	//serve
	http.Handle("/socket.io/", server)
	fs := http.FileServer(http.Dir("../client/pages"))
	http.Handle("/", fs)

	//signify socket connections
	server.OnConnect("/queue", handleQueueSockets)

	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))

}

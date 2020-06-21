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

//main handles initial socket connections and calls the needed functions for each connection. It does so concurrently so that it may keep listening.
func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	//signify handles
	http.Handle("/socket.io/", server)
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/queue", serveHandleQueue)

	//signify socket connections
	server.OnConnect("/queue", handleQueueSockets)

	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))

	//test_main() //I moved your thing here Marco
}

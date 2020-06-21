package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var serverActive bool = false

const (
	hardDebugMode bool = true
)

//clientPlayer is a struct that holds all the pointers and information about the player as well as the connection. Some variables can be null depending on the state
//the player is in. Thus, each concurrent operation is handled such that it is either the sub-operation of an operation that it knows will ensure the existence
//of the variables it wants to use or it is said operation itself.
type clientPlayer struct {
	conn        *websocket.Conn
	name        string
	messageType int
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

	//handle
	fs := http.FileServer(http.Dir("../client/pages"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", wsEndpoint)

	//server
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

//socket stuff

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func writeMessage(p *clientPlayer, data []byte) {
	if err := p.conn.WriteMessage(p.messageType, data); err != nil {
		log.Println(err)
		return
	}
}

//reader listens to all messages from a specific client
func reader(conn *websocket.Conn) {
	var p *clientPlayer
	defer handlepanic(p)
	p = initializePlayer(conn)
	for {
		// read in a message
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		p.messageType = messageType //in case message type changes. I am doing this just to make sure. I do not think it's that important, and I don't think it's good practice so feel free to remove

		fields, flag := parseUtilsAndSignal(string(data), 2) //is it 2 fields message?

		if flag == ok {
			//2 fields message!
		} else if flag == notPONG || fields == nil {
			//GET OUTTA HERE YOU NO PLAY PONG YOU MONSTER YOU GET OUT AAA
			conn.WriteMessage(messageType, []byte("PONG INVALID"))
			conn.Close()
			panic("PONG INVALID")
		}

		//another number of fields message!
		if len(fields) == 3 {

		} else {
			//invalid number of fields boi! GET OUTTA MY SERVER YE DEGENERATE
			conn.WriteMessage(messageType, []byte("PONG INVALID"))
			conn.Close()
			panic("PONG INVALID FIELD NUMBERS")
		}

		if hardDebugMode {
			// print out that message for extra clarity
			fmt.Println(p.name + ": " + string(data))
		}

	}
}

func initializePlayer(conn *websocket.Conn) *clientPlayer {
	// read in a message
	messageType, data, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}

	fields, flag := parseUtilsAndSignal(string(data), 2) //one for PONG one for name expected.

	if flag != ok {
		//FILTHY INVALID PROTOCOLS. GET OUTTA HERE!
		conn.WriteMessage(messageType, []byte("PONG INVALID"))
		conn.Close()
		panic("INVALID WITH FLAG: " + parseFlagToString(flag))
	}

	// print out that message for clarity
	fmt.Println(string(data))

	return &clientPlayer{
		conn:        conn,
		messageType: messageType,
		name:        fields[1],
	}
}

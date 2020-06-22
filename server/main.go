package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var serverActive bool = false

const (
	hardDebugMode bool = true
)

//2 fields are for direct messages. 3 fields are for direct + variable pass
var handledClientCalls = map[clientCallSpecification]clientMessageHandle{
	clientCallSpecification{isDirectMessage: true, msgBase: "PONG QUEUE"}:       queuePlayer,
	clientCallSpecification{isDirectMessage: false, msgBase: "PONG UPDATE-DIR"}: playerUpdateDesiredDirection,
}

type clientCallSpecification struct {
	msgBase         string
	isDirectMessage bool
}

type clientMessageHandle func(*playerRequest, string)

//clientPlayer is a struct that holds all the pointers and information about the player as well as the connection. Some variables can be null depending on the state
//the player is in. Thus, each concurrent operation is handled such that it is either the sub-operation of an operation that it knows will ensure the existence
//of the variables it wants to use or it is said operation itself.
type clientPlayer struct {
	queued               bool
	conn                 *websocket.Conn
	name                 string
	messageType          int
	activeGame           *game
	valid                bool
	writeChannel         chan *writeRequest
	disconnectChannel    chan interface{} //an empty interface is used as no data is passed, just the existence of the signal is enough
	m                    sync.RWMutex     //Thou shalt not both read and write from a connection from different threads. Alas, thou must nonetheless keep both threads listening- one from client and one for write requests within server.
	tendedPlayersElement *list.Element    //its position in the connected list.
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*:*")
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
	fs := http.FileServer(http.Dir("../client/render"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", wsEndpoint)

	//server
	log.Println("Serving at localhost:5000...")
	serverActive = true
	go queueSystem() //initialize the queue system that shall perpetually listen and listen for more and more players. This is its punishment for being a thread
	log.Fatal(http.ListenAndServe(":5000", nil))

	serverActive = false
}

//socket stuff

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //upgrade for all clients for now. Maybe Docker in less hacky application
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("PONG HEY THERE!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func writeMessage(p *clientPlayer, data []byte) {
	if !p.valid { //invalids have no place here
		return
	}
	if err := p.conn.WriteMessage(p.messageType, data); err != nil {
		log.Println(err)
		return
	}
}

var tendedPlayers *list.List = list.New()

var tendedPlayersMutex sync.RWMutex = sync.RWMutex{}

//READER

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
		conn:              conn,
		messageType:       messageType,
		name:              fields[1],
		valid:             true,
		disconnectChannel: make(chan interface{}, 3), //the two channel size numbers can be reduced if this is scaled. I am leaving high numbers for debugging purposes.
		writeChannel:      make(chan *writeRequest, 6),
	}
}

func integratePlayerIntoTendedMass(p *clientPlayer) {
	tendedPlayersMutex.Lock()
	p.tendedPlayersElement = tendedPlayers.PushBack(p)
	tendedPlayersMutex.Unlock()
}

//reader listens to all messages from a specific client and then passes it down to everything. it is basically God, if there was different instances of God for each person.
//It heeds our prayers so that we may gracefully play pong and marvel at the beauty of the front-end.
func reader(conn *websocket.Conn) {
	var p *clientPlayer
	defer handlepanic(p)
	p = initializePlayer(conn)
	integratePlayerIntoTendedMass(p)
	fmt.Println("Player initialized: " + p.name)
	go channelsListener(p) //one thread listens within, other listens out. the one who listens within also writes the messages.
	p.writeChanneledMessage("PONG NAME-OK")
	for p.valid && serverActive {
		// read in a message
		messageType, data, err := conn.ReadMessage()
		if !p.valid || !serverActive {
			return
		}
		p.m.RLock()
		fmt.Println(p.name + ": " + string(data))
		if err != nil {
			p.m.RUnlock()
			log.Println(err)
			p.writeChanneledMessage("PONG INVALID")
			handleDisconnectPlayer(p)
			panic("PONG INVALID")
		}
		p.messageType = messageType //in case message type changes. I am doing this just to make sure. I do not think it's that important, and I don't think it's good practice so feel free to remove

		fields, flag := parseUtilsAndSignal(string(data), 2) //is it 2 fields message?

		fmt.Println("right before checking flag")
		if flag == ok {
			//2 fields message!
			dataString := fields[0] + " " + fields[1]
			specification := clientCallSpecification{
				isDirectMessage: true,
				msgBase:         dataString,
			}
			if val, ok := handledClientCalls[specification]; ok {
				fmt.Println("hey we're handling client calls now")
				val(&playerRequest{message: dataString, p: p}, "")
			} else {
				p.m.RUnlock()
				p.writeChanneledMessage("PONG INVALID")
				handleDisconnectPlayer(p)
				panic("PONG INVALID DIRECT MESSAGE: " + string(data))
			}
		} else if flag == notPONG || fields == nil {
			//GET OUTTA HERE YOU NO PLAY PONG YOU MONSTER YOU GET OUT AAA
			p.m.RUnlock()
			p.writeChanneledMessage("PONG INVALID")
			handleDisconnectPlayer(p)
			panic("PONG INVALID")
		} else if len(fields) == 3 { //this is for directive and argument
			fmt.Println("3 field messages was sent.")
			dataString := fields[0] + " " + fields[1]
			specification := clientCallSpecification{
				isDirectMessage: false,
				msgBase:         dataString,
			}
			val, ok := handledClientCalls[specification]
			fmt.Println(ok)
			if ok {
				val(&playerRequest{message: dataString, p: p}, fields[2])
			} else {
				p.m.RUnlock()
				p.writeChanneledMessage("PONG INVALID")
				handleDisconnectPlayer(p)
				panic("PONG INVALID DIRECT MESSAGE: " + string(data))
			}
		} else {
			//invalid number of fields boi! GET OUTTA MY SERVER YE DEGENERATE
			p.m.RUnlock()
			p.writeChanneledMessage("PONG INVALID")
			handleDisconnectPlayer(p)
			panic("PONG INVALID FIELD NUMBER")
		}

		p.m.RUnlock()

		if hardDebugMode {
			// print out that message for extra clarity
			fmt.Println(p.name + ": " + string(data))
		}

	}
}

func (p *clientPlayer) writeChanneledMessage(msg string) {
	p.writeChannel <- &writeRequest{message: msg}
}

func handleDisconnectPlayer(p *clientPlayer) {
	p.disconnectChannel <- struct{}{} //first time seeing this might be weird. interface{} is type. interface{}{} is instantiation ;)
}

//INWARDS LISTENER AND WRITER

//channelsListener ensures that the channels of a player is always functional, but only if the player is valid.
func channelsListener(p *clientPlayer) {
	for p.valid && serverActive {
		if !p.valid {
			return
		}
		select {
		case w := <-p.writeChannel:
			p.m.Lock()
			writeMessage(p, []byte(w.message))
			p.m.Unlock()
			break
		case <-p.disconnectChannel:
			p.m.Lock()
			fmt.Println("disconnecting through channel: " + p.name)
			writeMessage(p, []byte("PONG CLOSE"))
			removeClient(p)
			p.conn.Close()
			p.m.Unlock()
			p = nil
			return
		}
	}
}

func removeClient(p *clientPlayer) {
	p.valid = false

	tendedPlayersMutex.Lock()
	tendedPlayers.Remove(p.tendedPlayersElement)
	tendedPlayersMutex.Unlock()

	if p.activeGame != nil {
		//oh noes must deal with other players. quitting scum!
	}

}

//INNER COMMUNICATION TYPES
type writeRequest struct {
	message string
}

type playerRequest struct {
	message string
	p       *clientPlayer
}

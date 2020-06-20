package main

import (
	"fmt"
	"log"
	"net"
)

var serverActive bool = false

type player struct {
	conn *net.Conn
	name string
}

//main handles initial socket connections and calls the needed functions for each connection. It does so concurrently so that it may keep listening.
func main() {
	li, err := net.Listen("tcp", ":52515")
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println("Now listening on port 52515...")

	serverActive = true

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("New connection from " + conn.RemoteAddr().String())
		go handleGameConnection(&conn)
	}
}

//handleGameConnection is the first function called concurrently for a client and it calls all other needed functions as well as constructs the client object.
func handleGameConnection(conn *net.Conn) {
	fmt.Println("Handling connection for: " + (*conn).RemoteAddr().String())

}

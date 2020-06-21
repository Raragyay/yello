package main

// import (
// 	"log"
// 	"net/http"
// )

// func serveIndex(w http.ResponseWriter, req *http.Request) {
// 	http.ServeFile(w, req, "../client/pages/index.html")
// }

// func serveHandleQueue(w http.ResponseWriter, req *http.Request) {
// 	name := req.FormValue("name")
// 	if len(name) == 0 {
// 		log.Println(req.RemoteAddr + " sent bad queue request with input: " + name)
// 		serveBadRequest(w, req)
// 		return
// 	}
// 	log.Println(req.RemoteAddr + " sent queue-up request with name: " + name)
// 	http.ServeFile(w, req, "../client/pages/queue.html")
// }

// func serveBadRequest(w http.ResponseWriter, req *http.Request) {
// 	http.ServeFile(w, req, "../client/pages/badRequest.html")
// }

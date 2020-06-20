package main

import (
	"net/http"
)

func serveIndex(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "../client/pages/index.html")
}

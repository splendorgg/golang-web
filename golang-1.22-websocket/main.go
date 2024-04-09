package main

import (
	"log"
	"net/http"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	go handleMessages()
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", serveWS)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

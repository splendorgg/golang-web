package main

import (
	"log"
	"net/http"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {
	hub := newHub()
	go hub.Run()
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(":3080", nil))
}

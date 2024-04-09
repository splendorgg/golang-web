package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}
type Message struct {
	Text string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan string)
)

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Error writing message:", err)
				return
			}
		}
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket error: ", err)
		return
	}
	defer wsconn.Close()

	client := &Client{
		conn: wsconn,
	}
	clients[client] = true

	for {

		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("readin error:", err)
			delete(clients, client)
			break
		}

		fmt.Println("received message", string(msg))

		broadcast <- string(msg)
	}

}

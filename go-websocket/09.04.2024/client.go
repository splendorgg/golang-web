package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	wsconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket error: ", err)
		return
	}
	defer wsconn.Close()
	id := uuid.New()
	client := &Client{
		id:   id.String(),
		hub:  hub,
		conn: wsconn,
		send: make(chan []byte),
	}
	client.hub.register <- client

	go client.readMessages()
	go client.writeMessages()

}

func (client *Client) readMessages() {
	defer func() {
		client.conn.Close()
		client.hub.unregister <- client
	}()

	for {
		_, readmsg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("readin error:", err)
			break
		}
		msg := &WSMessage{}
		reader := bytes.NewReader(readmsg)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			log.Println(err)
		}
		client.hub.broadcast <- &Message{ClientID: client.id, Text: msg.Text}

	}

}

func (client *Client) writeMessages() {
	defer func() {
		client.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				return
			}

		}
	}
}

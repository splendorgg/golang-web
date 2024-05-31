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

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	id := uuid.New()
	client := &Client{id: id.String(), hub: hub, conn: conn, send: make(chan []byte)}
	client.hub.register <- client

	go client.writeMessage()
	go client.readMessage()

}

func (c *Client) readMessage() {
	defer func() {
		c.conn.Close()
		c.hub.unregister <- c
	}()

	for {
		_, text, err := c.conn.ReadMessage()
		log.Printf("value %v", string(text))
		if err != nil {
			log.Println(err)
			break
		}

		msg := &WSMessage{}
		reader := bytes.NewReader(text)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(msg)
		if err != nil {
			log.Println(err)
		}
		c.hub.broadcast <- &Message{ClientID: c.id, Text: msg.Text}
	}
}

func (c *Client) writeMessage() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(msg)
			}

			if err := w.Close(); err != nil {
				return
			}

		}
	}
}

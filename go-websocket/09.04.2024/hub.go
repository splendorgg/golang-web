package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

type Message struct {
	ClientID string
	Text     string
}
type WSMessage struct {
	Text    string      `json:"text"`
	Headers interface{} `json:"HEADERS"`
}
type Hub struct {
	clients    map[*Client]bool
	messages   []*Message
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
			log.Printf("client registered %s:", client.id)
			for _, msg := range hub.messages {
				client.send <- getMessageTemplate(msg)
			}
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				log.Printf("client unregistered %s:", client.id)
				close(client.send)
				delete(hub.clients, client)
			}
		case message := <-hub.broadcast:
			hub.messages = append(hub.messages, message)
			for client := range hub.clients {
				select {
				case client.send <- getMessageTemplate(message):
				default:
					close(client.send)
					delete(hub.clients, client)
				}
				fmt.Println(client)
			}
		}
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Fatalf("template parsing:%s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("template parsing:%s", err)
	}
	return renderedMessage.Bytes()

}

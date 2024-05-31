package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request)bool  {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn,err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	
	for{
		messageType,message,err:=conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return 
		}
		log.Printf("Received message: %v",string(message))
		if err :=conn.WriteMessage(messageType,message);err !=nil{
			log.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws",handler)
	log.Fatal(http.ListenAndServe(":8080",nil))
}

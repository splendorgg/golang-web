/* package main

import (
	"fmt"
	"net/http"





	"github.com/gorilla/websocket"
)

type Room struct {
	Name     string
	Members  []*websocket.Conn
	Messages []string
}

var rooms map[string]*Room

func createRoom(name string) {
	rooms[name] = &Room{
		Name:     name,
		Members:  make([]*websocket.Conn, 0),
		Messages: make([]string, 0),
	}
}

func joinRoom(conn *websocket.Conn, roomName string) {
	room, ok := rooms[roomName]
	if ok {
		room.Members = append(room.Members, conn)
		// Gönderen için hoş geldin mesajı gönder
		conn.WriteMessage(websocket.TextMessage, []byte("Hoş geldiniz "+roomName+" odasına"))
	} else {
		// Oda bulunamadı
		conn.WriteMessage(websocket.TextMessage, []byte("Oda bulunamadı"))
	}
}

func main() {
	rooms = make(map[string]*Room)

	http.HandleFunc("/createRoom", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("name")
		createRoom(roomName)
		fmt.Fprintf(w, "Oda oluşturuldu: %s", roomName)
	})

	http.HandleFunc("/joinRoom", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("name")
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			fmt.Println("WebSocket bağlantısı oluşturulamadı:", err)
			return
		}
		defer conn.Close()

		joinRoom(conn, roomName)
	})

	http.ListenAndServe(":8080", nil)
}
*/

package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Client struct {
	conn     *websocket.Conn
	username string
}

type Room struct {
	name     string
	clients  map[*Client]bool
	messages chan []byte
}

var clients = make(map[*Client]bool)
var rooms = make(map[string]*Room)
var clientRoom = make(map[*Client]*Room)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Hatası:", err)
		return
	}

	username := r.URL.Query().Get("username")
	roomName := r.URL.Query().Get("room")

	client := &Client{
		conn:     conn,
		username: username,
	}

	joinRoom(client, roomName)
	defer leaveRoom(client)

	handleMessage(client)
}

func joinRoom(client *Client, roomName string) {
	room, exists := rooms[roomName]
	if !exists {
		room = &Room{
			name:     roomName,
			clients:  make(map[*Client]bool),
			messages: make(chan []byte),
		}
		rooms[roomName] = room

		go func() {
			for {
				message := <-room.messages
				for client := range room.clients {
					err := client.conn.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						log.Println("Mesaj Gönderme Hatası:", err)
						return
					}
				}
			}
		}()
	}

	room.clients[client] = true
	clientRoom[client] = room
}

func leaveRoom(client *Client) {
	if room, exists := clientRoom[client]; exists {
		delete(room.clients, client)
		delete(clientRoom, client)
	}
}

func handleMessage(client *Client) {
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Mesaj Okuma Hatası:", err)
			return
		}

		broadcastMessage(message, client)
	}
}

func broadcastMessage(message []byte, sender *Client) {
	username := sender.username

	for client := range clients {
		if clientRoom[client] == clientRoom[sender] {
			msg := fmt.Sprintf("%s: %s", username, message)
			err := client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Mesaj Gönderme Hatası:", err)
				return
			}
		}
	}
}
func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	roomName := r.URL.Query().Get("room")

	// İstemciyi oluştur ve odaya katıl
	client := &Client{
		conn:     nil, // WebSocket bağlantısı henüz oluşturulmadı
		username: username,
	}
	joinRoom(client, roomName)
}

func leaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Kullanıcı adını ve odanın adını al
	username := r.URL.Query().Get("username")
	roomName := r.URL.Query().Get("room")

	// İstemciyi bul ve odadan çıkart
	for client := range clients {
		if client.username == username {
			if clientRoom[client] != nil && clientRoom[client].name == roomName {
				leaveRoom(client)
				break
			}
		}
	}
}
func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	http.HandleFunc("/join", joinRoomHandler)
	http.HandleFunc("/leave", leaveRoomHandler)
}

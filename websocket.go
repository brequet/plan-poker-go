package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	// Add more client-specific fields here if needed
}

type Room struct {
	clients map[*Client]bool
	// Add more room-specific fields here if needed
}

var upgrader = websocket.Upgrader{}
var rooms = make(map[string]*Room)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO: check origin

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	log.Println("New connection !")

	// Get the room name from the request URL or other parameters
	roomName := "room1" // Replace with actual room name (you'll need to extract it from the request)

	// Create a new client
	client := &Client{conn: conn}

	// Get the room for the given name or create a new one if it doesn't exist
	room, ok := rooms[roomName]
	if !ok {
		room = &Room{
			clients: make(map[*Client]bool),
		}
		rooms[roomName] = room
	}

	// Add the client to the room
	room.clients[client] = true

	log.Printf("Room %s (%d client(s) in the room)", roomName, len(room.clients))
	log.Println("client ->", client)
	go handleMessages(client, room)
}

func handleMessages(client *Client, room *Room) {
	defer func() {
		// Clean up resources when client disconnects
		log.Println("Client deconnecting", client)
		delete(room.clients, client)
		client.conn.Close()
	}()

	for {
		// Read message from the client
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Println("msg", message)
		// Handle different message types here (e.g., poker planning estimates)

		// Broadcast the received message to all other clients in the room
		for c := range room.clients {
			if c != client {
				err := c.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error writing message:", err)
					break
				}
			}
		}
	}
}

func main() {
	log.Println("Starting server..")
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}

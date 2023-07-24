package main

import (
	"encoding/json"
	"log"
	"net/http"
	// "os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	log.Printf("Room %s (%d client(s) in the room) : new client %s", roomName, len(room.clients), client.conn.RemoteAddr().String())
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

		log.Println("msg :", string(message))
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

// HTTP endpoint to create a new room
func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get room name and other data
	var roomData struct {
		RoomName string `json:"roomName"`
		// Other relevant fields
	}
	if err := json.NewDecoder(r.Body).Decode(&roomData); err != nil {
		log.Printf("Error in createRoomHandler: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a new room with the provided room name
	// newRoom := createNewRoom(roomData.RoomName)
	log.Printf("Creating room : %s\n", roomData.RoomName)

	// Respond with the room details (e.g., room ID) to the frontend
	response := struct {
		RoomID string `json:"roomId"`
		// Other relevant fields
	}{
		RoomID: roomData.RoomName,
		// RoomID: newRoom.ID,
		// Populate other response data if needed
	}

	// Send the response back to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	log.Println("Server started on 127.0.0.1:8080")

	router := mux.NewRouter()

	router.HandleFunc("/ws", wsHandler)
	router.HandleFunc("/api/room", createRoomHandler)

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	// originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

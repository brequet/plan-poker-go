package main

import (
	"encoding/json"
	"log"
	"net/http"

	rm "github.com/baptiste-requet/plan-poker-go/rooms-manager"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	user     *rm.User
	roomCode string
}

var upgrader = websocket.Upgrader{}

var clients = make(map[*Client]bool)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO: check origin

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	client := &Client{
		conn:     conn,
		user:     nil,
		roomCode: "",
	}
	clients[client] = true
	log.Printf("New connection from %s", client.conn.RemoteAddr())

	go handleMessages(client)
}

func handleMessages(client *Client) {
	defer func() {
		// Clean up resources when client disconnects
		log.Println("Client deconnecting", client)
		if client.roomCode != "" {
			rm.DisconnectUserFromRoom(client.user, client.roomCode)
		}
		delete(clients, client)
		client.conn.Close()
	}()

	for {
		var receivedMessage Message
		err := client.conn.ReadJSON(&receivedMessage)
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			break
		}

		// React to different message types
		log.Println("Received message", receivedMessage)
		switch receivedMessage.Type {
		case JOIN_ROOM:
			var joinRoomMessage JoinRoomMessage
			if err := json.Unmarshal(receivedMessage.Payload, &joinRoomMessage); err != nil {
				log.Println("Unmarshal error join")
				return
			}
			log.Println("joinRoomMessage", joinRoomMessage)

			log.Printf("User '%s' joined room '%s'", joinRoomMessage.Nickname, joinRoomMessage.RoomCode)
			// TODO: Handle user joining the room, e.g., store the user in the room


		case USER_JOINED:
			var userJoinedMessage UserJoinedMessage
			if err := json.Unmarshal(receivedMessage.Payload, &userJoinedMessage); err != nil {
				// Handle unmarshal error
				log.Println("Unmarshal error joined")
				return
			}
			log.Println("userJoinedMessage", userJoinedMessage)

			// Handle the USER_JOINED message and use userJoinedMessage.UserName


		default:
			log.Println("Received unsupported message type:", receivedMessage.Type)
		}

		// Read message from the client
		// _, message, err := client.conn.ReadMessage()
		// if err != nil {
		// 	log.Println("Error reading message:", err)
		// 	break
		// }

		// log.Println("msg :", string(message))
		// // Handle different message types here (e.g., poker planning estimates)

		// // Broadcast the received message to all other clients in the room
		// for c := range clients { // TODO: only client in room
		// 	if c != client {
		// 		err := c.conn.WriteMessage(websocket.TextMessage, message)
		// 		if err != nil {
		// 			log.Println("Error writing message:", err)
		// 			break
		// 		}
		// 	}
		// }
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
	newRoom := rm.CreateRoom(roomData.RoomName)
	log.Printf("Created room : [%s] %s\n", newRoom.Code, newRoom.Name)

	// Respond with the room details (e.g., room ID) to the frontend
	response := struct {
		RoomCode string `json:"roomCode"`
		RoomName string `json:"roomName"`
	}{
		RoomCode: newRoom.Code,
		RoomName: newRoom.Name,
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

	/*
		TODO: endpoint for
		- fetching room info (name, code, connected user..)
		- fetching all rooms (for admin purpose)
		-
	*/

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	// originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

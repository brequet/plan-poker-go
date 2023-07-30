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

// TODO: better logging (with levels)

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
		handleDisconnection(client)
	}()

	for {
		var receivedMessage ReceiveMessage
		err := client.conn.ReadJSON(&receivedMessage)
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			break
		}

		switch receivedMessage.Type {
		case JOIN_ROOM:
			var joinRoomMessage JoinRoomMessage
			if err := json.Unmarshal(receivedMessage.Payload, &joinRoomMessage); err != nil {
				log.Fatalf("Could not unmarshal %s message for message type: %s", receivedMessage.Type, err)
				break
			}
			onRoomJoinEvent(client, joinRoomMessage)

		case USER_JOINED: // TODO: remove this event, this is an event back->front, makes no sense here
			var userJoinedMessage UserJoinedMessage
			if err := json.Unmarshal(receivedMessage.Payload, &userJoinedMessage); err != nil {
				log.Fatalf("Could not unmarshal %s message for message type: %s", receivedMessage.Type, err)
				break
			}
			log.Println("userJoinedMessage", userJoinedMessage)

		default:
			log.Println("Received unsupported message type:", receivedMessage.Type)
		}
	}
}

func handleDisconnection(client *Client) {
	// Clean up resources when client disconnects
	log.Println("Client disconnecting", client.conn.RemoteAddr().String())

	if client.roomCode != "" && client.user != nil {
		rm.DisconnectUserFromRoom(client.user, client.roomCode)

		disconnectionMessage := SendMessage{
			Type: USER_DISCONNECTED,
			Payload: DisconnectionMessage{
				User: User{
					UserName: client.user.Nickname,
					Uuid:     client.user.Uuid,
				},
			},
		}
		go broadcastMessageToOtherClientsInRoom(disconnectionMessage, client.user, client.roomCode)
	}

	delete(clients, client)
	client.conn.Close()
}

func onRoomJoinEvent(client *Client, joinRoomMessage JoinRoomMessage) { // todo in event.go ?
	user := rm.ConnectNewUserToRoom(joinRoomMessage.Nickname, joinRoomMessage.RoomCode)
	if user == nil {
		log.Fatalf("User %s could not join the room %s", joinRoomMessage.Nickname, joinRoomMessage.RoomCode)
		return
	}
	client.roomCode = joinRoomMessage.RoomCode
	client.user = user

	connectedUsers := []User{} // TODO: add in message existing users in the room
	for _, connectedUser := range rm.GetAllUserFromRoomByRoomCode(client.roomCode) {
		if connectedUser != user {
			connectedUsers = append(connectedUsers, User{
				UserName: connectedUser.Nickname,
				Uuid: connectedUser.Uuid,
			})
		}
	}

	confirmConnexionMessage := SendMessage{
		Type: CONFIRM_CONNECTION,
		Payload: ConfirmConnectionMessage{
			User: User{
				UserName: user.Nickname,
				Uuid:     user.Uuid,
			},
			ConnectedUsers: connectedUsers,
		},
	}
	err := client.conn.WriteJSON(confirmConnexionMessage)
	if err != nil {
		log.Fatalf("Could not send room joining confirmation to user %s [%s]", user.Nickname, client.roomCode)
		return
	}

	userJoinedMessage := SendMessage{
		Type: USER_JOINED,
		Payload: UserJoinedMessage{
			User: User{
				UserName: user.Nickname,
				Uuid:     user.Uuid,
			},
		},
	}
	go broadcastMessageToOtherClientsInRoom(userJoinedMessage, user, joinRoomMessage.RoomCode)
}

func getAllOtherClientsInRoomByRoomCode(excludedUser *rm.User, roomCode string) (foundClients []*Client) {
	for _, user := range rm.GetAllUserFromRoomByRoomCode(roomCode) {
		if excludedUser.Uuid == user.Uuid {
			continue
		}

		for client, _ := range clients {
			if user.Uuid == client.user.Uuid {
				foundClients = append(foundClients, client)
			}
		}
	}
	return foundClients
}

func broadcastMessageToClients(message interface{}, clients []*Client) {
	for _, client := range clients {
		client.conn.WriteJSON(message)
	}
}

func broadcastMessageToOtherClientsInRoom(message interface{}, user *rm.User, roomCode string) {
	broadcastMessageToClients(message, getAllOtherClientsInRoomByRoomCode(user, roomCode))
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

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	// originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

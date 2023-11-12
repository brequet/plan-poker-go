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
		case JOIN_ROOM: // TODO: generic func(messageType, callback func)
			var joinRoomMessage JoinRoomMessage
			if err := json.Unmarshal(receivedMessage.Payload, &joinRoomMessage); err != nil {
				log.Printf("Could not unmarshal %s message for message type: %s", receivedMessage.Type, err)
				break
			}
			onRoomJoinEvent(client, joinRoomMessage)

		case SUBMIT_ESTIMATE:
			var submitEstimateMessage SubmitEstimateMessage
			if err := json.Unmarshal(receivedMessage.Payload, &submitEstimateMessage); err != nil {
				log.Printf("Could not unmarshal %s message for message type: %s", receivedMessage.Type, err)
				break
			}
			onSubmitEstimateEvent(client, submitEstimateMessage)

		case REVEAL_ESTIMATE:
			var revealEstimateMessage RevealEstimateMessage
			if err := json.Unmarshal(receivedMessage.Payload, &revealEstimateMessage); err != nil {
				log.Printf("Could not unmarshal %s message for message type: %s", receivedMessage.Type, err)
				break
			}
			onRevealEstimateEvent(client, revealEstimateMessage)

		case RESET_PLANNING:
			onResetPlanningEvent(client)

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
			Payload: UserDisconnectedMessage{
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
	// TODO: check if user exist (ip:port ?) -> ex: if user F5 refresh page, keep connection if possible
	user := rm.ConnectNewUserToRoom(joinRoomMessage.Nickname, joinRoomMessage.RoomCode)
	if user == nil {
		log.Printf("User %s could not join the room %s", joinRoomMessage.Nickname, joinRoomMessage.RoomCode)
		// TODO: return if fails and adapt in front -> e.g. if user joins nickname choices while the room still exist, and then validate the name when the rooms expire, we get a error 500
		return
	}
	client.roomCode = joinRoomMessage.RoomCode
	client.user = user

	connectedUsers := []User{}
	for _, connectedUser := range rm.GetAllUserFromRoomByRoomCode(client.roomCode) {
		connectedUsers = append(connectedUsers, User{
			UserName: connectedUser.Nickname,
			Uuid:     connectedUser.Uuid,
			Estimate: connectedUser.Estimate,
		})
	}

	room := rm.FindRoomByRoomCode(client.roomCode)

	confirmConnexionMessage := SendMessage{
		Type: CONFIRM_CONNECTION,
		Payload: ConfirmConnectionMessage{
			User: User{
				UserName: user.Nickname,
				Uuid:     user.Uuid,
			},
			ConnectedUsers: connectedUsers,
			Room: Room{
				RoomCode:           room.Code,
				RoomName:           room.Name,
				IsEstimateRevealed: room.IsEstimateRevealed,
			},
		},
	}
	err := client.conn.WriteJSON(confirmConnexionMessage)
	if err != nil {
		log.Printf("Could not send room joining confirmation to user %s [%s]: %s", user.Nickname, client.roomCode, err)
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
	go broadcastMessageToOtherClientsInRoom(userJoinedMessage, client.user, client.roomCode)
}

func onSubmitEstimateEvent(client *Client, submitEstimateMessage SubmitEstimateMessage) {
	err := rm.SubmitEstimate(client.user, client.roomCode, submitEstimateMessage.Estimate)
	if err != nil {
		return
	}

	confirmEstimateSubmission := SendMessage{
		Type: CONFIRM_ESTIMATE_SUBMISSION,
		Payload: ConfirmEstimateSubmissionMessage{
			Estimate: submitEstimateMessage.Estimate,
		},
	}
	err = client.conn.WriteJSON(confirmEstimateSubmission)
	if err != nil {
		log.Printf("Could not send estimate submission confirmation to user %s [%s]: %s", client.user.Nickname, client.roomCode, err)
		return
	}

	userSubmittedMessage := SendMessage{
		Type: USER_SUBMITTED_ESTIMATE,
		Payload: UserSubmittedEstimate{
			User: User{
				UserName: client.user.Nickname,
				Uuid:     client.user.Uuid,
			},
			Estimate: submitEstimateMessage.Estimate,
		},
	}
	go broadcastMessageToAllClientsInRoom(userSubmittedMessage, client.roomCode)
}

func onRevealEstimateEvent(client *Client, revealEstimateMessage RevealEstimateMessage) {
	newShouldRevealEstimate, err := rm.ToggleShouldRevealEstimateForRoom(client.roomCode)
	if err != nil {
		log.Printf("Could not reveal estimates for room [%s]: %s", client.roomCode, err)
		return
	}

	revealEstimateMessageToSend := SendMessage{
		Type: REVEAL_ESTIMATE,
		Payload: RevealEstimateMessage{
			ShouldReveal: newShouldRevealEstimate,
		},
	}
	go broadcastMessageToAllClientsInRoom(revealEstimateMessageToSend, client.roomCode)
}

func onResetPlanningEvent(client *Client) {
	err := rm.ResetPlanningForRoom(client.roomCode)
	if err != nil {
		log.Printf("Could not reset planning for room [%s]: %s", client.roomCode, err)
		return
	}

	planningResetMessage := SendMessage{
		Type:    PLANNING_RESETED,
		Payload: PlanningResetedMessage{},
	}
	go broadcastMessageToAllClientsInRoom(planningResetMessage, client.roomCode)
}

func getAllClientsInRoomByRoomCode(roomCode string) (foundClients []*Client) {
	for _, user := range rm.GetAllUserFromRoomByRoomCode(roomCode) {
		for client, _ := range clients {
			if user.Uuid == client.user.Uuid {
				foundClients = append(foundClients, client)
			}
		}
	}
	return foundClients
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

func broadcastMessageToAllClientsInRoom(message interface{}, roomCode string) {
	broadcastMessageToClients(message, getAllClientsInRoomByRoomCode(roomCode))
}

func broadcastMessageToOtherClientsInRoom(message interface{}, excludedUser *rm.User, roomCode string) {
	broadcastMessageToClients(message, getAllOtherClientsInRoomByRoomCode(excludedUser, roomCode))
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

	if roomData.RoomName == "" {
		// RoomName empty, return an error response
		http.NotFound(w, r) /// TODO: bad param
		return
	}

	// Create a new room with the provided room name
	newRoom := rm.CreateRoom(roomData.RoomName)
	log.Printf("Created room : [%s] %s\n", newRoom.Code, newRoom.Name)

	// Respond with the room details (e.g., room ID) to the frontend
	response := mapRmRoomToRoom(newRoom)

	// Send the response back to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HTTP endpoint to get a room
func getRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the room code from the URL path parameter
	roomCode := mux.Vars(r)["roomCode"]

	// Retrieve the room from your data store using the room code
	room := rm.FindRoomByRoomCode(roomCode)

	if room == nil {
		// Room not found, return an error response
		http.NotFound(w, r)
		return
	}

	// Respond with the room details (e.g., room ID) to the frontend
	response := mapRmRoomToRoom(room)

	// Send the response back to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	log.Println("Server started on 127.0.0.1:8080")

	router := mux.NewRouter()

	router.HandleFunc("/ws", wsHandler)
	router.HandleFunc("/api/room", createRoomHandler).Methods("POST")
	router.HandleFunc("/api/room/{roomCode}", getRoomHandler).Methods("GET")

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

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

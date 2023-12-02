package rest

import (
	"encoding/json"
	"log"

	"net/http"

	rm "github.com/baptiste-requet/plan-poker-go/rooms-manager"
	"github.com/gorilla/websocket"
)

type User struct {
	UserName string `json:"userName"`
	Uuid     string `json:"uuid"`
	Estimate string `json:"estimate"`
}

type Room struct {
	RoomCode           string `json:"roomCode"`
	RoomName           string `json:"roomName"`
	IsEstimateRevealed bool   `json:"isEstimateRevealed"`
}

type MessageType string

/*
	Both received and send message
*/

const (
	REVEAL_ESTIMATE MessageType = "reveal_estimate"
)

type RevealEstimateMessage struct {
	ShouldReveal bool `json:"shouldReveal"`
}

/*
	Received messages
*/

const (
	JOIN_ROOM       MessageType = "join_room"
	SUBMIT_ESTIMATE MessageType = "submit_estimate"
	RESET_PLANNING  MessageType = "reset_planning"
)

type ReceiveMessage struct {
	Type    MessageType
	Payload json.RawMessage
}

type JoinRoomMessage struct {
	RoomCode string `json:"roomCode"`
	Nickname string `json:"nickname"`
}

type SubmitEstimateMessage struct {
	Estimate string `json:"estimate"`
}

type ResetPlanningMessage struct {
}

/*
	Sended messages
*/

const (
	USER_JOINED                 MessageType = "user_joined"
	USER_DISCONNECTED           MessageType = "user_disconnected"
	CONFIRM_CONNECTION          MessageType = "confirm_connection"
	CONFIRM_ESTIMATE_SUBMISSION MessageType = "confirm_estimate_submission"
	USER_SUBMITTED_ESTIMATE     MessageType = "user_submitted_estimate"
	ESTIMATE_SUBMITTED          MessageType = "estimate_submitted"
	ESTIMATE_REVEALED           MessageType = "estimate_revealed"
	PLANNING_RESETED            MessageType = "planning_reseted"
)

type SendMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type UserJoinedMessage struct {
	User User `json:"user"`
}

type UserDisconnectedMessage struct {
	User User `json:"user"`
}

type ConfirmConnectionMessage struct {
	User           User   `json:"user"`
	ConnectedUsers []User `json:"connectedUsers"`
	Room           Room   `json:"room"`
}

type ConfirmEstimateSubmissionMessage struct {
	Estimate string `json:"estimate"`
}

type UserSubmittedEstimate struct {
	User     User   `json:"user"`
	Estimate string `json:"estimate"`
}

type PlanningResetedMessage struct {
}

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

	// TODO if room doesn't exist, close connection

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
	// TODO: check if user exist (ip ?) -> ex: if user F5 refresh page, keep connection if possible
	user := rm.ConnectNewUserToRoom(joinRoomMessage.Nickname, joinRoomMessage.RoomCode)
	if user == nil {
		// TODO: return if fails and adapt in front -> e.g. if user joins nickname choices while the room still exist, and then validate the name when the rooms expire, we get a error 500
		// Or maybe just do not accept the web socket connection ?
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
	if client.user == nil {
		return
	}

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
			if client.user != nil && user.Uuid == client.user.Uuid {
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

		for client := range clients {
			if client.user != nil && user.Uuid == client.user.Uuid {
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

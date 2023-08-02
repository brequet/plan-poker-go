package main

import (
	"encoding/json"

	rm "github.com/baptiste-requet/plan-poker-go/rooms-manager"
)

type User struct {
	UserName string `json:"userName"`
	Uuid     string `json:"uuid"`
}

type Room struct {
	RoomCode string `json:"roomCode"`
	RoomName string `json:"roomName"`
}

type MessageType string

/*
	Received messages
*/

const (
	JOIN_ROOM       MessageType = "join_room"
	SUBMIT_ESTIMATE MessageType = "submit_estimate"
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
	// TaskId   string `json:"taskId"`
	Estimate string `json:"estimate"`
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
	RESET_PLANNING              MessageType = "reset_planning"
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

func mapRmRoomToRoom(rmRoom *rm.Room) *Room {
	return &Room{
		RoomCode: rmRoom.Code,
		RoomName: rmRoom.Name,
	}
}

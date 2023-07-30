package main

import (
	"encoding/json"
)

const (
	JOIN_ROOM          MessageType = "join_room"
	USER_JOINED        MessageType = "user_joined"
	USER_DISCONNECTED  MessageType = "user_disconnected"
	CONFIRM_CONNECTION MessageType = "confirm_connection"
	SUBMIT_ESTIMATE    MessageType = "submit_estimate"
	ESTIMATE_REVEALED  MessageType = "estimate_revealed"
	RESET_PLANNING     MessageType = "reset_planning"
)

type MessageType string

type User struct {
	UserName string `json:"userName"`
	Uuid     string `json:"uuid"`
}


type ReceiveMessage struct {
	Type    MessageType
	Payload json.RawMessage
}

type SendMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

type JoinRoomMessage struct {
	RoomCode string `json:"roomCode"`
	Nickname string `json:"nickname"`
}

type UserJoinedMessage struct {
	User User `json:"user"`
}

type DisconnectionMessage struct {
	User User `json:"user"`
}

type ConfirmConnectionMessage struct {
	User User `json:"user"`
	ConnectedUsers []User `json:"ConnectedUsers"`
}

type SubmitEstimateMessage struct {
	TaskId   string `json:"taskId"`
	Estimate int    `json:"estimate"`
}

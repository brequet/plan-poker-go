package main

import (
	"encoding/json"
)

const (
	JOIN_ROOM         MessageType = "join_room"
	USER_JOINED       MessageType = "user_joined"
	SUBMIT_ESTIMATE   MessageType = "submit_estimate"
	ESTIMATE_REVEALED MessageType = "estimate_revealed"
	RESET_PLANNING    MessageType = "reset_planning"
)

type MessageType string

type Message struct {
	Type    MessageType
	Payload json.RawMessage
}

type JoinRoomMessage struct {
	RoomCode string `json:"roomCode"`
	Nickname string `json:"nickname"`
}

type UserJoinedMessage struct {
	UserName string `json:"userName"`
}

type SubmitEstimateMessage struct {
	TaskId   string `json:"taskId"`
	Estimate int    `json:"estimate"`
}

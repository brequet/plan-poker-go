package roomsmanager

import (
	"testing"
)

func TestCreateRoom(t *testing.T) {
	roomName := "TestRoom"
	createdRoom := CreateRoom(roomName)

	// Check if the created room is not nil
	if createdRoom == nil {
		t.Errorf("Expected a non-nil room, but got nil")
	}

	// Check if the room exists in the rooms map
	if _, exist := rooms[createdRoom.code]; !exist {
		t.Errorf("Expected the created room to exist in the rooms map, but it does not")
	}

	// Check if the room's name matches the provided name
	if createdRoom.name != roomName {
		t.Errorf("Expected room name %s, but got %s", roomName, createdRoom.name)
	}
}

func TestFindRoomByRoomCode(t *testing.T) {
	roomName := "TestRoom"
	createdRoom := CreateRoom(roomName)

	foundRoom := FindRoomByRoomCode(createdRoom.code)

	// Check if the found room is not nil
	if foundRoom == nil {
		t.Errorf("Expected a non-nil room, but got nil")
	}

	// Check if the found room's code matches the provided code
	if foundRoom.code != createdRoom.code {
		t.Errorf("Expected room code %s, but got %s", createdRoom.code, foundRoom.code)
	}

	// Check finding a non-existent room
	nonExistentRoomCode := "NonExistentRoom"
	notFoundRoom := FindRoomByRoomCode(nonExistentRoomCode)

	// Check if the not found room is nil
	if notFoundRoom != nil {
		t.Errorf("Expected a nil room for a non-existent room, but got a non-nil room")
	}
}

func TestConnectNewUserToRoom(t *testing.T) {
	roomName := "TestRoom"
	createdRoom := CreateRoom(roomName)

	nickname := "TestUser"
	user := ConnectNewUserToRoom(nickname, createdRoom.code)

	// Check if the user connected successfully
	if user == nil {
		t.Errorf("Expected the user to connect to the room, but it didn't happen")
	}

	// Check if the user exists in the room's users map
	if _, exist := createdRoom.users[user]; !exist {
		t.Errorf("Expected the user to exist in the room's users map, but it doesn't")
	}

	// Check connecting to a non-existent room
	nonExistentRoomCode := "NonExistentRoom"
	userConnectedToNonExistentRoom := ConnectNewUserToRoom(nickname, nonExistentRoomCode)

	// Check if the user failed to connect to the non-existent room
	if userConnectedToNonExistentRoom != nil {
		t.Errorf("Expected the user to fail connecting to a non-existent room, but it didn't happen")
	}
}

func TestDisconnectUserFromRoom(t *testing.T) {
	roomName := "TestRoom"
	createdRoom := CreateRoom(roomName)

	nickname := "TestUser"
	ConnectNewUserToRoom(nickname, createdRoom.code)

	user := &User{nickname: nickname}
	userDisconnected := DisconnectUserFromRoom(user, createdRoom.code)

	// Check if the user disconnected successfully
	if !userDisconnected {
		t.Errorf("Expected the user to disconnect from the room, but it didn't happen")
	}

	// Check if the user doesn't exist in the room's users map after disconnection
	if _, exist := createdRoom.users[user]; exist {
		t.Errorf("Expected the user to be removed from the room's users map, but it still exists")
	}

	// Check disconnecting from a non-existent room
	nonExistentRoomCode := "NonExistentRoom"
	userDisconnectedFromNonExistentRoom := DisconnectUserFromRoom(user, nonExistentRoomCode)

	// Check if the user failed to disconnect from the non-existent room
	if userDisconnectedFromNonExistentRoom {
		t.Errorf("Expected the user to fail disconnecting from a non-existent room, but it didn't happen")
	}
}

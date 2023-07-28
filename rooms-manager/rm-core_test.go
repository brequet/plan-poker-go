package roomsmanager

import (
	"testing"
)

// TestCreateRoom tests the CreateRoom function
func TestCreateRoom(t *testing.T) {
	// Test creating a new room
	roomName := "TestRoom"
	roomCode := CreateRoom(roomName)

	// Check if the room code is not empty
	if roomCode == "" {
		t.Errorf("Expected a non-empty room code, got an empty string")
	}

	// Check if the room with the same name exists
	existingRoomCode := CreateRoom(roomName)
	if existingRoomCode != roomCode {
		t.Errorf("Expected the existing room code to match the previously created room code")
	}
}

// TestFindRoomCodeByName tests the findRoomCodeByName function
func TestFindRoomCodeByName(t *testing.T) {
	// Create a room and get its code
	roomName := "TestRoom"
	roomCode := CreateRoom(roomName)

	// Test finding the room code by its name
	foundRoomCode := findRoomCodeByName(roomName)

	// Check if the found room code matches the expected room code
	if foundRoomCode != roomCode {
		t.Errorf("Expected room code %s, but found %s", roomCode, foundRoomCode)
	}

	// Test finding a non-existent room code
	nonExistentRoomName := "NonExistentRoom"
	nonExistentRoomCode := findRoomCodeByName(nonExistentRoomName)

	// Check if the non-existent room code is an empty string
	if nonExistentRoomCode != "" {
		t.Errorf("Expected an empty room code for a non-existent room, but found %s", nonExistentRoomCode)
	}
}

// TestGenerateRoomCode tests the generateRoomCode function
func TestGenerateRoomCode(t *testing.T) {
	// Generate a room code
	roomCode := generateRoomCode()

	// Check if the generated room code has the correct length
	if len(roomCode) != ROOM_CODE_SIZE {
		t.Errorf("Expected room code length of %d, but found %d", ROOM_CODE_SIZE, len(roomCode))
	}
}

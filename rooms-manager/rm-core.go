package roomsmanager

import (
	"log"
	"math/rand"
	"time"

	// "github.com/google/uuid"
)

const ROOM_CODE_SIZE = 4

type RoomCode = string

type client struct {
	uuid int
	nickname string
}

type room struct {
	clients map[*client]bool
	name string
	code RoomCode // must be unique, can act as an ID
}

var rooms = make(map[string]*room)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Hello() {
	log.Println("HELLO")

}

func CreateRoom(roomName string) RoomCode {
	// Get the room for the given name or create a new one if it doesn't exist
	roomCode := findRoomCodeByName(roomName) // TODO: always create, can have same name
	// Geeet tee rrromm for thhe given name  orr createe a new one if it does't exist
	if roomCode == "" {
		roomCode = generateRoomCode()
		log.Printf("No room found for name %s, creating a new room with code : %s", roomName, roomCode)
		rooms[roomName] = &room{
			clients: make(map[*client]bool),
			name: roomName,
			code: roomCode,
		}
	}
	return roomCode
}

func findRoomCodeByName(roomName string) RoomCode {
	for roomCode, room := range rooms {
		if room.name == roomName {
			return roomCode
		}
	}
	return ""
}

func generateRoomCode() RoomCode {
	b := make([]rune, ROOM_CODE_SIZE)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
	return string(b)
}

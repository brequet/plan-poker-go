package roomsmanager

import (
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

const ROOM_CODE_SIZE = 4

type User struct {
	uuid     string
	nickname string
}

type Room struct {
	users map[*User]bool
	Name  string
	Code  string // must be unique, can act as an ID
}

var rooms = make(map[string]*Room)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateRoom(roomName string) *Room {
	roomCode := generateUniqueRoomCode()
	log.Printf("Creating room with name '%s', generated code : '%s'", roomName, roomCode)
	createdRoom := &Room{
		users: make(map[*User]bool),
		Name:  roomName,
		Code:  roomCode,
	}
	rooms[roomCode] = createdRoom
	return createdRoom
	// TODO: timeout 5min no one connected -> delete room
}

func FindRoomByRoomCode(roomCode string) *Room {
	room, exist := rooms[roomCode]
	if !exist {
		log.Printf("Room with code '%s' not found", roomCode)
		return nil
	}
	return room
}

func ConnectNewUserToRoom(nickname, roomCode string) *User { // TODO: return errors saying why user cant connect (ex: room doesn't exist)
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		return nil
	}

	newUser := createUser(nickname)
	room.users[newUser] = true
	log.Printf("User '%s' joined the room : '%s'", nickname, roomCode)

	return newUser
}

func DisconnectUserFromRoom(user *User, roomCode string) bool {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		return false
	}

	log.Printf("Removing user with nickname '%s', from room : '%s'", user.nickname, roomCode)
	delete(room.users, user)
	// TODO: timeout if no user left -> delete room
	return true
}

func createUser(nickname string) *User {
	// Get the string representation of the UUID
	uuidStr := uuid.New().String()
	log.Printf("Creating user with nickname '%s', uuid : '%s'", nickname, uuidStr)
	return &User{
		nickname: nickname,
		uuid:     uuidStr,
	}
}

func generateUniqueRoomCode() string {
	var generatedRoomCode string
	for {
		generatedRoomCode = generateRoomCode()
		if _, exist := rooms[generatedRoomCode]; !exist {
			break // Exit the loop if the generated code is unique
		}
	}
	return generatedRoomCode
}

func generateRoomCode() string {
	b := make([]rune, ROOM_CODE_SIZE)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

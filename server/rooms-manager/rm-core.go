package roomsmanager

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const ROOM_CODE_SIZE = 4

type User struct {
	Uuid     string
	Nickname string
	Estimate string
}

type Room struct {
	Users              map[*User]bool
	Name               string
	Code               string // must be unique, can act as an ID
	IsEstimateRevealed bool
	lastUpdated        time.Time
}

var rooms = make(map[string]*Room)

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
	startRoomCleanup()
}

func CreateRoom(roomName string) *Room {
	if roomName == "" {
		return nil
	}
	roomCode := generateUniqueRoomCode()
	log.Printf("[%s] Creating room with name '%s'", roomCode, roomName)
	createdRoom := &Room{
		Users:              make(map[*User]bool),
		Name:               roomName,
		Code:               roomCode,
		IsEstimateRevealed: false,
		lastUpdated:        time.Now(),
	}
	rooms[roomCode] = createdRoom
	return createdRoom
}

func FindRoomByRoomCode(roomCode string) *Room {
	room, exist := rooms[roomCode]
	if !exist {
		log.Printf("Room not found from code [%s]", roomCode)
		return nil
	}
	return room
}

func GetAllUserFromRoomByRoomCode(roomCode string) (users []*User) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		return users
	}

	for user, ok := range room.Users {
		if ok {
			users = append(users, user)
		}
	}
	return users
}

func ConnectNewUserToRoom(nickname, roomCode string) *User {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		log.Printf("User '%s' could not join the room: the room [%s] was not found.", nickname, roomCode)
		return nil
	}

	newUser := createUser(nickname)
	room.Users[newUser] = true
	room.lastUpdated = time.Now()
	log.Printf("[%s] User %s joined the room", roomCode, newUser.Uuid)

	return newUser
}

func DisconnectUserFromRoom(user *User, roomCode string) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		return
	}

	room.lastUpdated = time.Now()

	log.Printf("[%s] User '%s' disconnected", roomCode, user.Nickname)
	delete(room.Users, user)
}

func SubmitEstimate(user *User, roomCode string, estimate string) (err error) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		errMsg := fmt.Sprintf("User %s could not submit estimate {%s}: the room [%s] was not found", user.Uuid, estimate, roomCode)
		return errors.New(errMsg)
	}

	room.lastUpdated = time.Now()

	log.Printf("[%s] User %s submitted estimate {%s}", roomCode, user.Uuid, estimate)
	user.Estimate = estimate

	return nil
}

func ToggleShouldRevealEstimateForRoom(roomCode string) (bool, error) {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		errMsg := fmt.Sprintf("Cannot toggle estimate reveal: the room [%s] was not found", roomCode)
		return false, errors.New(errMsg)
	}

	room.lastUpdated = time.Now()
	room.IsEstimateRevealed = !room.IsEstimateRevealed
	if room.IsEstimateRevealed {
		log.Printf("[%s] Estimates revealed", roomCode)
	} else {
		log.Printf("[%s] Estimates hidden", roomCode)
	}

	return room.IsEstimateRevealed, nil
}

func ResetPlanningForRoom(roomCode string) error {
	room := FindRoomByRoomCode(roomCode)
	if room == nil {
		errMsg := fmt.Sprintf("Cannot reset the planning: the room [%s] was not found", roomCode)
		return errors.New(errMsg)
	}

	room.lastUpdated = time.Now()
	room.IsEstimateRevealed = false
	for user := range room.Users {
		user.Estimate = ""
	}

	log.Printf("[%s] The estimates were reset", roomCode)

	return nil
}

func createUser(nickname string) *User {
	// Get the string representation of the UUID
	uuidStr := uuid.New().String()
	log.Printf("Creating user with nickname '%s', uuid : '%s'", nickname, uuidStr)
	return &User{
		Nickname: nickname,
		Uuid:     uuidStr,
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

func startRoomCleanup() {
	// Define a ticker to run the cleanup at a regular interval
	ticker := time.NewTicker(time.Minute * 1) // Adjust the interval as needed

	// Run the cleanup routine in a goroutine
	go func() {
		for range ticker.C {
			cleanupEmptyRooms()
		}
	}()
}

func cleanupEmptyRooms() {
	currentTimestamp := time.Now()
	maxIdleDuration := time.Minute * 5
	// Iterate over the rooms and remove those that are empty and idle for too long
	for roomCode, room := range rooms {
		if len(room.Users) == 0 && currentTimestamp.Sub(room.lastUpdated) > maxIdleDuration {
			delete(rooms, roomCode)
			log.Printf("[%s] Room deleted due to inactivity", roomCode)
		}
	}
}

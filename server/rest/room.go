package rest

import (
	"encoding/json"
	"log"
	"net/http"

	rm "github.com/baptiste-requet/plan-poker-go/rooms-manager"
	"github.com/gorilla/mux"
)

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get room name and other data
	var roomData struct {
		RoomName string `json:"roomName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&roomData); err != nil {
		log.Printf("Error in createRoomHandler: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if roomData.RoomName == "" {
		// RoomName empty, return an error response
		http.NotFound(w, r) // TODO: bad param
		return
	}

	// Create a new room with the provided room name
	newRoom := rm.CreateRoom(roomData.RoomName)

	// Respond with the room details (e.g., room ID) to the frontend
	response := mapRmRoomToRoom(newRoom)

	// Send the response back to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the room code from the URL path parameter
	roomCode := mux.Vars(r)["roomCode"]

	// Retrieve the room from your data store using the room code
	room := rm.FindRoomByRoomCode(roomCode)

	if room == nil {
		http.NotFound(w, r)
		return
	}

	// Respond with the room details (e.g., room ID) to the frontend
	response := mapRmRoomToRoom(room)

	// Send the response back to the frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func mapRmRoomToRoom(rmRoom *rm.Room) *Room {
	return &Room{
		RoomCode: rmRoom.Code,
		RoomName: rmRoom.Name,
	}
}

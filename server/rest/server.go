package rest

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// TODO: better logging (with levels)

func Start() error {
	addrAndPort := "0.0.0.0:8080"

	router := mux.NewRouter()

	// General endpoints
	router.HandleFunc("/api/health", healthCheckHandler).Methods("GET")

	// Handled in websocket.go
	router.HandleFunc("/api/ws", wsHandler)

	// Handled in room.go
	router.HandleFunc("/api/room", createRoomHandler).Methods("POST")
	router.HandleFunc("/api/room/{roomCode}", getRoomHandler).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	// originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")}) TODO
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("Server started on ", addrAndPort)
	err := http.ListenAndServe(addrAndPort, handlers.CORS(originsOk, headersOk, methodsOk)(router))

	return err
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

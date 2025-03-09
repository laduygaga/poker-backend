package handlers

import (
	"log"
	"net/http"
	"poker-backend/internal/services"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	playerID := r.URL.Query().Get("player_id")
	go services.HandleConnection(conn, playerID)
	go services.HandleBroadcast()
}

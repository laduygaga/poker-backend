package services

import (
	"context"
	"encoding/json"
	"log"
	"poker-backend/internal/models"
	"poker-backend/pkg/database"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string) // Map client -> playerID
var broadcast = make(chan models.GameState)

func HandleConnection(conn *websocket.Conn, playerID string) {
    log.Printf("Player connected: %s", playerID)
    clients[conn] = playerID
    defer func() {
        log.Printf("Player disconnected: %s", playerID)
        delete(clients, conn)
        conn.Close()
    }()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading message:", err)
            break
        }

        var gameState models.GameState
        json.Unmarshal(msg, &gameState)
        broadcast <- gameState
    }
}

func deepCopyGameState(gameState models.GameState) models.GameState {
    data, _ := json.Marshal(gameState)
    var copiedState models.GameState
    json.Unmarshal(data, &copiedState)
    return copiedState
}

func HandleBroadcast() {
    for {
        gameState := <-broadcast
        for client, playerID := range clients {
			playerState := deepCopyGameState(gameState)
            log.Printf("Broadcasting to player %s", playerID)

            for i := range playerState.Players {
                if playerState.Players[i].ID != playerID {
                    log.Printf("Hiding hole cards for player %s", playerState.Players[i].ID)
                    playerState.Players[i].HoleCards = []string{}
                } else {
                    log.Printf("Keeping hole cards for player %s: %v", playerID, playerState.Players[i].HoleCards)
                }
            }

            data, err := json.Marshal(playerState)
            if err != nil {
                log.Println("Error marshaling player state:", err)
                continue
            }

            err = client.WriteMessage(websocket.TextMessage, data)
            if err != nil {
                log.Printf("Error broadcasting to player %s: %v", playerID, err)
                client.Close()
                delete(clients, client)
            } else {
                log.Printf("Successfully sent state to player %s", playerID)
            }
        }
    }
}

func BroadcastGameState(tableID string) error {
    ctx := context.Background()
    data, err := database.GetRedisClient().Get(ctx, "table:"+tableID).Result()
    if err != nil {
        return err
    }

    var gameState models.GameState
    if err := json.Unmarshal([]byte(data), &gameState); err != nil {
        return err
    }

    // Gửi trạng thái game qua channel broadcast
    broadcast <- gameState
    return nil
}

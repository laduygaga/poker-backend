package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"poker-backend/internal/models"
	"poker-backend/internal/utils"
	"poker-backend/pkg/database"
)

func CreateTable(tableID string, playerID string) error {
	deck := utils.NewDeck()
	deck.Shuffle()
	gameState := models.GameState{
		TableID: tableID,
		Players: []models.Player{{ID: playerID, Chips: 1000}},
		Deck:    deck.Cards,
		Pot:     0,
		CurrentTurn: playerID,
	}

	data, _ := json.Marshal(gameState)
	return database.GetRedisClient().Set(context.TODO(), "table:"+tableID, data, 0).Err()
}

func JoinTable(tableID, playerID string) error {
	data, err := database.GetRedisClient().Get(context.TODO(), "table:"+tableID).Result()
	if err != nil {
		return err
	}

	var gameState models.GameState
	json.Unmarshal([]byte(data), &gameState)
	gameState.Players = append(gameState.Players, models.Player{ID: playerID, Chips: 1000})
	updatedData, _ := json.Marshal(gameState)
	return database.GetRedisClient().Set(context.TODO(), "table:"+tableID, updatedData, 0).Err()
}

func DealCards(tableID string) error {
    ctx := context.Background()
    data, err := database.GetRedisClient().Get(ctx, "table:"+tableID).Result()
    if err != nil {
        return err
    }

    var gameState models.GameState
    if err := json.Unmarshal([]byte(data), &gameState); err != nil {
        return err
    }

    deck := &utils.Deck{Cards: gameState.Deck}
    for i := range gameState.Players {
        gameState.Players[i].HoleCards = deck.Deal(2)
        log.Printf("Dealt cards to %s: %v", gameState.Players[i].ID, gameState.Players[i].HoleCards)
    }
    gameState.Deck = deck.Cards

    updatedData, err := json.Marshal(gameState)
    if err != nil {
        return err
    }
    err = database.GetRedisClient().Set(ctx, "table:"+tableID, updatedData, 0).Err()
    if err != nil {
        return err
    }

    return BroadcastGameState(tableID)
}

func HandlePlayerAction(tableID, playerID, action string, amount int) error {
    ctx := context.Background()
    data, err := database.GetRedisClient().Get(ctx, "table:"+tableID).Result()
    if err != nil {
        return err
    }

    var gameState models.GameState
    json.Unmarshal([]byte(data), &gameState)

    // Tìm người chơi hiện tại
    var currentPlayer *models.Player
    for i := range gameState.Players {
        if gameState.Players[i].ID == playerID {
            currentPlayer = &gameState.Players[i]
            break
        }
    }
    if currentPlayer == nil {
        return errors.New("player not found")
    }

    // In hole_cards của người chơi hiện tại trước khi thực hiện hành động
    log.Printf("Before action: Player %s hole cards: %v", playerID, currentPlayer.HoleCards)

    // Kiểm tra lượt chơi
    if gameState.CurrentTurn != playerID {
        return errors.New("not your turn")
    }

    // Xử lý hành động
    switch action {
    case "check":
    case "call":
    case "raise":
        if amount <= 0 || amount > currentPlayer.Chips {
            return errors.New("invalid raise amount")
        }
        currentPlayer.Chips -= amount
        gameState.Pot += amount
    case "fold":
        currentPlayer.IsFolded = true
    default:
        return errors.New("invalid action")
    }

    // Chuyển lượt chơi
    gameState.CurrentTurn = getNextPlayer(gameState)

    // In hole_cards của người chơi hiện tại sau khi thực hiện hành động
    log.Printf("After action: Player %s hole cards: %v", playerID, currentPlayer.HoleCards)

    // Lưu trạng thái mới
    updatedData, err := json.Marshal(gameState)
    if err != nil {
        return err
    }
    err = database.GetRedisClient().Set(ctx, "table:"+tableID, updatedData, 0).Err()
    if err != nil {
        return err
    }

    // Gửi trạng thái game qua WebSocket
    return BroadcastGameState(tableID)
}

func getNextPlayer(gameState models.GameState) string {
    currentIndex := -1
    for i, player := range gameState.Players {
        if player.ID == gameState.CurrentTurn {
            currentIndex = i
            break
        }
    }

    for i := 1; i <= len(gameState.Players); i++ {
        nextIndex := (currentIndex + i) % len(gameState.Players)
        if !gameState.Players[nextIndex].IsFolded {
            return gameState.Players[nextIndex].ID
        }
    }
    return gameState.CurrentTurn // Nếu không tìm thấy người chơi tiếp theo
}

func DealCommunityCards(tableID, stage string) error {
    ctx := context.Background()
    data, err := database.GetRedisClient().Get(ctx, "table:"+tableID).Result()
    if err != nil {
        return err
    }

    var gameState models.GameState
    json.Unmarshal([]byte(data), &gameState)
    deck := &utils.Deck{Cards: gameState.Deck}

    switch stage {
    case "flop":
        gameState.CommunityCards = deck.Deal(3)
    case "turn":
        gameState.CommunityCards = append(gameState.CommunityCards, deck.Deal(1)...)
    case "river":
        gameState.CommunityCards = append(gameState.CommunityCards, deck.Deal(1)...)
    default:
        return errors.New("invalid stage")
    }

    gameState.Deck = deck.Cards
    updatedData, _ := json.Marshal(gameState)
    err = database.GetRedisClient().Set(ctx, "table:"+tableID, updatedData, 0).Err()
    if err != nil {
        return err
    }

    return BroadcastGameState(tableID)
}

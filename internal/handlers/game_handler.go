package handlers

import (
	"encoding/json"
	"net/http"
	"poker-backend/internal/services"
)

func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TableID  string `json:"table_id"`
		PlayerID string `json:"player_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	err := services.CreateTable(req.TableID, req.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func JoinTableHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TableID  string `json:"table_id"`
		PlayerID string `json:"player_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	err := services.JoinTable(req.TableID, req.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DealCardsHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        TableID string `json:"table_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err := services.DealCards(req.TableID)
    if err != nil {
        http.Error(w, "Failed to deal cards: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Cards dealt successfully"})
}

func PlayerActionHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        TableID  string `json:"table_id"`
        PlayerID string `json:"player_id"`
        Action   string `json:"action"` // "check", "call", "raise", "fold"
        Amount   int    `json:"amount"` // Số tiền cược (nếu raise)
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err := services.HandlePlayerAction(req.TableID, req.PlayerID, req.Action, req.Amount)
    if err != nil {
        http.Error(w, "Failed to handle action: "+err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Action processed successfully"})
}

func DealCommunityCardsHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        TableID string `json:"table_id"`
        Stage   string `json:"stage"` // "flop", "turn", "river"
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err := services.DealCommunityCards(req.TableID, req.Stage)
    if err != nil {
        http.Error(w, "Failed to deal community cards: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Community cards dealt successfully"})
}

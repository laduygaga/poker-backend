package main

import (
	"log"
	"net/http"
	"poker-backend/internal/config"
	"poker-backend/internal/middleware" 
	"poker-backend/internal/handlers"
	"poker-backend/pkg/database"
	"github.com/gorilla/mux"
)

func main() {
	// Khởi tạo config
	cfg := config.LoadConfig()

	// Kết nối MongoDB và Redis
	database.InitMongoDB(cfg.MongoURI)
	database.InitRedis(cfg.RedisAddr)

	// Khởi tạo router chính
	r := mux.NewRouter()

	// Public routes (không cần token)
	r.HandleFunc("/api/auth/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", handlers.LoginHandler).Methods("POST")

	// Protected routes (yêu cầu token)
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware) // Sử dụng middleware.AuthMiddleware
	protected.HandleFunc("/game/create-table", handlers.CreateTableHandler).Methods("POST")
	protected.HandleFunc("/game/join-table", handlers.JoinTableHandler).Methods("POST")
	protected.HandleFunc("/game/deal-cards", handlers.DealCardsHandler).Methods("POST")
	protected.HandleFunc("/game/action", handlers.PlayerActionHandler).Methods("POST")
	protected.HandleFunc("/game/deal-community-cards", handlers.DealCommunityCardsHandler).Methods("POST")

	// WebSocket endpoint
	r.HandleFunc("/ws", handlers.SocketHandler)

	log.Printf("Server started at :%s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}

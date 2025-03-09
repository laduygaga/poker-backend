package services

import (
	"context"
	"poker-backend/internal/models"
	"poker-backend/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func Register(username, password string) error {
	collection := database.GetMongoClient().Database("poker").Collection("users")
	user := models.User{
		Username: username,
		Password: password, // Trong thực tế, cần hash password
		Chips:    1000,
	}
	_, err := collection.InsertOne(context.TODO(), user)
	return err
}

func Login(username, password string) (string, error) {
	collection := database.GetMongoClient().Database("poker").Collection("users")
	var user models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil || user.Password != password { // Trong thực tế, cần so sánh password đã hash
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	return tokenString, err
}

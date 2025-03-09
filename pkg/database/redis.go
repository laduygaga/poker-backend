package database

import (
	"context"
	"log"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis(addr string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to Redis")
}

func GetRedisClient() *redis.Client {
	return redisClient
}

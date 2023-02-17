package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func InitRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:   0,
	})

	// checking if redis is connected
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis Connection Failed",
			err)
	}

	log.Println("Redis Successfully Connected.",
		"Ping", pong)

	return redisClient
}

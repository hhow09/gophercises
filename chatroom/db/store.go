package db

import (
	"github.com/go-redis/redis/v8"
)

type Store struct {
	redis *redis.Client
}

func NewStore() *Store {
	return &Store{redis: InitRedisClient()}
}

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

func (s *Store) CreateUser(ctx context.Context, params CreateUserParams) error {
	userstr, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
		return err
	}
	status := s.redis.Set(context.Background(), USER_KEY, userstr, 0)
	fmt.Println(status)
	return nil
}

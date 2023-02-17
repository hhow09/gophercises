package main

import (
	"log"

	"github.com/hhow09/gophercises/chatroom/db"
	"github.com/hhow09/gophercises/chatroom/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	store := db.NewStore()
	server, err := server.NewServer(store)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = server.Serve()
	if err != nil {
		log.Fatal(err)
		return
	}
}

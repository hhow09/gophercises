package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hhow09/gophercises/chatroom/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	_, err := server.InitServer()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("visit http://localhost:%v/chatroom\n", os.Getenv("WEB_HOST"))
}

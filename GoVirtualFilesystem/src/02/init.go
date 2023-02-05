package main

import (
	"fmt"
	"log"

	"github.com/chzyer/readline"
)

// setName gets a custom username from the current user.
func setName() string {
	rl, err := readline.New(">")
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()
	rl.CaptureExitSignal()
	var username string
	for {
		fmt.Println("enter user name")
		input, err := rl.Readline()
		if err != nil {
			log.Fatal(err)
		}
		if input == "1" {
			username = "Guest"
			break
		}
		if input != "" {
			username = input
			break
		}
	}
	return username
}

// initUser initializes the user object on startup.
func initUser() *user {
	username := setName()
	currentUser := createUser(username)
	return currentUser
}

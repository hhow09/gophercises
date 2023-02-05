package main

import (
	"log"
	"strings"
)

// shellLoop runs the main shell loop for the filesystem.
func shellLoop(currentUser *user) {
	filesystem := initFilesystem()
	prompt := currentUser.initPrompt()
	for {
		input, err := prompt.Readline()
		if err != nil {
			log.Fatal(err)
			return
		}
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			continue
		}
		filesystem.execute(input)
	}
}

func main() {
	currentUser := initUser()
	shellLoop(currentUser)
}

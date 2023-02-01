package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		if string(message) == "close" {
			ws.Close()
		}
		if string(message) == "ping" {
			message = []byte("pong")
		} else {
			message = []byte(fmt.Sprintln("echo: ", string(message)))
		}
		if err = ws.WriteMessage(mt, message); err != nil {
			break
		}
	}
}

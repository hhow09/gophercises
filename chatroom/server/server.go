package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"os"

	socketio "github.com/googollee/go-socket.io"
)

func InitServer() (*socketio.Server, error) {
	router := gin.New()
	router.LoadHTMLGlob("./static/*")

	server := socketio.NewServer(nil)

	_, err := server.Adapter(&socketio.RedisAdapterOptions{
		Addr:   fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:     0,
		Prefix: "socket.io",
	})
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("user connected ID:", s.ID())
		ok := server.JoinRoom("/", "bcast", s)
		fmt.Println("join room bcast", ok)

		go func(s socketio.Conn) {
			server.BroadcastToRoom("/", "bcast", "connected", s.ID())
		}(s)
		return nil
	})

	// set
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		server.BroadcastToRoom("/", "bcast", "said", s.ID(), msg) // broadcast to all users
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed user ID:", s.ID(), reason)
	})

	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.GET("/chatroom", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "chat.tmpl", gin.H{"Debug": true})
	})
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()
	if err := router.Run(fmt.Sprintf(":%v", os.Getenv("WEB_HOST"))); err != nil {
		log.Fatal("failed run app: ", err)
		return nil, err
	}
	return server, nil
}

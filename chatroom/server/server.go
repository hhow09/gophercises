package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hhow09/gophercises/chatroom/db"

	"os"

	socketio "github.com/googollee/go-socket.io"
)

const (
	MAIN_NS        = "/"
	PUBLIC_CHAT_NS = "/public_chat"
)

type Server struct {
	store        *db.Store
	router       *gin.Engine
	socketServer *socketio.Server
}

func NewServer(store *db.Store) (*Server, error) {
	socketServer := socketio.NewServer(nil)
	_, err := socketServer.Adapter(&socketio.RedisAdapterOptions{
		Addr:   fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:     0,
		Prefix: "socket.io",
	})
	if err != nil {
		log.Fatal("init socket server error:", err)
		return nil, err
	}
	server := &Server{store: store, router: gin.New(), socketServer: socketServer}
	err = server.setupSocketRouter()
	if err != nil {
		log.Fatal("init socket server error:", err)
		return nil, err
	}
	server.setupHTTPRouter()

	return server, nil
}

func (s *Server) setupSocketRouter() error {
	s.socketServer.OnConnect(MAIN_NS, func(sc socketio.Conn) error {
		log.Println("user connected ID:", sc.ID())
		ok := s.socketServer.JoinRoom(MAIN_NS, "bcast", sc)
		fmt.Println("join room bcast", ok)

		go func(sc socketio.Conn) {
			s.socketServer.BroadcastToRoom(MAIN_NS, "bcast", "connected", sc.ID())
		}(sc)
		return nil
	})

	s.socketServer.OnEvent(PUBLIC_CHAT_NS, "msg", func(sc socketio.Conn, msg string) string {
		s.socketServer.BroadcastToRoom(MAIN_NS, "bcast", "said", sc.ID(), msg) // broadcast to all users
		return "recv " + msg
	})

	s.socketServer.OnEvent(MAIN_NS, "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	s.socketServer.OnError(MAIN_NS, func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	s.socketServer.OnDisconnect(MAIN_NS, func(s socketio.Conn, reason string) {
		log.Println("closed user ID:", s.ID(), reason)
	})

	return nil
}

func (s *Server) setupHTTPRouter() {
	s.router.LoadHTMLGlob("./static/*")
	s.router.GET("/socket.io/*any", AuthMiddleware(), gin.WrapH(s.socketServer))
	s.router.POST("/socket.io/*any", AuthMiddleware(), gin.WrapH(s.socketServer))
	s.router.GET("/public_chat", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "public_chat.tmpl", gin.H{"Debug": true})
	})
}

func (s *Server) Serve() error {
	go func() {
		if err := s.socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		} else {
			fmt.Println("socketio success listen")
		}
	}()
	defer s.socketServer.Close()
	if err := s.router.Run(fmt.Sprintf(":%v", os.Getenv("WEB_HOST"))); err != nil {
		log.Fatal("failed run app: ", err)
		return err
	}
	fmt.Printf("visit http://localhost:%v/public_chat\n", os.Getenv("WEB_HOST"))
	return nil
}

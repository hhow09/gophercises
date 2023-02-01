package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hhow09/gophercises/ws_demo/handler"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func main() {
	flag.Parse()
	http.HandleFunc("/ping", handler.WsHandler)
	http.HandleFunc("/", serveHome)
	fmt.Printf("visit: http://localhost%s/ \n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

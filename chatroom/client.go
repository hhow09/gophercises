// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/queue-b/gocketio"
)

func init() {
	godotenv.Load()
}

func main() {
	flag.Parse()
	log.SetFlags(1)

	u := url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%s", os.Getenv("WEB_HOST"))}
	manager, err := gocketio.DialContext(context.Background(), u.String(), gocketio.DefaultManagerConfig())
	mainClient, err := manager.Namespace("/")
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}
	client, err := manager.Namespace("/chat")

	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	mainClient.On("connected", func(msg string) {
		log.Printf("user-%v connected\n", msg)
	})
	mainClient.On("said", func(id, msg string) {
		log.Printf("user-%v said:%v\n", id, msg)
	})
	mainClient.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		reader.Discard(reader.Buffered())
		command := string(data)
		err := client.Emit("msg", command)

		if err != nil {
			log.Printf("err %v", err)
		}
	}
}

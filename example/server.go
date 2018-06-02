package main

import (
	"log"
	"time"

	"github.com/ChaimHong/kfkrpc"
	"github.com/ChaimHong/kfkrpc/example/service1"
)

var generate bool

func Server(sid uint16) {
	server := kfkrpc.NewServer(sid)
	server.Register(&service1.ServiceA{})

	log.Println("start server...")
	go server.Serve(newSaramaClient())

	for {
		time.Sleep(10 * time.Second)
	}
}

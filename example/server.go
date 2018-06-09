package main

import (
	"log"

	"github.com/ChaimHong/kfkrpc"
	"github.com/ChaimHong/kfkrpc/example/service1"
)

func Server(sid uint16, addr string) {
	server := kfkrpc.NewServer(sid)
	server.Register(&service1.ServiceA{})

	go server.Serve(addr)

	log.Println("start server...", sid, addr)
}

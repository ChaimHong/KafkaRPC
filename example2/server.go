package main

import (
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/ChaimHong/util"

	"github.com/ChaimHong/kfkrpc/example/service1"
)

var generate bool

func Server() {
	server := rpc.NewServer()
	server.Register(&service1.ServiceA{})

	log.Println("start server...")
	lis, e := net.Listen("tcp", "0.0.0.0:20001")
	util.CheckPanic(e)
	go server.Accept(lis)

	for {
		time.Sleep(10 * time.Second)
	}
}

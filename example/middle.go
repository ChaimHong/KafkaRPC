package main

import (
	"github.com/ChaimHong/kfkrpc"
)

func MiddleWare(addr string) {
	mw, err := kfkrpc.NewMiddleWare(addr)
	if err != nil {
		panic(err)
	}
	mw.Server()
}

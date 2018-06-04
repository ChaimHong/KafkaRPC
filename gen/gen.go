package main

import (
	"github.com/ChaimHong/kfkrpc"
	"github.com/funny/fastbin"
)

func main() {
	fastbin.Register(&kfkrpc.KFKMessage{})
	fastbin.Register(&kfkrpc.ResponeMsg{})
	fastbin.GenCode()
}

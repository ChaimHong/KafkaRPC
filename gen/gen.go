package main

import (
	"github.com/ChaimHong/kfkrpc"
	"github.com/funny/fastbin"
)

func main() {
	fastbin.Register(&kfkrpc.KFKMessage{})
	fastbin.Register(&kfkrpc.ResponeMsg{})
	fastbin.Register(&kfkrpc.MiddleReply{})
	fastbin.Register(&kfkrpc.RpcRequest{})
	fastbin.Register(&kfkrpc.RpcResponse{})

	// fastbin.Register(&service2.A{})
	fastbin.GenCode()
}

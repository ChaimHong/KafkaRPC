package main

import (
	"fmt"

	"github.com/ChaimHong/kfkrpc"
	"github.com/ChaimHong/kfkrpc/example/service1"
)

func Client(sid uint16) {
	client := kfkrpc.NewClient(newSaramaClient(), sid)
	done := make(chan bool, 1)

	reply := &service1.AOut{}
	client.Call(1, service1.ServiceA_A,
		&kfkrpc.Request{
			Args:  &service1.AIn{V: 2},
			Reply: reply,
		}, func(err error) {
			done <- true
		})
	<-done

	fmt.Println(reply)

	reply2 := ""
	var arg2 int = 2
	client.Call(1, service1.ServiceA_B,
		&kfkrpc.Request{
			Args:  &arg2,
			Reply: &reply2,
		}, func(err error) {
			done <- true
		})
	<-done

	fmt.Println(reply2)

}

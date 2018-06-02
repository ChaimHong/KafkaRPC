package main

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/ChaimHong/kfkrpc/example/service1"
	"github.com/ChaimHong/util"
)

func Client() {
	conn, e := net.Dial("tcp", "127.0.0.1:20001")
	util.CheckPanic(e)
	client := rpc.NewClient(conn)
	{
		reply := &service1.AOut{}
		err := client.Call(service1.ServiceA_A, &service1.AIn{V: 2}, reply)
		util.CheckPanic(err)
		fmt.Println(reply)
	}
	{
		reply := &service1.AOut{}
		err := client.Call(service1.ServiceA_A, &service1.AIn{V: 2}, reply)
		util.CheckPanic(err)
		fmt.Println(reply)
	}
	// {
	// 	reply := &service1.AOut{}
	// 	err := client.Call(service1.ServiceA_A, &service1.AIn{V: 2}, reply)
	// 	util.CheckPanic(err)
	// 	fmt.Println(reply)
	// }
	return
	{

		reply2 := ""
		args := int(2)
		e2 := client.Call(service1.ServiceA_B, &args, &reply2)
		util.CheckPanic(e2)
		fmt.Println(reply2)
	}
}

package main

import (
	"flag"

	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
	"github.com/funny/fastbin"

	"github.com/ChaimHong/kfkrpc"
	"github.com/ChaimHong/kfkrpc/example/service1"
)

func main() {
	isSrv := false
	gen := false
	var sid int
	flag.BoolVar(&isSrv, "is_server", true, "is server")
	flag.IntVar(&sid, "sid", 0, "sid")
	flag.BoolVar(&gen, "gen", false, "gen")

	flag.Parse()

	if gen {
		server := kfkrpc.NewServer(0)
		server.Register(&service1.ServiceA{})

		fastbin.GenCode()
		return
	}

	if isSrv {
		Server(uint16(sid))
	} else {
		Client(uint16(sid))
	}
}

func newSaramaClient() sarama.Client {
	config := sarama.NewConfig()
	addrs := []string{
		"localhost:9092",
	}
	gClient, e := sarama.NewClient(addrs, config)
	util.CheckPanic(e)

	return gClient
}

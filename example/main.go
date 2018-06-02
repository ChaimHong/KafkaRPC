package main

import (
	"flag"

	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
)

func main() {
	isSrv := false
	var sid int
	flag.BoolVar(&isSrv, "is_server", true, "is server")
	flag.IntVar(&sid, "sid", 0, "sid")

	flag.Parse()

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

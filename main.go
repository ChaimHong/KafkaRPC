package kfkrpc

import (
	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
)

var gClose = make(chan int, 1)

func newSaramaClient() sarama.Client {
	config := sarama.NewConfig()
	addrs := []string{
		"localhost:9092",
	}
	gClient, e := sarama.NewClient(addrs, config)
	util.CheckPanic(e)

	return gClient
}

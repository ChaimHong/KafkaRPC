package kfkrpc

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
)

var gClient sarama.Client
var gClose = make(chan int, 1)

func main() {
	getSaramaClient()

	client := gClient
	consumer, err := sarama.NewConsumerFromClient(client)
	util.CheckPanic(err)

	topics, err0 := consumer.Topics()
	util.CheckPanic(err0)
	util.Printfn(consumer, topics)
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	filename := "offset.txt"
	oLastOffset := int64(-2)
	if _, e := os.Stat(filename); e == nil {
		f0, e0 := os.Open(filename)
		util.CheckPanic(e0)
		d, e1 := ioutil.ReadAll(f0)
		util.CheckPanic(e1)
		f0.Close()

		v, e2 := strconv.Atoi(string(d))
		util.CheckPanic(e2)
		oLastOffset = int64(v)
		oLastOffset++
	}

	partitionConsumer, err := consumer.ConsumePartition("test", 0, oLastOffset)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d %s %s\n", msg.Offset, msg.Key, msg.Value)
			oLastOffset = msg.Offset
			consumed++
		case <-signals:
			close(gClose)
		case <-gClose:
			break ConsumerLoop
		}
	}

	f, err := os.Create("offset.txt")
	util.CheckPanic(err)
	f.Write([]byte(strconv.Itoa(int(oLastOffset))))
	f.Close()

	log.Printf("Consumed: %d\n", consumed)
}

func getSaramaClient() {
	if gClient == nil {
		config := sarama.NewConfig()
		addrs := []string{
			"localhost:9092",
		}
		var e error
		gClient, e = sarama.NewClient(addrs, config)
		util.CheckPanic(e)
	}
	return
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

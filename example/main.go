package main

import (
	"flag"
	"time"

	"github.com/ChaimHong/fastbin"
	"github.com/ChaimHong/kfkrpc"
	"github.com/ChaimHong/kfkrpc/example/service1"
	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
)

type config struct {
	RpcServers    map[int]string  `json:"rpc_servers"`
	MiddleServers []string        `json:"middle_servers"`
	Mq            kfkrpc.MqConfig `json:"mq"`
}

func main() {
	var (
		isSrv      = false
		gen        = false
		middle     = false
		configFile = ""
		sid        int
	)

	flag.BoolVar(&isSrv, "is_server", false, "is server")
	flag.IntVar(&sid, "sid", 0, "sid")
	flag.BoolVar(&gen, "gen", false, "gen")
	flag.BoolVar(&middle, "middle", false, "middle")
	flag.StringVar(&configFile, "config", "", "config file")

	flag.Parse()

	if gen {
		server := kfkrpc.NewServer(0)
		server.Register(&service1.ServiceA{})
		fastbin.GenCode()
		return
	}

	cfg := new(config)
	err := util.LoadJsonConfig(cfg, configFile)
	if err != nil {
		panic(err)
	}

	kfkrpc.LoadMqConfig(&cfg.Mq)

	if middle {
		MiddleWare(cfg.MiddleServers[sid])
	} else {
		if isSrv {
			Server(uint16(sid), cfg.RpcServers[sid])
		} else {
			Client(uint16(sid), cfg.MiddleServers)
		}
	}

	for {
		time.Sleep(10 * time.Second)
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

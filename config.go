package kfkrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/funny/debug"
)

var (
	gMqConfig *MqConfig
	Logger    = log.New(os.Stdout, "", log.LstdFlags)
)

type EnDecError struct {
	Err interface{}
}

func (err *EnDecError) Error() string {
	return fmt.Sprintf("%s", err.Err)
}

func LoadMqConfig(mqc *MqConfig) {
	gMqConfig = mqc
}

type MqConfig struct {
	Servers           []string `json:"servers"`
	Ak                string   `json:"username"`
	Password          string   `json:"password"`
	CertFile          string   `json:"cert_file"`
	RequestTopics     []string `json:"requestTopics"`
	RequestConsumerId string   `json:"requestConsumerGroup"`
}

func getConfig() (clusterCfg *cluster.Config) {
	clusterCfg = cluster.NewConfig()

	clusterCfg.Net.SASL.Enable = true
	clusterCfg.Net.SASL.User = gMqConfig.Ak
	clusterCfg.Net.SASL.Password = gMqConfig.Password
	clusterCfg.Net.SASL.Handshake = true

	certBytes, err := ioutil.ReadFile(gMqConfig.CertFile)
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("kafka consumer failed to parse root certificate")
	}

	clusterCfg.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	clusterCfg.Net.TLS.Enable = true
	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	clusterCfg.Group.Return.Notifications = true

	clusterCfg.Version = sarama.V0_10_0_0
	if err = clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		Logger.Println(msg)
		panic(msg)
	}

	return
}

func getRequestConsumer() (consumer *cluster.Consumer) {
	clusterCfg := getConfig()
	var err error
	consumer, err = cluster.NewConsumer(gMqConfig.Servers, gMqConfig.RequestConsumerId, gMqConfig.RequestTopics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		fmt.Println(msg)
		panic(msg)
	}

	return
}

func getProducer() (producer sarama.SyncProducer) {
	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = true
	mqConfig.Net.SASL.User = gMqConfig.Ak
	mqConfig.Net.SASL.Password = gMqConfig.Password
	mqConfig.Net.SASL.Handshake = true
	mqConfig.ClientID = "zdgj_s0"

	Logger.Println(mqConfig.Net.SASL)

	certBytes, err := ioutil.ReadFile(gMqConfig.CertFile)
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("kafka producer failed to parse root certificate")
	}

	mqConfig.Net.TLS.Config = &tls.Config{
		//Certificates:       []tls.Certificate{},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	mqConfig.Net.TLS.Enable = true
	mqConfig.Producer.Return.Successes = true

	if err = mqConfig.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka producer config invalidate. config: %v. err: %v", *gMqConfig, err)
		fmt.Println(msg)
		panic(msg)
	}

	producer, err = sarama.NewSyncProducer(gMqConfig.Servers, mqConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)

		fmt.Printf("%#v\n", debug.StackTrace(0).Frames())
		fmt.Println(msg)
		panic(msg)
	}

	return
}

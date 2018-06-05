package kfkrpc

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/funny/debug"
)

var (
	configPath string
	configName = "kafka.json"
	gMqConfig  = new(MqConfig)
)

func init() {
	var err error

	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configPath = filepath.Join(workPath, "conf")
	LoadJsonConfig(gMqConfig, configName)
	log.Printf("init %#v ", gMqConfig)
}

type MqConfig struct {
	Servers           []string `json:"servers"`
	Ak                string   `json:"username"`
	Password          string   `json:"password"`
	CertFile          string   `json:"cert_file"`
	RequestTopics     []string `json:"requestTopics"`
	RequestConsumerId string   `json:"requestConsumerGroup"`
	ResponeTopics     []string `json:"responeTopics"`
	ResponeConsumerId string   `json:"responeConsumerGroup"`
}

func LoadJsonConfig(config interface{}, filename string) {
	var err error
	var decoder *json.Decoder

	file := OpenFile(filename)
	defer file.Close()

	decoder = json.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		msg := fmt.Sprintf("Decode json fail for config file at %s. Error: %v", filename, err)
		panic(msg)
	}

	json.Marshal(config)
}

func GetFullPath(filename string) string {
	return filepath.Join(configPath, filename)
}

func OpenFile(filename string) *os.File {
	fullPath := filepath.Join(configPath, filename)

	var file *os.File
	var err error

	if file, err = os.Open(fullPath); err != nil {
		msg := fmt.Sprintf("Can not load config at %s. Error: %v", fullPath, err)
		panic(msg)
	}

	return file
}

func getConfig() (clusterCfg *cluster.Config) {
	clusterCfg = cluster.NewConfig()

	log.Printf("confing %#v", gMqConfig)

	clusterCfg.Net.SASL.Enable = true
	clusterCfg.Net.SASL.User = gMqConfig.Ak
	clusterCfg.Net.SASL.Password = gMqConfig.Password
	clusterCfg.Net.SASL.Handshake = true

	certBytes, err := ioutil.ReadFile(GetFullPath(gMqConfig.CertFile))
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
	clusterCfg.Group.Mode = ConsumerModePartitions

	clusterCfg.Version = sarama.V0_10_0_0
	if err = clusterCfg.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka consumer config invalidate. config: %v. err: %v", *clusterCfg, err)
		fmt.Println(msg)
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

func getResponeConsumer() (consumer *cluster.Consumer) {
	clusterCfg := getConfig()
	var err error
	consumer, err = cluster.NewConsumer(gMqConfig.Servers, gMqConfig.ResponeConsumerId, gMqConfig.ResponeTopics, clusterCfg)
	if err != nil {
		msg := fmt.Sprintf("Create kafka consumer error: %v. config: %v", err, clusterCfg)
		fmt.Println(msg)
		panic(msg)
	}

	return
}

func getProducer() (producer sarama.SyncProducer) {
	// sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = true
	mqConfig.Net.SASL.User = gMqConfig.Ak
	mqConfig.Net.SASL.Password = gMqConfig.Password
	mqConfig.Net.SASL.Handshake = true
	mqConfig.ClientID = "zdgj_s0"

	log.Println(mqConfig.Net.SASL)

	certBytes, err := ioutil.ReadFile(GetFullPath(gMqConfig.CertFile))
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

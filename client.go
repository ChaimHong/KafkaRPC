package kfkrpc

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type Client struct {
	producer  sarama.SyncProducer
	consumer  *cluster.Consumer
	sid       uint16
	calls     sync.Map // map[string]*pendingCall
	done      chan bool
	decBuffer *bytes.Buffer
	reading   sync.Mutex
}

type pendingCall struct {
	done chan bool
	data []byte
}

func NewClient(sclient sarama.Client, sid uint16) *Client {
	c := new(Client)
	c.producer = getProducer()
	c.consumer = getResponeConsumer()

	c.decBuffer = new(bytes.Buffer)

	c.sid = sid

	go c.loopMsg()

	return c
}

type Request struct {
	Args  interface{}
	Reply interface{}
}

func (c *Client) Call(sid uint16, serviceMethod string, r *Request, cb func(err error)) {
	correlationId := fmt.Sprintf("%d-%d-%d", c.sid, sid, time.Now().UnixNano())
	pending := &pendingCall{}
	pending.done = make(chan bool, 1)
	c.calls.Store(correlationId, pending)

	encBuffer := new(bytes.Buffer)
	encBuffer.Reset()

	kfkMsg := new(KFKMessage)
	kfkMsg.ServiceMethod = serviceMethod
	kfkMsg.CorrelationId = correlationId
	kfkMsg.ServerId = sid

	argMsg := getIMessage(r.Args)
	argMsgBytes := make([]byte, argMsg.Size())
	argMsg.Marshal(argMsgBytes)

	kfkMsg.Body = argMsgBytes

	iKfkMsg := getIMessage(kfkMsg)
	kfkMsgBytes := make([]byte, iKfkMsg.Size())

	iKfkMsg.Marshal(kfkMsgBytes)

	partition, offset, err0 := c.producer.SendMessage(&sarama.ProducerMessage{
		Topic: gMqConfig.RequestTopics[0],
		Key:   sarama.StringEncoder(correlationId),
		Value: sarama.ByteEncoder(kfkMsgBytes),
	})

	if err0 != nil {
		cb(err0)
		log.Printf("send err %v", err0)
		return
	}
	log.Println(partition, offset)
	log.Printf("input %#v %#v", correlationId, kfkMsgBytes)

	<-pending.done

	log.Printf("done %#v %#v", correlationId, pending.data)

	c.decRespone(pending.data, r.Reply)

	cb(nil)

	c.calls.Delete(correlationId)
	return
}

func (c *Client) decRespone(data []byte, reply interface{}) {
	c.reading.Lock()
	defer c.reading.Unlock()
	c.decBuffer.Reset()

	_, err := c.decBuffer.Write(data)
	util.CheckPanic(err)

	out := new(ResponeMsg)
	iOut := getIMessage(out)
	iOut.Unmarshal(c.decBuffer.Bytes())

	reply.(IMessage).Unmarshal(out.Body)

	fmt.Printf("%#v\n", reply)
}

func (c *Client) loopMsg() {
ConsumerLoop:
	for {
		log.Println("wait message...")
		select {
		case msg, more := <-c.consumer.Messages():
			if !more {
				log.Println("message close")
				break ConsumerLoop
			}
			log.Printf("Consumed message offset %d %s \n", msg.Offset, msg.Key)
			key := string(msg.Key)
			if done, ok := c.calls.Load(key); ok {
				done.(*pendingCall).data = msg.Value
				done.(*pendingCall).done <- true

				c.consumer.MarkOffset(msg, "")
			}
			log.Println("consumer", msg.Offset)
		case err, more := <-c.consumer.Errors():
			if more {
				log.Printf("Kafka consumer error: %v", err.Error())
			}
		case ntf, more := <-c.consumer.Notifications():
			if more {
				log.Printf("Kafka consumer rebalance: %v", ntf)
			}
		case <-c.done:
			break ConsumerLoop
		}
	}

	log.Println("loop message end ")
}

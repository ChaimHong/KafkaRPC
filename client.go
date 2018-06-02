package kfkrpc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
)

type Client struct {
	producer sarama.AsyncProducer
	consumer sarama.PartitionConsumer
	sid      uint16

	calls sync.Map // map[string]*pendingCall
	done  chan bool

	decBuffer *bytes.Buffer

	reading sync.Mutex
}

type pendingCall struct {
	done chan bool
	data []byte
}

func NewClient(sclient sarama.Client, sid uint16) *Client {
	c := new(Client)
	var err error
	c.producer, err = sarama.NewAsyncProducerFromClient(sclient)
	util.CheckPanic(err)

	consumer, e1 := sarama.NewConsumerFromClient(sclient)
	util.CheckPanic(e1)
	c.consumer, err = consumer.ConsumePartition("Reply", 0, 0)
	util.CheckPanic(err)

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
	go func() {
		correlationId := fmt.Sprintf("%d-%d-%d", c.sid, sid, time.Now().UnixNano())
		pending := &pendingCall{}
		pending.done = make(chan bool, 1)
		c.calls.Store(correlationId, pending)

		encBuffer := new(bytes.Buffer)
		encBuffer.Reset()
		enc := gob.NewEncoder(encBuffer)

		kfkMsg := new(KFKMessage)
		kfkMsg.ServiceMethod = serviceMethod
		kfkMsg.CorrelationId = correlationId
		kfkMsg.To = sid
		kfkMsg.From = c.sid

		enc.Encode(r.Args)
		kfkMsg.Body = encBuffer.Bytes()

		encBuffer.Reset()
		enc.Encode(kfkMsg)
		log.Println(kfkMsg, encBuffer.Bytes())

		select {
		case c.producer.Input() <- &sarama.ProducerMessage{
			Topic: "Request",
			Key:   sarama.StringEncoder(correlationId),
			Value: sarama.ByteEncoder(encBuffer.Bytes()),
		}:
		case err := <-c.producer.Errors():
			log.Println("Failed to produce message", err)
		}

		log.Println("args send ok")
		encBuffer.Reset()
		<-pending.done

		c.decRespone(pending.data, r.Reply)

		cb(nil)

		c.calls.Delete(correlationId)
		return
	}()
}

func (c *Client) decRespone(data []byte, reply interface{}) {
	c.reading.Lock()
	defer c.reading.Unlock()
	c.decBuffer.Reset()
	dec := gob.NewDecoder(c.decBuffer)

	_, err := c.decBuffer.Write(data)
	util.CheckPanic(err)

	out := new(KFKMessage)
	e0 := dec.Decode(out)
	util.CheckPanic(e0)

	c.decBuffer.Reset()
	c.decBuffer.Write(out.Body)

	dec.Decode(reply)
	c.decBuffer.Reset()

	fmt.Printf("%#v\n", reply)
}

func (c *Client) loopMsg() {
ConsumerLoop:
	for {
		select {
		case msg := <-c.consumer.Messages():
			log.Printf("Consumed message offset %d %s \n", msg.Offset, msg.Key)
			key := string(msg.Key)
			if done, ok := c.calls.Load(key); ok {
				fmt.Println("consumer", msg.Value)
				done.(*pendingCall).data = msg.Value
				done.(*pendingCall).done <- true
			}
			log.Println("consumer", msg.Offset)
		case <-c.done:
			break ConsumerLoop
		}
	}

	log.Println("loop message end ")
}

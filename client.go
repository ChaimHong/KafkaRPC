package kfkrpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
)

type Client struct {
	producer sarama.SyncProducer

	sid       uint16
	calls     sync.Map // map[string]*pendingCall
	done      chan bool
	decBuffer *bytes.Buffer
	reading   sync.Mutex

	// middleware
	middlesId uint32
	middles   sync.Map // map [uint32]*MiddleClient
}

type pendingCall struct {
	done chan bool
	err  string
	data []byte
}

func NewClient(sid uint16, middles []string) *Client {
	c := new(Client)
	c.producer = getProducer()

	c.decBuffer = new(bytes.Buffer)

	c.sid = sid

	for _, addr := range middles {
		addr = strings.Replace(addr, "0.0.0.0", "127.0.0.1", 1)
		middle, err := newMiddleClient(addr, sid, c)
		if err == nil {
			id := atomic.AddUint32(&c.middlesId, 1)
			c.middles.Store(id, middle)
		} else {
			Logger.Printf("middle %s could not connect:%s", addr, err)
		}
	}

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
	kfkMsg.ClientId = c.sid

	kfkMsg.Body = getIMessageBytes(r.Args)
	kfkMsgBytes := getIMessageBytes(kfkMsg)

	_, _, err0 := c.producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: gMqConfig.RequestTopics[0],
			Key:   sarama.StringEncoder(correlationId),
			Value: sarama.ByteEncoder(kfkMsgBytes),
		})

	if err0 != nil {
		cb(err0)
		Logger.Printf("send err %v", err0)
		return
	}

	<-pending.done

	if pending.err == "" {
		getIMessage(r.Reply).Unmarshal(pending.data)
		cb(nil)
	} else {
		cb(errors.New(pending.err))
	}

	c.calls.Delete(correlationId)
	return
}

type MiddleClient struct {
	id   uint
	conn net.Conn
}

func newMiddleClient(addr string, sid uint16, client *Client) (*MiddleClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	buf := [2]byte{}
	binary.LittleEndian.PutUint16(buf[:], sid)
	_, err = conn.Write(buf[:])
	if err != nil {
		return nil, err
	}

	mc := &MiddleClient{
		conn: conn,
	}
	go mc.loop(client)

	return mc, nil
}

type MiddleReply struct {
	CorrelationId string
	Error         string
	Body          []byte
}

func (mc *MiddleClient) loop(client *Client) (err error) {
	for {
		var msg []byte
		msg, err = Receive(mc.conn)
		if err != nil {
			break
		}

		reply := new(MiddleReply)
		getIMessage(reply).Unmarshal(msg)

		if done, ok := client.calls.Load(reply.CorrelationId); ok {
			call := done.(*pendingCall)
			call.data = reply.Body
			call.err = reply.Error
			call.done <- true
		}

		select {
		case <-client.done:
			break
		default:
		}
	}

	client.middles.Delete(mc.id)
	Logger.Println("MiddleClient loop end")

	return
}

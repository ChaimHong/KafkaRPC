package kfkrpc

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"sync"

	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
)

type Server struct {
	producer sarama.AsyncProducer
	consumer sarama.PartitionConsumer
	sid      uint16

	decBuffer *bytes.Buffer
	dec       *gob.Decoder
	encBuffer *bytes.Buffer
	enc       *gob.Encoder

	server  *rpc.Server
	cSeqMap sync.Map // map[uint64]*KFKMessage
	cSeq    uint64   //
}

func NewServer(sid uint16) *Server {
	srv := &Server{
		server:    rpc.NewServer(),
		sid:       sid,
		decBuffer: new(bytes.Buffer),
		encBuffer: new(bytes.Buffer),
	}

	return srv
}

func (srv *Server) Serve(client sarama.Client) {
	var err error
	srv.producer, err = sarama.NewAsyncProducerFromClient(client)
	util.CheckPanic(err)

	consumer, e1 := sarama.NewConsumerFromClient(client)
	util.CheckPanic(e1)
	srv.consumer, err = consumer.ConsumePartition("Request", 0, 0)
	util.CheckPanic(err)

	srv.server.ServeCodec(srv)
}

func (srv *Server) Register(rcvr interface{}) {
	err := srv.server.Register(rcvr)
	if err != nil {
		panic(err)
	}
}

var _ rpc.ServerCodec = (*Server)(nil)

var errmsg = "sid is not match"
var errorSidNotMatch = errors.New(errmsg)

func (srv *Server) ReadRequestHeader(r *rpc.Request) (e error) {
	select {
	case msg := <-srv.consumer.Messages():
		log.Printf("Consumed message offset %d %s \n", msg.Offset, msg.Key)

		srv.decBuffer.Reset()
		// TODO optimize
		srv.dec = gob.NewDecoder(srv.decBuffer)

		kfkMessage := new(KFKMessage)
		srv.decBuffer.Write(msg.Value)
		log.Println(srv.decBuffer.Bytes())

		e = srv.dec.Decode(kfkMessage)
		util.CheckPanic(e)

		fmt.Println(kfkMessage)

		// 不匹配抛出异常
		if kfkMessage.To != srv.sid {
			srv.decBuffer.Reset()
			return errorSidNotMatch
		}

		srv.cSeq++
		srv.cSeqMap.Store(srv.cSeq, kfkMessage)
		r.Seq = srv.cSeq

		r.ServiceMethod = kfkMessage.ServiceMethod

		srv.decBuffer.Reset()

		_, e = srv.decBuffer.Write(kfkMessage.Body)
		util.CheckPanic(e)

		log.Println("read header end")
		return nil
	}
}

func (srv *Server) ReadRequestBody(body interface{}) (err error) {
	if body == nil {
		srv.decBuffer.Reset()
		return nil
	}

	err = srv.dec.Decode(body)
	srv.decBuffer.Reset()
	return nil
}

func (srv *Server) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	defer srv.cSeqMap.Delete(r.Seq)
	if r.Error == errmsg {
		return nil
	}

	rawMessage, ok := srv.cSeqMap.Load(r.Seq)
	if !ok {
		return errors.New("seq do not map")
	}

	fromMsg := rawMessage.(*KFKMessage)
	fromMsg.To, fromMsg.From = fromMsg.From, fromMsg.To

	srv.encBuffer.Reset()
	srv.enc = gob.NewEncoder(srv.encBuffer)

	srv.enc.Encode(body)
	fromMsg.Body = srv.encBuffer.Bytes()
	srv.encBuffer.Reset()

	srv.enc.Encode(fromMsg)
	msgBytes := srv.encBuffer.Bytes()

	fmt.Println("reply", body)

	select {
	case srv.producer.Input() <- &sarama.ProducerMessage{
		Topic: "Reply",
		Key:   sarama.StringEncoder(fromMsg.CorrelationId),
		Value: sarama.ByteEncoder(msgBytes),
	}:
	case err := <-srv.producer.Errors():
		log.Println("Failed to produce message", err)
	}

	fmt.Println("reply ok ")

	return nil
}

func (srv *Server) Close() error {
	err0 := srv.producer.Close()
	err1 := srv.consumer.Close()
	err := ""
	if err0 != nil {
		err += err0.Error()
	}
	if err1 != nil {
		err += err1.Error()
	}

	if err == "" {
		return nil
	}

	return errors.New(err)
}

package kfkrpc

import (
	"bytes"
	"errors"
	"log"
	"net/rpc"
	"reflect"
	"sync"

	"github.com/ChaimHong/util"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/funny/fastbin"
)

type Server struct {
	producer sarama.SyncProducer
	consumer *cluster.Consumer

	sid uint16

	decBuffer *bytes.Buffer

	server  *rpc.Server
	cSeqMap sync.Map // map[uint64]*KFKMessage
	cSeq    uint64   //
}

func NewServer(sid uint16) *Server {
	srv := &Server{
		server:    rpc.NewServer(),
		sid:       sid,
		decBuffer: new(bytes.Buffer),
	}

	return srv
}

func (srv *Server) Serve(client sarama.Client) {
	srv.producer = getProducer()
	srv.consumer = getRequestConsumer()

	srv.server.ServeCodec(srv)
}

func (srv *Server) Register(rcvr interface{}) {
	err := srv.server.Register(rcvr)
	if err != nil {
		panic(err)
	}

	typ := reflect.TypeOf(rcvr)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		mtype := method.Type

		argType := mtype.In(1)

		replyType := mtype.In(2)

		fastbin.RegisterType(argType)
		fastbin.RegisterType(replyType)
	}
}

var _ rpc.ServerCodec = (*Server)(nil)

var errmsg = "sid is not match"
var errorSidNotMatch = errors.New(errmsg)

func (srv *Server) ReadRequestHeader(r *rpc.Request) (e error) {
	select {
	case msg, more := <-srv.consumer.Messages():
		if !more {
			log.Println("consumer is close")
			return
		}
		return srv.decHeader(r, msg)
	case err, more := <-srv.consumer.Errors():
		if more {
			log.Printf("Kafka consumer error: %v", err.Error())
		}
	case ntf, more := <-srv.consumer.Notifications():
		if more {
			log.Printf("Kafka consumer rebalance: %v", ntf)
		}
	}

	return nil
}

func (srv *Server) decHeader(r *rpc.Request, msg *sarama.ConsumerMessage) (err error) {
	defer func() {
		if e0 := recover(); e0 != nil {
			log.Printf("dec header err %s", e0)
			err = errorSidNotMatch
		}
	}()

	var e error
	kfkMessage := new(KFKMessage)
	srv.decBuffer.Reset()
	_, e = srv.decBuffer.Write(msg.Value)
	util.CheckPanic(e)

	getIMessage(kfkMessage).Unmarshal(srv.decBuffer.Bytes())

	log.Printf("header %#v", kfkMessage)

	// 不匹配抛出异常
	if kfkMessage.ServerId != srv.sid {
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

	srv.consumer.MarkOffset(msg, "")
	return
}

func (srv *Server) ReadRequestBody(body interface{}) (err error) {
	if body == nil {
		srv.decBuffer.Reset()
		return nil
	}

	body.(IMessage).Unmarshal(srv.decBuffer.Bytes())

	log.Printf("header %#v", body)
	return nil
}

func (srv *Server) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if r.Error == errmsg {
		return nil
	}

	rawMessage, ok := srv.cSeqMap.Load(r.Seq)
	srv.cSeqMap.Delete(r.Seq)

	if !ok {
		return errors.New("seq do not map")
	}

	reply := body.(IMessage)
	bytes := make([]byte, reply.Size())
	reply.Marshal(bytes)

	rsp := new(ResponeMsg)
	rsp.Error = r.Error
	rsp.Body = bytes

	iRsp := getIMessage(rsp)
	msgBytes := make([]byte, iRsp.Size())
	iRsp.Marshal(msgBytes)

	partition, offset, err0 := srv.producer.SendMessage(&sarama.ProducerMessage{
		Topic: gMqConfig.ResponeTopics[0],
		Key:   sarama.StringEncoder(rawMessage.(*KFKMessage).CorrelationId),
		Value: sarama.ByteEncoder(msgBytes),
	})

	if err0 != nil {
		log.Printf("send err %v", err0)
		return err0
	}

	log.Println("send ok", partition, offset, err0)

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

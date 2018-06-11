package kfkrpc

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/funny/cmd"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type MiddleWare struct {
	consumer *cluster.Consumer
	srvNodes sync.Map // map[uint16]*rpc.Client to rpc.Server

	lsn     net.Listener
	clients map[uint16]net.Conn //
	done    chan bool
}

type Node struct {
	Id   uint16
	Addr string
}

func NewMiddleWare(addr string) (mw *MiddleWare, err error) {
	mw = &MiddleWare{
		consumer: getRequestConsumer(),
		done:     make(chan bool, 1),
		clients:  make(map[uint16]net.Conn, 10),
	}

	mw.lsn, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return mw, nil
}

func (mw *MiddleWare) Server() {
	// server for client
	go mw.server()

	// loop kfk message
	go mw.loop()

	mw.registerCommands()

	go func() {
		signalChan1 := make(chan os.Signal, 10)
		signalChan2 := make(chan os.Signal, 10)

		signal.Notify(signalChan1, syscall.SIGTERM)
		signal.Notify(signalChan2, syscall.SIGUSR1)

		for {
			select {
			case <-signalChan1:
				log.Println("chan 1")
			case <-signalChan2:
				log.Println("chan 2")
			}
		}
	}()

	go shellLoop()
}

func shellLoop() {
	input := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("shellLoop> ")

		command, _, err := input.ReadLine()
		log.Println(command, err)
		if err != nil {
			break
		}

		log.Println(command, err)

		if _, ok := cmd.Process(string(command)); !ok {
			Logger.Printf("unknow command: %s", string(command))
			cmd.Help(os.Stderr)
		}
	}
}

func (mw *MiddleWare) registerCommands() {
	cmd.Register(
		"help",
		"Print this screen",
		func() {
			cmd.Help(os.Stderr)
		},
	)
	cmd.Register(
		"add node ([0-9]{0,4}) (.*)",
		"add node 1 127.0.0.1:3000",
		func(args []string) {
			sid, err := strconv.Atoi(args[1])
			if err != nil {
				return
			}
			mw.AddNode(Node{Id: uint16(sid), Addr: args[2]})
		},
	)
}

func (mw *MiddleWare) server() error {
	for {
		conn, err := mw.lsn.Accept()
		if err != nil {
			return err
		}
		go func() {
			sidBytes := [2]byte{}
			_, err := conn.Read(sidBytes[:])
			sid := binary.LittleEndian.Uint16(sidBytes[:])
			mw.clients[sid] = conn

			if err != nil {
				conn.Close()
				return
			}
		}()
	}
}

func (mw *MiddleWare) loop() {
ConsumerLoop:
	for {
		select {
		case msg, ok := <-mw.consumer.Messages():
			if !ok {
				Logger.Println("message close")
				break ConsumerLoop
			}
			Logger.Printf("Consumed message offset %d %s \n", msg.Offset, msg.Key)
			if err, request := newKFKMessage(msg.Value); err == nil {
				mw.doConsumer(msg, request)
			} else {
				Logger.Println(err)
			}
		case err, ok := <-mw.consumer.Errors():
			if ok {
				Logger.Printf("Kafka consumer error: %v", err.Error())
			}
		case ntf, ok := <-mw.consumer.Notifications():
			if ok {
				Logger.Printf("Kafka consumer rebalance: %v", ntf)
			}
		case <-mw.done:
			break ConsumerLoop
		}
	}

	Logger.Println("loop message end ")
}

func (mw *MiddleWare) doConsumer(msg *sarama.ConsumerMessage, request *KFKMessage) {
	v, ok := mw.srvNodes.Load(request.ServerId)
	if !ok {
		Logger.Printf("server id(%d) is not exist", request.ServerId)
		return
	}

	client := v.(*rpc.Client)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				Logger.Printf("doconsumer err:%s", err)
			}
		}()

		var reply = &RpcResponse{}
		err := client.Call(request.ServiceMethod, request.Body, reply)
		if err == io.ErrUnexpectedEOF || err == rpc.ErrShutdown {
			client.Close()
			mw.srvNodes.Delete(client)
		} else {
			// 返回给请求的服务器
			conn, ok := mw.clients[request.ClientId]
			if !ok {
				panic(fmt.Sprintf("client(%d)'s conn is not exist", request.ClientId))
			}

			rsp := &MiddleReply{CorrelationId: string(msg.Key)}
			if err == nil {
				rsp.Body = *reply.Result
			} else {
				rsp.Error = err.Error()
			}

			err := Send(conn, getIMessageBytes(rsp))
			if err == nil {
				mw.consumer.MarkOffset(msg, "")
			} else {
				Logger.Printf("send err %s", err)
			}
		}
	}()
}

func (mw *MiddleWare) Close() {
	close(mw.done)
}

type clientCodec struct {
	rw io.ReadWriteCloser //

	req *RpcRequest
	rsp *RpcResponse
}

func newClientCodec(rw io.ReadWriteCloser) *clientCodec {
	return &clientCodec{
		rw:  rw,
		req: new(RpcRequest),
		rsp: new(RpcResponse),
	}
}

var _ rpc.ClientCodec

func (cc *clientCodec) WriteRequest(req *rpc.Request, args interface{}) error {
	cc.req.reset()

	cc.req.Method = req.ServiceMethod
	cc.req.Id = req.Seq
	cc.req.Args = args.([]byte)

	return Send(cc.rw, getIMessageBytes(cc.req))
}

func (cc *clientCodec) ReadResponseHeader(rsp *rpc.Response) (err error) {
	defer func() {
		if e0 := recover(); e0 != nil {
			err = &EnDecError{e0}
		}
	}()

	var msg []byte
	msg, err = Receive(cc.rw)
	if err != nil {
		return
	}

	cc.rsp.reset()
	getIMessage(cc.rsp).Unmarshal(msg)
	rsp.Seq = cc.rsp.Id
	rsp.Error = cc.rsp.Error

	return nil
}

func (cc *clientCodec) ReadResponseBody(x interface{}) error {
	if x == nil {
		return nil
	}

	if reply, ok := x.(*RpcResponse); ok {
		*reply = *cc.rsp
	} else {
		Logger.Println("middleware read body err")
	}

	return nil
}

func (cc *clientCodec) Close() error {
	return cc.rw.Close()
}

func (mw *MiddleWare) AddNode(node Node) (err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", node.Addr)
	if err != nil {
		return
	}

	client := rpc.NewClientWithCodec(newClientCodec(conn))

	mw.srvNodes.Store(node.Id, client)

	Logger.Println("add node success", node)
	return
}

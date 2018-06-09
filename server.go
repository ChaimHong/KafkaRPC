package kfkrpc

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"reflect"
	"sync"

	"github.com/funny/fastbin"
)

type Server struct {
	sid     uint16
	server  *rpc.Server
	cSeqMap sync.Map // map[uint64]*KFKMessage
	cSeq    uint64   //
}

func NewServer(sid uint16) *Server {
	srv := &Server{
		server: rpc.NewServer(),
		sid:    sid,
	}

	return srv
}

func (srv *Server) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go srv.server.ServeCodec(newServerCodec(conn))
	}
}

func (srv *Server) Register(rcvr interface{}) error {
	err := srv.server.Register(rcvr)
	if err != nil {
		return err
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

	return nil
}

type serverCodec struct {
	rw  io.ReadWriteCloser
	req *RpcRequest

	mutex   sync.Mutex // protects seq, pending
	seq     uint64
	pending map[uint64]uint64
}

func newServerCodec(rwc io.ReadWriteCloser) *serverCodec {
	return &serverCodec{
		rw:      rwc,
		req:     new(RpcRequest),
		pending: make(map[uint64]uint64, 10),
	}
}

func (sc *serverCodec) ReadRequestHeader(r *rpc.Request) (e error) {
	var msg []byte
	msg, e = Receive(sc.rw)
	if e != nil {
		return
	}

	sc.req.reset()
	getIMessage(sc.req).Unmarshal(msg)

	r.ServiceMethod = sc.req.Method

	sc.mutex.Lock()
	sc.seq++
	sc.pending[sc.seq] = sc.req.Id
	r.Seq = sc.seq
	sc.mutex.Unlock()

	return nil
}

func (sc *serverCodec) ReadRequestBody(body interface{}) (err error) {
	defer func() {
		if e0 := recover(); e0 != nil {
			err = errors.New(fmt.Sprintf("%s", e0))
		}
	}()
	bytes := sc.req.Args
	im := getIMessage(body)
	im.Unmarshal(bytes)

	return nil
}

func (sc *serverCodec) WriteResponse(r *rpc.Response, x interface{}) (err error) {
	sc.mutex.Lock()
	b, ok := sc.pending[r.Seq]
	if !ok {
		sc.mutex.Unlock()
		return errors.New("invalid sequence number in response")
	}
	delete(sc.pending, r.Seq)
	sc.mutex.Unlock()

	resp := &RpcResponse{Id: b}
	if r.Error == "" {
		bytes := getIMessageBytes(x)
		resp.Result = new(RawBytes)
		*(resp.Result) = bytes
	} else {
		resp.Error = r.Error
	}

	return Send(sc.rw, getIMessageBytes(resp))
}

func (sc *serverCodec) Close() error {
	return sc.rw.Close()
}

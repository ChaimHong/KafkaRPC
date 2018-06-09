package kfkrpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type KFKMessage struct {
	ServiceMethod string // 服务名称
	CorrelationId string // application use - correlation identifier
	ServerId      uint16 //
	ClientId      uint16
	Body          []byte
}

func getIMessage(msg interface{}) IMessage {
	return msg.(IMessage)
}

func getIMessageBytes(msg interface{}) (bytes []byte) {
	im := getIMessage(msg)
	bytes = make([]byte, im.Size())
	im.Marshal(bytes)
	return
}

type IMessage interface {
	Size() int
	Marshal([]byte) int
	Unmarshal([]byte) int
}

type ResponeMsg struct {
	Error string
	Body  []byte
}

func newKFKMessage(data []byte) (err error, kfkMsg *KFKMessage) {
	defer func() {
		if e0 := recover(); e0 != nil {
			err = errors.New(fmt.Sprintf("%s", e0))
			kfkMsg = nil
		}
	}()

	kfkMsg = new(KFKMessage)
	getIMessage(kfkMsg).Unmarshal(data)

	return
}

func Send(w io.Writer, data []byte) error {
	buff := make([]byte, 4+len(data))
	binary.LittleEndian.PutUint32(buff[:4], uint32(len(data)))
	copy(buff[4:], data)
	_, err := w.Write(buff)
	return err
}

func Receive(r io.Reader) (msg []byte, err error) {
	var head = [4]byte{}
	_, err = r.Read(head[:])
	if err != nil {
		return nil, err
	}

	len := binary.LittleEndian.Uint32(head[:])
	msg = make([]byte, len)
	_, err = r.Read(msg)

	return
}

type RpcRequest struct {
	Method string
	Id     uint64
	Args   []byte
}

func (req *RpcRequest) reset() {
	req.Method = ""
	req.Id = 0
	req.Args = nil
}

type RawBytes []byte

type RpcResponse struct {
	Id     uint64
	Result *RawBytes
	Error  string
}

func (rsp *RpcResponse) reset() {
	rsp.Id = 0
	rsp.Error = ""
	rsp.Result = nil
}

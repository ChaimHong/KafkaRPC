package kfkrpc

import "github.com/funny/fastapi"

type Message struct {
	FuncName  string //
	RPCID     string //
	Producter int16  //
	Consumer  int16  //
	Args      string // 值
	Replys    string //
}

type KFKMessage struct {
	ServiceMethod string
	CorrelationId string // application use - correlation identifier
	To            uint16 // application use - address to reply to (ex: RPC) serverId
	From          uint16 //
	Body          []byte
}

func getFastApiMessage(msg interface{}) fastapi.Message {
	return msg.(fastapi.Message)
}

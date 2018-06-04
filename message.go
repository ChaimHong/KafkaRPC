package kfkrpc

type Message struct {
	FuncName  string //
	RPCID     string //
	Producter int16  //
	Consumer  int16  //
	Args      string // 值
	Replys    string //
}

type KFKMessage struct {
	ServiceMethod string // 服务名称
	CorrelationId string // application use - correlation identifier
	ServerId      uint16 //
	Body          []byte
}

func getIMessage(msg interface{}) IMessage {
	return msg.(IMessage)
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

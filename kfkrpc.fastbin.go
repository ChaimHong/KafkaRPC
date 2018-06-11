package kfkrpc

import "encoding/binary"
import "github.com/ChaimHong/gobuf"

var _ gobuf.Struct = (*ResponeMsg)(nil)

func (s *ResponeMsg) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.Error))) + len(s.Error)
	size += gobuf.UvarintSize(uint64(len(s.Body))) + len(s.Body)
	return size
}

func (s *ResponeMsg) Marshal(b []byte) int {
	var n int
	n += binary.PutUvarint(b[n:], uint64(len(s.Error)))
	copy(b[n:], s.Error)
	n += len(s.Error)
	n += binary.PutUvarint(b[n:], uint64(len(s.Body)))
	copy(b[n:], s.Body)
	n += len(s.Body)
	return n
}

func (s *ResponeMsg) Unmarshal(b []byte) int {
	var n int
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Error = string(b[n : n+int(l)])
		n += int(l)
	}
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Body = make([]byte, l)
		copy(s.Body, b[n:n+int(l)])
		n += int(l)
	}
	return n
}

var _ gobuf.Struct = (*MiddleReply)(nil)

func (s *MiddleReply) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.CorrelationId))) + len(s.CorrelationId)
	size += gobuf.UvarintSize(uint64(len(s.Error))) + len(s.Error)
	size += gobuf.UvarintSize(uint64(len(s.Body))) + len(s.Body)
	return size
}

func (s *MiddleReply) Marshal(b []byte) int {
	var n int
	n += binary.PutUvarint(b[n:], uint64(len(s.CorrelationId)))
	copy(b[n:], s.CorrelationId)
	n += len(s.CorrelationId)
	n += binary.PutUvarint(b[n:], uint64(len(s.Error)))
	copy(b[n:], s.Error)
	n += len(s.Error)
	n += binary.PutUvarint(b[n:], uint64(len(s.Body)))
	copy(b[n:], s.Body)
	n += len(s.Body)
	return n
}

func (s *MiddleReply) Unmarshal(b []byte) int {
	var n int
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.CorrelationId = string(b[n : n+int(l)])
		n += int(l)
	}
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Error = string(b[n : n+int(l)])
		n += int(l)
	}
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Body = make([]byte, l)
		copy(s.Body, b[n:n+int(l)])
		n += int(l)
	}
	return n
}

var _ gobuf.Struct = (*RpcRequest)(nil)

func (s *RpcRequest) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.Method))) + len(s.Method)
	size += 8
	size += gobuf.UvarintSize(uint64(len(s.Args))) + len(s.Args)
	return size
}

func (s *RpcRequest) Marshal(b []byte) int {
	var n int
	n += binary.PutUvarint(b[n:], uint64(len(s.Method)))
	copy(b[n:], s.Method)
	n += len(s.Method)
	binary.LittleEndian.PutUint64(b[n:], uint64(s.Id))
	n += 8
	n += binary.PutUvarint(b[n:], uint64(len(s.Args)))
	copy(b[n:], s.Args)
	n += len(s.Args)
	return n
}

func (s *RpcRequest) Unmarshal(b []byte) int {
	var n int
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Method = string(b[n : n+int(l)])
		n += int(l)
	}
	s.Id = uint64(binary.LittleEndian.Uint64(b[n:]))
	n += 8
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Args = make([]byte, l)
		copy(s.Args, b[n:n+int(l)])
		n += int(l)
	}
	return n
}

var _ gobuf.Struct = (*RpcResponse)(nil)

func (s *RpcResponse) Size() int {
	var size int
	size += 8
	size += 1
	if s.Result != nil {
		size += gobuf.UvarintSize(uint64(len(*s.Result))) + len(*s.Result)
	}
	size += gobuf.UvarintSize(uint64(len(s.Error))) + len(s.Error)
	return size
}

func (s *RpcResponse) Marshal(b []byte) int {
	var n int
	binary.LittleEndian.PutUint64(b[n:], uint64(s.Id))
	n += 8
	if s.Result != nil {
		b[n] = 1
		n++
		n += binary.PutUvarint(b[n:], uint64(len(*s.Result)))
		copy(b[n:], *s.Result)
		n += len(*s.Result)
	} else {
		b[n] = 0
		n++
	}
	n += binary.PutUvarint(b[n:], uint64(len(s.Error)))
	copy(b[n:], s.Error)
	n += len(s.Error)
	return n
}

func (s *RpcResponse) Unmarshal(b []byte) int {
	var n int
	s.Id = uint64(binary.LittleEndian.Uint64(b[n:]))
	n += 8
	if b[n] != 0 {
		n += 1
		val1 := new(RawBytes)
		{
			l, x := binary.Uvarint(b[n:])
			n += x
			*val1 = make([]byte, l)
			copy(*val1, b[n:n+int(l)])
			n += int(l)
		}
		s.Result = val1
	} else {
		n += 1
	}
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Error = string(b[n : n+int(l)])
		n += int(l)
	}
	return n
}

var _ gobuf.Struct = (*KFKMessage)(nil)

func (s *KFKMessage) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.ServiceMethod))) + len(s.ServiceMethod)
	size += gobuf.UvarintSize(uint64(len(s.CorrelationId))) + len(s.CorrelationId)
	size += 2
	size += 2
	size += gobuf.UvarintSize(uint64(len(s.Body))) + len(s.Body)
	return size
}

func (s *KFKMessage) Marshal(b []byte) int {
	var n int
	n += binary.PutUvarint(b[n:], uint64(len(s.ServiceMethod)))
	copy(b[n:], s.ServiceMethod)
	n += len(s.ServiceMethod)
	n += binary.PutUvarint(b[n:], uint64(len(s.CorrelationId)))
	copy(b[n:], s.CorrelationId)
	n += len(s.CorrelationId)
	binary.LittleEndian.PutUint16(b[n:], uint16(s.ServerId))
	n += 2
	binary.LittleEndian.PutUint16(b[n:], uint16(s.ClientId))
	n += 2
	n += binary.PutUvarint(b[n:], uint64(len(s.Body)))
	copy(b[n:], s.Body)
	n += len(s.Body)
	return n
}

func (s *KFKMessage) Unmarshal(b []byte) int {
	var n int
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.ServiceMethod = string(b[n : n+int(l)])
		n += int(l)
	}
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.CorrelationId = string(b[n : n+int(l)])
		n += int(l)
	}
	s.ServerId = uint16(binary.LittleEndian.Uint16(b[n:]))
	n += 2
	s.ClientId = uint16(binary.LittleEndian.Uint16(b[n:]))
	n += 2
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Body = make([]byte, l)
		copy(s.Body, b[n:n+int(l)])
		n += int(l)
	}
	return n
}

package kfkrpc

import "encoding/binary"
import "github.com/ChaimHong/gobuf"

var _ gobuf.Struct = (*KFKMessage)(nil)

func (s *KFKMessage) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.ServiceMethod))) + len(s.ServiceMethod)
	size += gobuf.UvarintSize(uint64(len(s.CorrelationId))) + len(s.CorrelationId)
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
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.Body = make([]byte, l)
		copy(s.Body, b[n:n+int(l)])
		n += int(l)
	}
	return n
}

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

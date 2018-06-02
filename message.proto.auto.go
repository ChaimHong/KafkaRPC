package kfkrpc

import "encoding/binary"
import "github.com/ChaimHong/gobuf"

var _ gobuf.Struct = (*KFKMessage)(nil)

func (s *KFKMessage) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.CorrelationId))) + len(s.CorrelationId)
	size += 2
	size += 2
	size += gobuf.UvarintSize(uint64(len(s.Body))) + len(s.Body)
	return size
}

func (s *KFKMessage) Marshal(b []byte) int {
	var n int
	n += binary.PutUvarint(b[n:], uint64(len(s.CorrelationId)))
	copy(b[n:], s.CorrelationId)
	n += len(s.CorrelationId)
	binary.LittleEndian.PutUint16(b[n:], uint16(s.To))
	n += 2
	binary.LittleEndian.PutUint16(b[n:], uint16(s.From))
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
		s.CorrelationId = string(b[n : n+int(l)])
		n += int(l)
	}
	s.To = uint16(binary.LittleEndian.Uint16(b[n:]))
	n += 2
	s.From = uint16(binary.LittleEndian.Uint16(b[n:]))
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

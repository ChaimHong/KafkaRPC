package service1

import "encoding/binary"
import "github.com/ChaimHong/gobuf"

var _ gobuf.Struct = (*AIn)(nil)

func (s *AIn) Size() int {
	var size int
	size += gobuf.VarintSize(int64(s.V))
	return size
}

func (s *AIn) Marshal(b []byte) int {
	var n int
	n += binary.PutVarint(b[n:], int64(s.V))
	return n
}

func (s *AIn) Unmarshal(b []byte) int {
	var n int
	{
		v, x := binary.Varint(b[n:])
		s.V = int(v)
		n += x
	}
	return n
}

var _ gobuf.Struct = (*AOut)(nil)

func (s *AOut) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.V))) + len(s.V)
	return size
}

func (s *AOut) Marshal(b []byte) int {
	var n int
	n += binary.PutUvarint(b[n:], uint64(len(s.V)))
	copy(b[n:], s.V)
	n += len(s.V)
	return n
}

func (s *AOut) Unmarshal(b []byte) int {
	var n int
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.V = string(b[n : n+int(l)])
		n += int(l)
	}
	return n
}

var _ gobuf.Struct = (*BIn)(nil)

func (s *BIn) Size() int {
	var size int
	size += gobuf.VarintSize(int64(s.B))
	return size
}

func (s *BIn) Marshal(b []byte) int {
	var n int
	n += binary.PutVarint(b[n:], int64(s.B))
	return n
}

func (s *BIn) Unmarshal(b []byte) int {
	var n int
	{
		v, x := binary.Varint(b[n:])
		s.B = int(v)
		n += x
	}
	return n
}

var _ gobuf.Struct = (*BOut)(nil)

func (s *BOut) Size() int {
	var size int
	size += gobuf.UvarintSize(uint64(len(s.B))) + len(s.B)
	return size
}

func (s *BOut) Marshal(b []byte) int {
	var n int
	n += binary.PutUvarint(b[n:], uint64(len(s.B)))
	copy(b[n:], s.B)
	n += len(s.B)
	return n
}

func (s *BOut) Unmarshal(b []byte) int {
	var n int
	{
		l, x := binary.Uvarint(b[n:])
		n += x
		s.B = string(b[n : n+int(l)])
		n += int(l)
	}
	return n
}

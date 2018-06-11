package service2

import "encoding/binary"
import "github.com/ChaimHong/gobuf"

var _ gobuf.Struct = (*A)(nil)

func (s *A) Size() int {
	var size int
	size += 1
	if s.F1 != nil {
		size += gobuf.UvarintSize(uint64(len(*s.F1))) + len(*s.F1)
	}
	return size
}

func (s *A) Marshal(b []byte) int {
	var n int
	if s.F1 != nil {
		b[n] = 1
		n++
		n += binary.PutUvarint(b[n:], uint64(len(*s.F1)))
		copy(b[n:], *s.F1)
		n += len(*s.F1)
	} else {
		b[n] = 0
		n++
	}
	return n
}

func (s *A) Unmarshal(b []byte) int {
	var n int
	if b[n] != 0 {
		n += 1
		val1 := new(Raw)
		{
			l, x := binary.Uvarint(b[n:])
			n += x
			*val1 = make([]byte, l)
			copy(*val1, b[n:n+int(l)])
			n += int(l)
		}
		s.F1 = val1
	} else {
		n += 1
	}
	return n
}

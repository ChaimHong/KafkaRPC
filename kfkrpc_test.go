package kfkrpc_test

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	b := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[:], 6)
	fmt.Println(b)
	c := []byte{4}
	copy(b[1:], c)
	fmt.Println(b, c)

}

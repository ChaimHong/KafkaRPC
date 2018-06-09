package kfkrpc_test

import (
	"encoding/binary"
	"fmt"
	"testing"
)

// func TestRV(t *testing.T) {
// 	v := new(service1.AIn)
// 	i := reflect.New(reflect.TypeOf(v).Elem())
// 	log.Printf("%#v\n", i)

// 	_ = i.(kfkrpc.IMessage)
// }

func TestBytes(t *testing.T) {
	b := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint16(b[:], 6)
	fmt.Println(b)
	c := []byte{4}
	copy(b[1:], c)
	fmt.Println(b, c)

}

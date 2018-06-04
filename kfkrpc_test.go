package kfkrpc_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/ChaimHong/kfkrpc"
	"github.com/ChaimHong/kfkrpc/example/service1"
)

func TestRV(t *testing.T) {
	v := new(service1.AIn)
	i := reflect.New(reflect.TypeOf(v).Elem())
	log.Printf("%#v\n", i)

	_ = i.(kfkrpc.IMessage)
}

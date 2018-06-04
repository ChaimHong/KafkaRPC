package service1_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/ChaimHong/kfkrpc/example/service1"
)

func TestRV(t *testing.T) {
	v := new(service1.AIn)
	i := reflect.New(reflect.TypeOf(v).Elem())
	log.Printf("%#v\n", i)
}

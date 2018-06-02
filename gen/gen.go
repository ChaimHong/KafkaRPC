package main

import (
	"encoding/json"
	"flag"
	"go/types"
	"os"
	"reflect"

	"github.com/ChaimHong/ReflectType2GoType"
	"github.com/ChaimHong/gobuf/parser"
	"github.com/ChaimHong/gobuf/plugins/go"
	"github.com/ChaimHong/kfkrpc"
	"github.com/ChaimHong/util"
)

func main() {
	var filename string
	flag.StringVar(&filename, "filename", "message.proto.auto.go", "kfkrpcMessage.proto")

	flag.Parse()

	newConver := rtype2gtype.NewConver()
	var m kfkrpc.KFKMessage
	t := reflect.TypeOf(m)
	typeStruct, _ := newConver.Conver(t)

	structs := map[string]*types.Struct{
		t.Name(): typeStruct.(*types.Struct),
	}

	doc, err := parser.ParseData("kfkrpc", nil, structs, nil)
	util.CheckPanic(err)

	jsonData, err2 := json.MarshalIndent(doc, "", " ")
	util.CheckPanic(err2)

	code, err3 := gosource.Gen(jsonData, false)
	util.CheckPanic(err3)

	f, e := os.Create(filename)
	util.CheckPanic(e)
	f.Write(code)
	f.Close()
}

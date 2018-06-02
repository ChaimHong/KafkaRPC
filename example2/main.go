package main

import (
	"flag"
)

func main() {
	isSrv := false

	flag.BoolVar(&isSrv, "is_server", true, "is server")

	flag.Parse()

	if isSrv {
		Server()
	} else {
		Client()
	}
}

package main

import (
	"fmt"
	"net"
)

func ListenOnSelfPort(ln net.Listener) {
	for {
		connection, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Print(connection)
	}
}
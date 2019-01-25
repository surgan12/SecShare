package clientProperties

import (
	"fmt"
	"net"
)

type ClientListen struct {
	List    []string
	PeerIP map[string]string
}

type ClientQuery struct {
	Name  []byte
	Query []byte
}

func ListenOnSelfPort(ln net.Listener) {
	for {
		connection, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Print(connection)
	}
}

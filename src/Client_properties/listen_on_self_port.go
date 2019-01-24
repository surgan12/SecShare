package Client_properties

import (
	"fmt"
	"net"
)

type Client_listen struct {
	List    []string
	Peer_IP map[string]string
}

type Client_Query struct {
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

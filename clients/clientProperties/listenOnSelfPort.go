package clientproperties

import (
	"fmt"
	"net"
)

// ClientListen stores list of clients and map of their IP
type ClientListen struct {
	List    []string
	PeerIP map[string]string
}

// ClientQuery stores name and query of clients
type ClientQuery struct {
	Name  []byte
	Query []byte
}

// ListenOnSelfPort listens for clients on network
func ListenOnSelfPort(ln net.Listener) {
	for {
		connection, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Print(connection)
	}
}

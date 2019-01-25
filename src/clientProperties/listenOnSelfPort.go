package clientproperties

import (
	"fmt"
	"net"
)

//struct to store list of clients and their IP addresses
type ClientListen struct {
	List    []string
	PeerIP map[string]string
}

//struct to store the name and query of clients
type ClientQuery struct {
	Name  []byte
	Query []byte
}

//function to listen on port for all the peers among the network
func ListenOnSelfPort(ln net.Listener) {
	for {
		connection, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Print(connection)
	}
}

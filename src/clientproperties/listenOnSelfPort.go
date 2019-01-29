package clientproperties

import (
	"fmt"
	"net"
)

// Client properties as stored in the server
type Client struct {
	Address          string
	Name             string
	ConnectionServer net.Conn
}

// ClientListen stores list of clients and map of their IP
type ClientListen struct {
	List   []string
	PeerIP map[string]string
	PeerListenPort map[string]string
}

// ClientQuery stores name and query of clients
type ClientQuery struct {
	Name  []byte
	Query []byte
	ClientListenPort []byte
}

// ClientJob stores the names, jobs and connection
type ClientJob struct {
	Name  string
	Query string
	Conn  net.Conn
	ClientListenPort string
}

// ListenOnSelfPort listens for clients on network
func ListenOnSelfPort(ln net.Listener) {
	for {
		connection, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Print(connection)
		fmt.Println("connected from client")

	}
}

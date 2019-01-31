package clientproperties

import (
	"fmt"
	"net"
)

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

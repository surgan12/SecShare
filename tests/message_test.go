package main

import (
	"fmt"
	"net"
	"testing"
	// cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	cp "../src/clientproperties"
)

func TestRequestMessage(t *testing.T) {

	ln, err := net.Listen("tcp", ":40000")
	if err != nil {
		fmt.Println("Error ", err, " in listening on address ", ln)
	}

	var list []string
	TestMapPeerIP := make(map[string]string)
	TestMapPeerListenPort := make(map[string]string)

	list = append(list, "abc")
	TestMapPeerIP["abc"] = "127.0.0.1"
	TestMapPeerListenPort["abc"] = "40000"
	activeClient := cp.ClientListen{List: list, PeerIP: TestMapPeerIP, PeerListenPort: TestMapPeerListenPort}

	name := "my_name"

	message_status := cp.RequestMessage(&activeClient, name, "abc", "hey buddy!")
	if message_status != "sent" {
		t.Fatal("Error in sending message ...")
	}

}
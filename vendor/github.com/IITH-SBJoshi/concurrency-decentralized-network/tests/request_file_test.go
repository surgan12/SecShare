package main

import (
	"fmt"
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	"net"
	"testing"
	// cp "../src/clientproperties"
)

func TestRequestSomeFile(t *testing.T) {
	ln, err := net.Listen("tcp", ":50000")
	if err != nil {
		fmt.Println("Error ", err, " in listening on address ", ln)
	}

	var list []string
	TestMapPeerIP := make(map[string]string)
	TestMapPeerListenPort := make(map[string]string)

	list = append(list, "abc")
	TestMapPeerIP["abc"] = "127.0.0.1"
	TestMapPeerListenPort["abc"] = "50000"
	activeClient := cp.ClientListen{List: list, PeerIP: TestMapPeerIP, PeerListenPort: TestMapPeerListenPort}

	name := "requester"
	// var TestFilesInDir []string
	// TestFilesInDir = append(TestFilesInDir, "image.jpg")
	// directoryFiles := cp.ClientFiles{FilesInDir: TestFilesInDir}

	request_status := cp.GetRequestedFile(&activeClient, name, "abc", "image.jpg")

	if request_status == "error_no_file" {
		t.Fatal("Error in requesting file ...")
	}
}
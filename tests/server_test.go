package main

import "testing"

import cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientProperties"

import sp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/ServerProperties"

func TestRemoveFromClient(t *testing.T) {

	clientList := []cp.Client{}
	cli := cp.Client{Address: "001", Name: "dummy"}
	clientList = append(clientList, cli)
	clientList = sp.RemoveFromClient(clientList, "dummy")
	if len(clientList) > 0 {
		t.Fatal("TestRemoveFromClient not working correctly")
	}
}

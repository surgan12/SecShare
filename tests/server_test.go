package main

import "testing"

import cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"

//import cp "../src/clientproperties"

import sp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/serverproperties"

//import sp "../src/serverproperties"

func TestRemoveFromClient(t *testing.T) {

	clientList := []cp.Client{}
	cli := cp.Client{Address: "001", Name: "dummy"}
	clientList = append(clientList, cli)
	clientList = sp.RemoveFromClient(clientList, "dummy")
	if len(clientList) > 0 {
		t.Fatal("RemoveFromClient not working correctly")
	}
}

func TestQueryDeal(t *testing.T) {
	name := "user"
	TestMap := make(map[string]string)
	TestMap[name] = "active"
	var list []string
	list = append(list, "user")
	cli := cp.ClientListen{List: list, PeerIP: TestMap}
	clt := cp.Client{Name: name, Address: ":8087"}
	var clients = []cp.Client{clt}
	cli = sp.QueryDeal(&clients, cli, name)
	if len(clients) > 0 || len(cli.List) > 0 || cli.PeerIP[name] != "" {
		t.Fatal("QueryDeal in the serverproperties not working correctly")
	}

}

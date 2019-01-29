package serverproperties

import (
	// cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	cp "../clientproperties"
)

// RemoveFromClient removes the client who quits from the list
func RemoveFromClient(clients []cp.Client, name string) []cp.Client {
	tempClients := []cp.Client{}
	for i := 0; i < len(clients); i++ {
		if clients[i].Name != name {
			tempClients = append(tempClients, clients[i])
		}
	}
	return tempClients
}

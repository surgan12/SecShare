package clientproperties

import (
	"../encryptionproperties"
	"encoding/json"
	// "fmt"
	"net"
)

//MyPeers list of connections dialed by current client
type MyPeers struct {
	Conn 	net.Conn
	PeerName string	
}

// FileRequest stores the queries and information about requester
type FileRequest struct {
	query         string
	myAddress     string
	myName        string
	requestedFile string
}

// sendingToServer function to send queries to server
func sendingToServer(name []byte, query []byte, conn net.Conn, 
					 queryType string, listenPort []byte) {
	objectToSend := ClientQuery{Name: name, Query: query, ClientListenPort: listenPort}
	encoder := json.NewEncoder(conn)
	encoder.Encode(objectToSend)
	if queryType == "quit" {
		conn.Close()
	}
}

package clientproperties

import (
	"encoding/json"
	// "fmt"
	// fp "../../fileproperties"
	// "crypto/rand"
	// "crypto/rsa"
	// "crypto/sha512"	
	fp "github.com/IITH-SBJoshi/concurrency-decentralized-network/fileproperties"
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
	List           []string
	PeerIP         map[string]string
	PeerListenPort map[string]string
}

// ClientQuery stores name and query of clients
type ClientQuery struct {
	Name             []byte
	Query            []byte
	ClientListenPort []byte
}

// ClientJob stores the names, jobs and connection
type ClientJob struct {
	Name             string
	Query            string
	Conn             net.Conn
	ClientListenPort string
}

//MyPeers list of connections dialed by current client
type MyPeers struct {
	Conn     net.Conn
	PeerName string
}
//MyReceivedFiles received files
type MyReceivedFiles struct {
	MyFileName string
	MyFile     []FilePartContents
	FilePartInfo fp.FilePartInfo
}

//FilePartContents contents of file in parts
type FilePartContents struct {
	Contents []byte
}
//BaseRequest request for file
type BaseRequest struct {
	RequestType string
	FileRequest 
	FilePartInfo fp.FilePartInfo
}

// FileRequest stores the queries and information about requester
type FileRequest struct {
	Query         string
	MyAddress     string
	MyName        string
	RequestedFile string
}

//SendingToServer function to send queries to server
func SendingToServer(name []byte, query []byte, conn net.Conn,
	queryType string, listenPort []byte) {

	objectToSend := ClientQuery{Name: name, Query: query, ClientListenPort: listenPort}
	encoder := json.NewEncoder(conn)
	encoder.Encode(objectToSend)
	if queryType == "quit" {
		conn.Close()
	}
}


package clientproperties

import (
	// fp "../../fileproperties"
	"encoding/json"
	"fmt"
	"sync"
	// "crypto/rand"
	// "crypto/rsa"
	// "crypto/sha512"
	fp "github.com/IITH-SBJoshi/concurrency-decentralized-network/fileproperties"
	"net"
)

// Struct to store details of specific client
type Client struct {
	Address          string
	Name             string
	ConnectionServer net.Conn
}

// Store - all client names,
//		   IP to name mapping of all clients
//		   Port at which all those clients are listening for P2P requests
type ClientListen struct {
	List           []string
	PeerIP         map[string]string
	PeerListenPort map[string]string
}

// To store name and query of a client
type ClientQuery struct {
	Name             []byte
	Query            []byte
	ClientListenPort []byte
}

// Stores the names, jobs, query and connection
type ClientJob struct {
	Name             string
	Query            string
	Conn             net.Conn
	ClientListenPort string
}

// List of peers to which a clinet dials
type MyPeers struct {
	Conn     net.Conn
	PeerName string
}

// To store the information of files received by client
type MyReceivedFiles struct {
	PartsReceived int
	MyFileName   string             // name of file
	MyFile       []FilePartContents // Contents of the file
	FilePartInfo fp.FilePartInfo    // Information of various file parts
}

// To store the information of the messages receievd
type MyReceivedMessages struct {
	Counter    int              // The counter from which client has to start reading the messages
	MyMessages []MessageRequest // slice of structs of type MessageRequest to store all message requests
}

// Contents of a part of file
type FilePartContents struct {
	Contents []byte
}

// A base request which is used as a generic request for all types of P2P queries
type BaseRequest struct {
	RequestType    string          // type of request
	FileRequest                    // Information about File requseter if its a file request
	FilePartInfo   fp.FilePartInfo // Information about file parts if its a file request
	MessageRequest                 // Details of message is its a message request
}

// FileRequest stores the queries and information about requester
type FileRequest struct {
	Query         string
	MyAddress     string
	MyName        string
	RequestedFile string
}

// Stores the information about requester, who is sending the message (for the receiver to reply back)
type MessageRequest struct {
	SenderQuery   string
	SenderAddress string
	SenderName    string
	Message       string
}

// ClientFiles stores the files in the "files" directory of client
type ClientFiles struct {
	FilesInDir []string
}

// SendingToServer - to send queries to server
func SendingToServer(name []byte, query []byte, conn net.Conn,
	queryType string, listenPort []byte) {

	objectToSend := ClientQuery{Name: name, Query: query, ClientListenPort: listenPort}
	encoder := json.NewEncoder(conn)
	encoder.Encode(objectToSend)
	if queryType == "quit" {
		conn.Close()
	}
}

// To display the most recent messages which haven't been seen yet
func DisplayRecentUnseenMessages(mymessages *MyReceivedMessages) {
	// locking it, so that new messages can't be written at the current moment
	var mutex = &sync.Mutex{}
	mutex.Lock() // locking
	
	if mymessages.Counter == len(mymessages.MyMessages) {
		fmt.Println("No Recent Unseen messages!!")
	} else {
		for i := mymessages.Counter; i < len(mymessages.MyMessages); i++ {
			fmt.Println(mymessages.MyMessages[i].SenderName, " - sent you a message : ", mymessages.MyMessages[i].Message)
		}
		mymessages.Counter = len(mymessages.MyMessages) // incrementing the counter to latest count, as we have read all recent messages
	}
	
	fmt.Println('\n')
	mutex.Unlock()
}

// To display Num recent messages
func DisplayNumRecentMessages(mymessages *MyReceivedMessages, recentCount int) {
	// locking it, so that new messages can't be written at the current moment
	var mutex = &sync.Mutex{}

	// if contain more than recentCount number of messages
	if len(mymessages.MyMessages)-recentCount >= 0 {
		mutex.Lock() // locking
		for i := len(mymessages.MyMessages) - recentCount; i < len(mymessages.MyMessages); i++ {
			fmt.Println(mymessages.MyMessages[i].SenderName, " - sent you a message : ", mymessages.MyMessages[i].Message)
		}
		fmt.Println('\n')
		mutex.Unlock()

	} else {
		// if contains less than recentCount number of messages
		fmt.Println("Displaying all messages!")
		mutex.Lock() // locking
		for i := 0; i < len(mymessages.MyMessages); i++ {
			fmt.Println(mymessages.MyMessages[i].SenderName, " - sent you a message : ", mymessages.MyMessages[i].Message)
		}
		fmt.Println('\n')
		mutex.Unlock()
	}

}

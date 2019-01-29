package clientproperties

import (
	"encoding/json"
	"fmt"
	"net"
	// client "../../clients"
)

// FileRequest stores the queries and information about requester
type FileRequest struct {
	query         string
	myAddress     string
	myName        string
	requestedFile string
}

// RequestSomeFile request files from peers on network
func RequestSomeFile(activeClient ClientListen, name string) {
	var senderName string // is the person who will send the file
	fmt.Println("Whom do you want to receive the file from ? : ")
	fmt.Scanln(&senderName)
	var fileName string
	fmt.Println("What file do you want ? ")
	fmt.Scanln(&fileName) // file we want to receive

	fileRequest := FileRequest{query: "receive_file", myAddress: activeClient.PeerIP[name],
		myName: name, requestedFile: "any song"}	

	fmt.Println("receive from ", activeClient.PeerListenPort[senderName])
	connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[senderName])
	fmt.Println(err)
	// for err != nil {
	// 	fmt.Println("Please enter a valid person name - ")
	// 	connection1, err1 := net.Dial("tcp", activeClient.PeerListenPort[senderName])
	// 	connection = connection1
	// 	err = err1
	// }

	encoder := json.NewEncoder(connection)
	encoder.Encode(fileRequest)
}
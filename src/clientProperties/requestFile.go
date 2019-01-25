package clientProperties

import (
	"encoding/json"
	"fmt"
	"net"
	// client "../../clients"
)

type FileRequest struct {
	query          string
	myAddress     string
	myName        string
	requestedFile string
}

func RequestSomeFile(activeClient ClientListen, name string) {
	var senderName string // is the person who will send the file
	fmt.Println("Whom do you want to receive the file from ? : ")
	fmt.Scanln(&senderName)
	var fileName string
	fmt.Println("What file do you want ? ")
	fmt.Scanln(&fileName) // file we want to receive

	fileRequest := FileRequest{query: "receive_file", myAddress: activeClient.PeerIP[name],
		myName: name, requestedFile: "any song"}

	connection, err := net.Dial("tcp", activeClient.PeerIP[senderName])
	
	for err != nil {
		fmt.Println("Please enter a valid person name - ")
	}

	encoder := json.NewEncoder(connection)
	encoder.Encode(fileRequest)
}

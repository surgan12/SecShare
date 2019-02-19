package clientproperties

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func sendFileRequestToPeer(connection net.Conn, fileRequest FileRequest) {
	//handle with care, FilePartInfo field for this truct is Nil. Will throw seg fault if accessed
	
	someRequest := BaseRequest{ RequestType: "receive_from_peer", FileRequest: fileRequest}
	encoder1 := json.NewEncoder(connection)
	encoder1.Encode(&someRequest)
}

// RequestSomeFile request files from peers on network
func RequestSomeFile(activeClient *ClientListen, name string, directoryFiles *ClientFiles) {
	// _, PublicKeyClient = GenerateKeyPair()
	var fileExist bool
	// fmt.Print(len(activeClient.List))
	var fileSenderName string // is the person who will send the file
	fmt.Print("Whom do you want to receive the file from ? : ")
	fmt.Scanln(&fileSenderName)
	var fileName string
	fmt.Print("What file do you want ? ")
	fmt.Scanln(&fileName) // file we want to receive
	for _, file := range directoryFiles.FilesInDir {

		if(strings.Compare(file, fileName) == 0){
			fileExist = true;
		}
	}
	if(fileExist == true){

		fileRequest := FileRequest{Query: "receive_file", 
				   		MyAddress: activeClient.PeerIP[name],
				   		MyName: name, RequestedFile: fileName}

		connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[fileSenderName])
		for err != nil {
			fmt.Println("Please enter a valid person name - ")
			connection1, err1 := net.Dial("tcp", ":" + activeClient.PeerListenPort[fileSenderName])
			connection = connection1
			err = err1
		}

		sendFileRequestToPeer(connection, fileRequest)
	} else {
		fmt.Println("files does not exist")
	}	

	// connection.Close()	// closing connection after one time requestb17e198f6aeb5753c2c193c
}

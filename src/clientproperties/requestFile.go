package clientproperties

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func SendFileRequestToPeer(connection net.Conn, fileRequest FileRequest) {
	//handle with care, FilePartInfo field for this struct is Nil. Will throw seg fault if accessed

	someRequest := BaseRequest{RequestType: "receive_from_peer", FileRequest: fileRequest}
	encoder1 := json.NewEncoder(connection)
	encoder1.Encode(&someRequest)
}

// RequestSomeFile function request files from peers on network
func RequestSomeFile(activeClient *ClientListen, name string, directoryFiles *ClientFiles, fileSenderName string, fileName string) string {
	// _, PublicKeyClient = GenerateKeyPair()

	var fileExist bool

	for _, file := range directoryFiles.FilesInDir {

		if strings.Compare(file, fileName) == 0 {
			fileExist = true
		}
	}

	if fileExist == true {

		fileRequest := FileRequest{Query: "receive_file",
			MyAddress: activeClient.PeerIP[name],
			MyName:    name, RequestedFile: fileName}

		connection, err := net.Dial("tcp", ":"+activeClient.PeerListenPort[fileSenderName])
		for err != nil {
			fmt.Println("Please enter a valid person name - ")
			connection1, err1 := net.Dial("tcp", ":"+activeClient.PeerListenPort[fileSenderName])
			connection = connection1
			err = err1
		}

		SendFileRequestToPeer(connection, fileRequest)
		request_status := "completed"
		return request_status

	} else {
		fmt.Println("files does not exist")
		request_status := "error_no_file"
		return request_status
	}

	// connection.Close()	// closing connection after one time requestb17e198f6aeb5753c2c193c
}

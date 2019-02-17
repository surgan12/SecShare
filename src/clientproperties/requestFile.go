package clientproperties

import (
	"encoding/json"
	"fmt"
	"net"
	// "crypto/rand"
	// "crypto/rsa"
	// "crypto/sha512"
)

func sendFileRequestToPeer(connection net.Conn, fileRequest FileRequest) {
	//handle with care, FilePartInfo field for this truct is Nil. Will throw seg fault if accessed
	// fmt.Println("file properties testing")
	// fmt.Println(FileRequest.myName)
	// fmt.Println(FileRequest.query)	
	someRequest := BaseRequest{ RequestType: "receive_from_peer", FileRequest: fileRequest}
	encoder1 := json.NewEncoder(connection)
	encoder1.Encode(&someRequest)
}

// RequestSomeFile request files from peers on network
func RequestSomeFile(activeClient *ClientListen, name string) {
	// _, PublicKeyClient = GenerateKeyPair()

	var fileSenderName string // is the person who will send the file
	fmt.Println("Whom do you want to receive the file from ? : ")
	fmt.Scanln(&fileSenderName)
	var fileName string
	fmt.Println("What file do you want ? ")
	fmt.Scanln(&fileName) // file we want to receive

	fileRequest := FileRequest{Query: "receive_file", 
				   MyAddress: activeClient.PeerIP[name],
				   MyName: name, RequestedFile: fileName}	

	// if !checkPeers(myPeers, fileSenderName) {
	// 	connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[fileSenderName])
	// 	for err != nil {
	// 		fmt.Println("Please enter a valid person name - ")
	// 		connection1, err1 := net.Dial("tcp", ":" + activeClient.PeerListenPort[fileSenderName])
	// 		connection = connection1
	// 		err = err1
	// 	}
	// 	currentPeer := MyPeers{conn: connection, PeerName : fileSenderName}
	// 	myPeers := append(myPeers, currentPeer)
	// 	connection.Write([]byte(name))
	// }

	connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[fileSenderName])
	for err != nil {
		fmt.Println("Please enter a valid person name - ")
		connection1, err1 := net.Dial("tcp", ":" + activeClient.PeerListenPort[fileSenderName])
		connection = connection1
		err = err1
	}

	sendFileRequestToPeer(connection, fileRequest)
	// connection.Close()	// closing connection after one time requestb17e198f6aeb5753c2c193c
}

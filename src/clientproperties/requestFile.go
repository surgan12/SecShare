package clientproperties

import (
	"encoding/json"
	"fmt"
	"net"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

func checkPeers (myPeers []MyPeers, checkName string) bool {
	for i := 0; i < len(myPeers); i++ {
		if(checkName == myPeers[i].PeerName)
			return true
	}
	return false 
}

// RequestSomeFile request files from peers on network
func RequestSomeFile(activeClient ClientListen, name string, myPeers []MyPeers) {
	_, PublicKeyClient = GenerateKeyPair()

	var senderName string // is the person who will send the file
	fmt.Println("Whom do you want to receive the file from ? : ")
	fmt.Scanln(&senderName)
	var fileName string
	fmt.Println("What file do you want ? ")
	fmt.Scanln(&fileName) // file we want to receive

	fileRequest := FileRequest{query: "receive_file", myAddress: activeClient.PeerIP[name],
		myName: name, requestedFile: "any song"}	

	if !checkPeers(myPeers, senderName) {
		connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[senderName])
		for err != nil {
			fmt.Println("Please enter a valid person name - ")
			connection1, err1 := net.Dial("tcp", activeClient.PeerListenPort[senderName])
			connection = connection1
			err = err1
		}
		currentPeer := MyPeers{conn: connection, PeerName : senderName} 
	}

	serverkeysClient = PerformHandshake(connection, PublicKeyClient)

	//Encryptions of file has been done here
	senderNameQuery := EncryptWithPublicKey([]byte(senderName), serverkeysClient)
	fileNameQuery := EncryptWithPublicKey([]byte(fileName), serverkeysClient)


	// fmt.Println(err)

	encoder := json.NewEncoder(connection)
	encoder.Encode(fileRequest)
}

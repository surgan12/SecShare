package clientproperties

import (
	"net"
	en "../encryptionproperties"
	"encoding/json"
	"fmt"
)

func checkPeers (myPeers []MyPeers, checkName string) bool {
	for i := 0; i < len(myPeers); i++ {
		if checkName == myPeers[i].PeerName {
			return true
		}
	}
	return false 
}

// RequestSomeFile request files from peers on network
func RequestSomeFile(activeClient ClientListen, name string, myPeers []MyPeers) {
	_, PublicKeyClient := en.GenerateKeyPair()

	var senderName string // is the person who will send the file
	fmt.Println("Whom do you want to receive the file from ? : ")
	fmt.Scanln(&senderName)
	var fileName string
	fmt.Println("What file do you want ? ")
	fmt.Scanln(&fileName) // file we want to receive

	fileRequest := FileRequest{query: "receive_file", myAddress: activeClient.PeerIP[name],
		myName: name, requestedFile: "any song"}

	fmt.Println("Value of checkPeers array is ", checkPeers(myPeers, senderName))	

	if !checkPeers(myPeers, senderName) {

		connection , err := net.Dial("tcp", ":" + activeClient.PeerListenPort[senderName])
		for err != nil {
			fmt.Println("Please enter a valid person name - ")
			connection1, err1 := net.Dial("tcp", activeClient.PeerListenPort[senderName])
			connection = connection1
			err = err1
		}
		currentPeer := MyPeers{Conn: connection , PeerName : senderName} 
	}

	serverkeysClient := en.PerformHandshake(connection, PublicKeyClient)

	//Encryptions of file has been done here
	senderNameQuery := en.EncryptWithPublicKey([]byte(senderName), serverkeysClient)
	fileNameQuery := en.EncryptWithPublicKey([]byte(fileName), serverkeysClient)


	// fmt.Println(err)

	encoder := json.NewEncoder(connection)
	encoder.Encode(fileRequest)
}

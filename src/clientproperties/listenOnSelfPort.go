package clientproperties

import (
	"fmt"
	"encoding/json"
	"net"
	fp "../../fileproperties"
	"sync"
)

//Check if the current guy is already my peer
// func checkPeers (myPeers []MyPeers, checkName string) bool {
// 	for i := 0; i < len(myPeers); i++ {
// 		if(checkName == myPeers[i].PeerName)
// 			return true
// 	}
// 	return false 
// }

var mutex = &sync.Mutex{} // Lock and unlock (Mutex)

func sendFileParts(newfilerequest FileRequest, allfileparts []fp.FilePartInfo, 
				   activeClient ClientListen, myname string) {
	for names, _ := range activeClient.PeerListenPort {
		if (names != myname){
			connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[names])
			for err != nil {
				fmt.Println("Error in dialing, dialing again ... ")
				connection1, err1 := net.Dial("tcp", ":" + activeClient.PeerListenPort[names])
				connection = connection1
				err = err1
			}
			fmt.Println("Connection established to send a file part to connection - ", connection)
		}
	}
}

func handleNewFileSendRequest(newfilerequest FileRequest, myname string, activeClient ClientListen) {
	if newfilerequest.myName == myname {
		
		allfileparts := fp.GetSplitFile(newfilerequest.requestedFile)
		sendFileParts(newfilerequest, allfileparts, activeClient, myname)

	} else {

		fmt.Println("Forwarding to file receiver")

	}
}

func handleReceivedFile(newrequest *BaseRequest, myfiles map[string]MyReceivedFiles) {

	var TotalFileParts int
	var filePartNum int 
	requestedFileName := newrequest.FilePartInfo.FileName
	if _, ok := myfiles[requestedFileName]; ok {
		// appending to already created object of this received file
		TotalFileParts = newrequest.FilePartInfo.TotalParts
		filePartNum = newrequest.FilePartInfo.PartNumber
    	
    	mutex.Lock()
    	myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
    	mutex.Unlock()

    	if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
    		concatenateFileParts(myfiles[requestedFileName])
    	}

	} else {
		// creating new received file object for my own file
		myfiles[requestedFileName] = MyReceivedFiles{}
		// myfiles[requestedFileName].MyFile := make([]FilePartContents, newrequest.FilePartInfo.TotalParts)
		// myfiles[requestedFileName].MyFileName = newrequest.FilePartInfo.FileName

		mutex.Lock()
		myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
    	mutex.Unlock()

    	if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
    		concatenateFileParts(myfiles[requestedFileName])
    	}		
	}
}

func handleReceivedRequest(connection net.Conn, activeClient ClientListen, myname string, 
						   myfiles map[string]MyReceivedFiles) {
	var newrequest BaseRequest
	newconn := json.NewDecoder(connection)
	newconn.Decode(newrequest)

	if newrequest.FileRequest.myName == myname {

		fmt.Println("Received some file part for myself ")
		handleReceivedFile(&newrequest, myfiles)

	} else {

		fmt.Println("Forwarding some received file part ")
		
		// myAddress is address of person asking for the file
		receiverAddress := newrequest.FileRequest.myAddress
		forwardConnection, forwardConnErr := net.Dial("tcp", receiverAddress)
		for forwardConnErr != nil {
			fmt.Println("Error in dialing, dialing again ... ")
			connection1, err1 := net.Dial("tcp", receiverAddress)
			forwardConnection = connection1
			forwardConnErr = err1
		}

		newSendRequest := newrequest
		newconn := json.NewEncoder(forwardConnection)
		newconn.Encode(newSendRequest)

	}
}

func handleConnection(connection net.Conn, activeClient ClientListen, myname string, 
				      myfiles map[string]MyReceivedFiles) {

	var newrequest BaseRequest
	newconn := json.NewDecoder(connection)
	newconn.Decode(newrequest)

	if newrequest.RequestType == "receive_from_peer" {

		fmt.Println("Request to receive a file from peer ")
		handleNewFileSendRequest(newrequest.FileRequest, myname, activeClient)

	} else if newrequest.RequestType == "received_some_file" {

		fmt.Println("Received some file part ")
		handleReceivedRequest(connection, activeClient, myname, myfiles)
	}

}

// ListenOnSelfPort listens for clients on network
func ListenOnSelfPort(ln net.Listener, myname string, activeClient ClientListen, 
					   myfiles map[string]MyReceivedFiles) {
	for {
		connection, err := ln.Accept()
		
		if err != nil {
			panic(err)
		}

		fmt.Print(connection)
		fmt.Println("connection received from client")

		go handleConnection(connection, activeClient, myname, myfiles)
	}
}
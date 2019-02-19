package clientproperties

import (
	"fmt"
	"encoding/json"
	"net"
	fp "../../fileproperties"
	// fp "github.com/IITH-SBJoshi/concurrency-decentralized-network/fileproperties"
	"sync"
)

var mutex = &sync.Mutex{} // Lock and unlock (Mutex)

func sendFileParts(newfilerequest FileRequest, allfileparts []fp.FilePartInfo, 
				   activeClient *ClientListen, myname string) {
	for names := range activeClient.PeerListenPort {
		// fmt.Println("cureent names : ", names)
		if (names != myname){
			// fmt.Println("cureent names : ", names)
			connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[names])
			for err != nil {
				fmt.Println("Error in dialing, dialing again ... ")
				connection1, err1 := net.Dial("tcp", ":" + activeClient.PeerListenPort[names])
				connection = connection1
				err = err1
			}
			fmt.Println("Connection established to send a file part to connection - ", connection)
			baseRequest := BaseRequest {RequestType : "received_some_file", FileRequest : newfilerequest,
										 FilePartInfo : allfileparts[0]}
			encoder := json.NewEncoder(connection)
			encoder.Encode(&baseRequest)

		}
	}
}

func handleNewFileSendRequest(newfilerequest FileRequest, myname string, activeClient *ClientListen) {
	// fmt.Println(newfilerequest.MyName)
	// fmt.Println(myname)
		
	allfileparts := fp.GetSplitFile(newfilerequest.RequestedFile, len(activeClient.List))
	// fmt.Println("received file from Happy")
	sendFileParts(newfilerequest, allfileparts, activeClient, myname)

}

func handleReceivedFile(newrequest BaseRequest, myfiles map[string]MyReceivedFiles) {

	var TotalFileParts int
	var filePartNum int 
	// fmt.Println("testing : 2")
	requestedFileName := newrequest.FilePartInfo.FileName
	// fmt.Println(newrequest.FilePartInfo.TotalParts)
	// fmt.Println(newrequest.FilePartInfo.PartNumber)
	// fmt.Println(requestedFileName)

	TotalFileParts = newrequest.FilePartInfo.TotalParts
	filePartNum = newrequest.FilePartInfo.PartNumber
	if _, ok := myfiles[requestedFileName]; ok {

    	myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
    	// mutex.Unlock()
		// fmt.Println("1")
    	if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
    		concatenateFileParts(myfiles[requestedFileName])
    	}
		// fmt.Println("1")
	} else {
		// creating new received file object for my own file
		myfiles[requestedFileName] = MyReceivedFiles{newrequest.FilePartInfo.FileName,
										 make([]FilePartContents, newrequest.FilePartInfo.TotalParts),
										 newrequest.FilePartInfo }
		mutex.Lock()
		myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
    	mutex.Unlock()
    	fmt.Println(myfiles[requestedFileName].MyFileName)
    	fmt.Println(len(myfiles[requestedFileName].MyFile))
    	if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
    		concatenateFileParts(myfiles[requestedFileName])
    	}
	}
}

func handleReceivedRequest(connection net.Conn, activeClient *ClientListen, myname string, 
						   myfiles map[string]MyReceivedFiles, newrequest BaseRequest) {
	if newrequest.FileRequest.MyName == myname {

		fmt.Println("Received some file part for myself ")
		
		handleReceivedFile(newrequest, myfiles)

	} else {

		fmt.Println("Forwarding some received file part ")
		
		// myAddress is address of person asking for the file
		receiverAddress := newrequest.FileRequest.MyAddress
		forwardConnection, forwardConnErr := net.Dial("tcp", receiverAddress)
		for forwardConnErr != nil {
			fmt.Println("Error in dialing, dialing again ... ")
			connection1, err1 := net.Dial("tcp", receiverAddress)
			forwardConnection = connection1
			forwardConnErr = err1
		}

		newSendRequest := newrequest
		newconn := json.NewEncoder(forwardConnection)
		newconn.Encode(&newSendRequest)

	}
}

func handleConnection(connection net.Conn, activeClient *ClientListen, myname string, 
				      myfiles map[string]MyReceivedFiles) {

	var newrequest BaseRequest
	newconn := json.NewDecoder(connection)
	newconn.Decode(&newrequest)

	if newrequest.RequestType == "receive_from_peer" {

		fmt.Println("Request to receive a file from peer ")
		handleNewFileSendRequest(newrequest.FileRequest, myname, activeClient)

	} else if newrequest.RequestType == "received_some_file" {

		fmt.Println("Received some file part ")
		handleReceivedRequest(connection, activeClient, myname, myfiles, newrequest)
	}

}

// ListenOnSelfPort listens for clients on network
func ListenOnSelfPort(ln net.Listener, myname string, activeClient *ClientListen, 
					   myfiles map[string]MyReceivedFiles) {
	for {
		connection, err := ln.Accept()
		
		if err != nil {
			panic(err)
		}
		// fmt.Print(connection)
		// fmt.Println(activeClient.PeerListenPort)
		// fmt.Println(activeClient.List)
		go handleConnection(connection, activeClient, myname, myfiles)
	}
}
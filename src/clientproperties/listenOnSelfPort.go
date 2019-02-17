package clientproperties

import (
	"fmt"
	"encoding/json"
	"net"
	// fp "../../fileproperties"
	fp "github.com/IITH-SBJoshi/concurrency-decentralized-network/fileproperties"
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
				   activeClient *ClientListen, myname string) {
	fmt.Println("send file parst m")
	fmt.Println(activeClient.PeerListenPort)
	for names := range activeClient.PeerListenPort {
		fmt.Println("cureent names : ", names)
		if (names != myname){
			fmt.Println("cureent names : ", names)
			connection, err := net.Dial("tcp", ":" + activeClient.PeerListenPort[names])
			for err != nil {
				fmt.Println("Error in dialing, dialing again ... ")
				connection1, err1 := net.Dial("tcp", ":" + activeClient.PeerListenPort[names])
				connection = connection1
				err = err1
			}
			fmt.Println("Connection established to send a file part to connection - ", connection)
			// fileRequest := FileRequest { Query : "receive_file", 
			// 							MyAddress : activeClient.PeerListenPort[myname],
			// 							MyName : myname, RequestedFile : filename}
			baseRequest := BaseRequest {RequestType : "received_some_file", FileRequest : newfilerequest,
										 FilePartInfo : allfileparts[0]}
			// fileParts := allfileparts
			encoder := json.NewEncoder(connection)
			encoder.Encode(&baseRequest)

		}
	}
}

func handleNewFileSendRequest(newfilerequest FileRequest, myname string, activeClient *ClientListen) {
	fmt.Println(newfilerequest.MyName)
	fmt.Println(myname)
	// if newfilerequest.MyName == myname {
		
	allfileparts := fp.GetSplitFile(newfilerequest.RequestedFile)
	fmt.Println("received file from Happy")
	sendFileParts(newfilerequest, allfileparts, activeClient, myname)

	// } else {

	// 	fmt.Println("Forwarding to file receiver")

	// }
}

func handleReceivedFile(newrequest BaseRequest, myfiles map[string]MyReceivedFiles) {

	var TotalFileParts int
	var filePartNum int 
	fmt.Println("testing : 2")
	requestedFileName := newrequest.FilePartInfo.FileName
	fmt.Println(newrequest.FilePartInfo.TotalParts)
	fmt.Println(newrequest.FilePartInfo.PartNumber)
	fmt.Println(requestedFileName)

	TotalFileParts = newrequest.FilePartInfo.TotalParts
	filePartNum = newrequest.FilePartInfo.PartNumber
	// fmt.Println(myfiles)
	// myfiles := make(map[string]MyReceivedFiles)
	// myfiles[requestedFileName] := MyReceivedFiles
	if _, ok := myfiles[requestedFileName]; ok {
		// appending to already created object of this received file
		// fmt.Println("1")
		// TotalFileParts = newrequest.FilePartInfo.TotalParts
		// filePartNum = newrequest.FilePartInfo.PartNumber
  //  		fmt.Println("1")

    	// mutex.Lock()
    	myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
    	// mutex.Unlock()
		fmt.Println("1")
    	if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
    		concatenateFileParts(myfiles[requestedFileName])
    	}
		fmt.Println("1")
	} else {
		// creating new received file object for my own file
		fmt.Println("1")
		// myfiles[requestedFileName] = MyReceivedFiles{}
		// myfiles[requestedFileName].MyFile = make([]FilePartContents, newrequest.FilePartInfo.TotalParts)
		// myfiles[requestedFileName].MyFileName = newrequest.FilePartInfo.FileName
		// myfiles := make(MyReceivedFiles)
		myfiles[requestedFileName] = MyReceivedFiles{newrequest.FilePartInfo.FileName,
										 make([]FilePartContents, newrequest.FilePartInfo.TotalParts),
										 newrequest.FilePartInfo }
		fmt.Println("2")
		mutex.Lock()
		myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
    	mutex.Unlock()
    	fmt.Println("3")
    	fmt.Println(myfiles[requestedFileName].MyFileName)
    	fmt.Println(len(myfiles[requestedFileName].MyFile))
    	if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
    		fmt.Println("inside if")
    		concatenateFileParts(myfiles[requestedFileName])
    	}
    	fmt.Println("4")		
	}
}

func handleReceivedRequest(connection net.Conn, activeClient *ClientListen, myname string, 
						   myfiles map[string]MyReceivedFiles, newrequest BaseRequest) {
	fmt.Println("testing  : 1")
	// var newrequest BaseRequest
	// newconn := json.NewDecoder(connection)
	// newconn.Decode(&newrequest)
	fmt.Println(newrequest.RequestType)
	fmt.Println(newrequest.FileRequest.MyName)
	fmt.Println(myname)
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
	// fmt.Println("testing : 1")

	var newrequest BaseRequest
	newconn := json.NewDecoder(connection)
	newconn.Decode(&newrequest)
	fmt.Println("request type : ", newrequest.RequestType)

	// var newrequest BaseRequest
	// newconn := json.NewDecoder(connection)
	// newconn.Decode(&newrequest)
	// fmt.Println("request type : ", newrequest.RequestType)


	// fmt.Println(newrequest)
	// fmt.Println(newrequest.FileRequest)
	// fmt.Println(newrequest.FileRequest.Query)
	// fmt.Println(newrequest.FileRequest.MyName)
	// fmt.Println(myname)

	if newrequest.RequestType == "receive_from_peer" {

		fmt.Println("Request to receive a file from peer ")
		// fmt.Println(newrequest.FileRequest.MyName)
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
		// fmt.Println("\n")
		fmt.Print(connection)
		fmt.Println("my name is : " , myname)
		// fmt.Println(myfiles)
		// var newrequest BaseRequest
		// newconn := json.NewDecoder(connection)
		// newconn.Decode(&newrequest)

		// var newRequestfile FileRequest
		// newconn := json.NewDecoder(connection)
		// newconn.Decode(&newRequestfile)

		// fmt.Println(newRequestfile)
		// fmt.Println("connection received from client")
		fmt.Println(activeClient.PeerListenPort)
		fmt.Println(activeClient.List)
		go handleConnection(connection, activeClient, myname, myfiles)
	}
}
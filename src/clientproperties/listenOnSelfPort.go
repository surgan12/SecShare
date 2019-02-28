package clientproperties

import (
	// fp "../../fileproperties"
	"encoding/json"
	"fmt"
	fp "github.com/IITH-SBJoshi/concurrency-decentralized-network/fileproperties"
	"net"
	"strings"
	"sync"
)

// creating locks for messages and files array
var mutexFiles = &sync.Mutex{}    // Lock and unlock (mutexFiles)
var mutexMessages = &sync.Mutex{} // Lock and unlock (mutexFiles)

// send various file parts to peers
func sendFileParts(newfilerequest FileRequest, allfileparts []fp.FilePartInfo,
	activeClient *ClientListen, myname string) int {

	countSent := 0
	for names := range activeClient.PeerListenPort {
		if names != myname {

			count := 0
			connection, err := net.Dial("tcp", ":"+activeClient.PeerListenPort[names])
			for err != nil {
				fmt.Println("Error in dialing to: ", names, " dialing again...")
				connection1, err1 := net.Dial("tcp", ":"+activeClient.PeerListenPort[names])
				connection = connection1
				err = err1
				count++
				if count > 100 {
					fmt.Println("Error in sending current file part - ", err)
				}
			}

			// sending the file part corresponding to this peer
			// fmt.Println("Connection established to send a file part to connection - ", connection)
			baseRequest := BaseRequest{RequestType: "received_some_file", FileRequest: newfilerequest,
				FilePartInfo: allfileparts[countSent]}
			encoder := json.NewEncoder(connection)
			encoder.Encode(&baseRequest)
			countSent++
		}

	}
	return countSent
}

// Handle request to send some file to a peer
func handleNewFileSendRequest(newfilerequest FileRequest, myname string, activeClient *ClientListen) {
	// getting all splits of file
	allfileparts := fp.GetSplitFile(newfilerequest.RequestedFile, len(activeClient.List))
	countSent := sendFileParts(newfilerequest, allfileparts, activeClient, myname)

}

// Handle some received file part for myself
func handleReceivedFile(newrequest BaseRequest, myfiles map[string]MyReceivedFiles) {

	var TotalFileParts int
	var filePartNum int

	requestedFileName := newrequest.FilePartInfo.FileName

	TotalFileParts = newrequest.FilePartInfo.TotalParts
	filePartNum = newrequest.FilePartInfo.PartNumber

	// If already exists in myfies, append the current part to corresponding file struct
	if _, ok := myfiles[requestedFileName]; ok {

		mutexFiles.Lock()
		myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
		mutexFiles.Unlock()
		// if all parts have been received, concatenate it and create new file
		if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
			concatenateFileParts(myfiles[requestedFileName])
		}

	} else {
		// creating new received file object for my own file
		myfiles[requestedFileName] = MyReceivedFiles{newrequest.FilePartInfo.FileName,
			make([]FilePartContents, newrequest.FilePartInfo.TotalParts),
			newrequest.FilePartInfo}

		// locking
		mutexFiles.Lock()
		myfiles[requestedFileName].MyFile[filePartNum].Contents = newrequest.FilePartInfo.FilePartContents
		mutexFiles.Unlock()

		// if all parts have been received, concatenate it and create new file
		if len(myfiles[requestedFileName].MyFile) == TotalFileParts {
			concatenateFileParts(myfiles[requestedFileName])
		}
	}
}

// Handle a request
func handleReceivedRequest(connection net.Conn, activeClient *ClientListen, myname string,
	myfiles map[string]MyReceivedFiles, newrequest BaseRequest) {

	// If the file receievd is for me
	if newrequest.FileRequest.MyName == myname {
		handleReceivedFile(newrequest, myfiles)

	} else {
		// if file is received to be forwarded to some other peer
		receiverAddress := newrequest.FileRequest.MyAddress
		forwardConnection, forwardConnErr := net.Dial("tcp", receiverAddress)
		for forwardConnErr != nil {
			fmt.Println("Error in dialing, dialing again ... ")
			connection1, err1 := net.Dial("tcp", receiverAddress)
			forwardConnection = connection1
			forwardConnErr = err1
		}

		// sending the request further
		newSendRequest := newrequest
		newconn := json.NewEncoder(forwardConnection)
		newconn.Encode(&newSendRequest)

	}
}

// checking existence of file to send it to some peer
func CheckFileExistence(request FileRequest, directoryFiles *ClientFiles) bool {
	// if the file exists
	for _, file := range directoryFiles.FilesInDir {

		if strings.Compare(file, request.RequestedFile) == 0 {
			return true
		}
	}
	return false
}

func handleConnection(connection net.Conn, activeClient *ClientListen, myname string,
	myfiles map[string]MyReceivedFiles, mymessages *MyReceivedMessages, directoryFiles *ClientFiles) {

	var newrequest BaseRequest
	newconn := json.NewDecoder(connection)
	newconn.Decode(&newrequest)

	// If peer is asking for a file, send the existence status of the file in my directory
	if newrequest.RequestType == "ask_for_file" {

		exists := CheckFileExistence(newrequest.FileRequest, directoryFiles)
		if exists {
			_ = RequestMessage(activeClient, myname, newrequest.FileRequest.MyName, "have the file you requested")
		} else {
			_ = RequestMessage(activeClient, myname, newrequest.FileRequest.MyName, "doesn't have the file you requested")
		}

		// if peer is asking me to send some file
	} else if newrequest.RequestType == "receive_from_peer" {

		handleNewFileSendRequest(newrequest.FileRequest, myname, activeClient)

		// if received some file part
	} else if newrequest.RequestType == "received_some_file" {

		handleReceivedRequest(connection, activeClient, myname, myfiles, newrequest)

		// If receievd some message
	} else if newrequest.RequestType == "receive_message" {
		mutexMessages.Lock()
		mymessages.MyMessages = append(mymessages.MyMessages, newrequest.MessageRequest)
		mutexMessages.Unlock()
	}

}

// ListenOnSelfPort listens for clients on network
func ListenOnSelfPort(ln net.Listener, myname string, activeClient *ClientListen, myfiles map[string]MyReceivedFiles,
	mymessages *MyReceivedMessages, directoryFiles *ClientFiles) {
	for {
		connection, err := ln.Accept()

		if err != nil {
			panic(err)
		}
		// Hanling the received connection
		go handleConnection(connection, activeClient, myname, myfiles, mymessages, directoryFiles)
	}
}

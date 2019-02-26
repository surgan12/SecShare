package clientproperties

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//sendMessageRequestToPeer encodes the baserequest
func sendMessageToPeer(connection net.Conn, messageRequest MessageRequest) {
	baseRequest := BaseRequest{RequestType: "receive_message", MessageRequest: messageRequest}
	encoder1 := json.NewEncoder(connection)
	encoder1.Encode(&baseRequest)
	connection.Close()
}

// RequestMessage takes message from client and dials to receiver
func RequestMessage(activeClient *ClientListen, name string, messageReceiverName string,
	message string) string {

// RequestChatting takes message from client and dials to receiver
func RequestChatting(activeClient *ClientListen, name string, messageSenderName string, 
	message string) {
	
	// fmt.Println(message)
	messageRequest := MessageRequest{
		SenderQuery: "message_request", SenderName: name,
		SenderAddress: activeClient.PeerIP[name], Message: message}

	connection, err := net.Dial("tcp", ":"+activeClient.PeerListenPort[messageReceiverName])
	for err != nil {
		fmt.Println("Please enter a valid person name - ")
		connection1, err1 := net.Dial("tcp", ":"+activeClient.PeerListenPort[messageReceiverName])
		connection = connection1
		err = err1
	}
	sendMessageToPeer(connection, messageRequest)
	message_status := "sent"
	return message_status
}

func MessageReceiverCredentials() (string, string) {

	var messageSenderName string
	var message string
	// fmt.Print(name)
	in := bufio.NewReader(os.Stdin)
	fmt.Print("Message (Person's name) : ")
	fmt.Scanln(&messageSenderName)
	fmt.Print("Message to send : ")
	// fmt.Scanln(&message)
	message, err := in.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return messageSenderName, message
}

func MessageReceiverCredentials() (string, string) {
	
	var messageSenderName string
	var message string
	// fmt.Print(name)
	in := bufio.NewReader(os.Stdin)
	fmt.Print("Whom do you want to chat to : ")
	fmt.Scanln(&messageSenderName)
	fmt.Print("What message do you want to send : ")
	// fmt.Scanln(&message)
	message, err := in.ReadString('\n')
	
	if err != nil {
		panic(err)
	}

	return messageSenderName, message
}
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

	messageRequest := MessageRequest{
		SenderQuery: "message_request", SenderName: name,
		SenderAddress: activeClient.PeerIP[name], Message: message}

	connection, err := net.Dial("tcp", ":"+activeClient.PeerListenPort[messageReceiverName])

	count := 0
	for err != nil {
		fmt.Println("Error in dialing to: ", messageReceiverName, " dialing again...")
		connection1, err1 := net.Dial("tcp", ":"+activeClient.PeerListenPort[messageReceiverName])
		connection = connection1
		err = err1
		count++
		if count > 10 {
			message_status := "not sent"
			return message_status
			break
		}
	}

	sendMessageToPeer(connection, messageRequest)
	message_status := "sent"
	return message_status
}

func MessageReceiverCredentials() (string, string) {

	// getting credentials of the person to send message to
	var messageReceiverName string
	var message string // what message to send?
	in := bufio.NewReader(os.Stdin)
	fmt.Print("Message (Person's name) : ")
	fmt.Scanln(&messageReceiverName)
	fmt.Print("Message to send : ")
	// fmt.Scanln(&message)
	message, err := in.ReadString('\n')

	if err != nil {
		panic(err)
	}

	return messageReceiverName, message
}
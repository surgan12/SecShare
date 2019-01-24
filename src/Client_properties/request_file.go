package Client_properties

import (
	"encoding/json"
	"fmt"
	"net"
	// client "../../clients"
)

type FileRequest struct {
	query          string
	my_address     string
	my_name        string
	requested_file string
}

func Request_some_file(active_client Client_listen, name string) {
	var sender_name string // is the person who will send the file
	fmt.Println("Whom do you want to receive the file from ? : ")
	fmt.Scanln(&sender_name)
	var file_name string
	fmt.Println("What file do you want ? ")
	fmt.Scanln(&file_name) // file we want to receive

	file_request := FileRequest{query: "receive_file", my_address: active_client.Peer_IP[name],
		my_name: name, requested_file: "any song"}

	connection, err := net.Dial("tcp", active_client.Peer_IP[sender_name])
	
	for err != nil {
		fmt.Println("Please enter a valid person name - ")
	}

	encoder := json.NewEncoder(connection)
	encoder.Encode(file_request)
}

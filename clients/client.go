package main

import (
	// cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	en "../src/encryptionproperties"
	cp "../src/clientproperties"
	"fmt"
	"net"
	"os"
	"encoding/json"
)

//initializing a wait group
// var wg sync.WaitGroup

func gettingPeersFromServer(c net.Conn, peers *[]string, msg *cp.ClientListen) {
	for {
		d := json.NewDecoder(c)
		d.Decode(msg)
		// fmt.Println("Current active peers are -> ")
		// fmt.Println(msg)
	}
}

func main() {

	var activeClient cp.ClientListen

	var myPeers []cp.MyPeers

	peers := []string{}

	// var myname string
	var name string
	var query string
	var listenPort string
	var flag = false

	_, PublicKey := en.GenerateKeyPair()

	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	ServerKey := en.PerformHandshake(conn, PublicKey)
	// fmt.Print(ServerKey.N)
	fmt.Println("The follwing queries are supported - quit, receive_file - <sender_name>")
	for {

		go gettingPeersFromServer(conn, &peers, &activeClient)
		if flag == false {
			fmt.Print("Enter your credentials : ")
			fmt.Scanln(&name)
			flag = true
			query = "login"

			nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
			queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)

			fmt.Println("Which port do you want to listen upon ? : ")
			fmt.Scanln(&listenPort)
			ln, err := net.Listen("tcp", ":"+listenPort)
			fmt.Println("error on listening is ", err)
			for err != nil {
				fmt.Println("Cant listen on this port, choose another : ")
				fmt.Scanln(&listenPort)
				ln1, err1 := net.Listen("tcp", listenPort)
				ln = ln1
				err = err1
			} 

			mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
			cp.sendingToServer(nameByte, queryByte, conn, query, mylistenport)
			go cp.ListenOnSelfPort(ln)
			continue

		} else {
	
			fmt.Print("What do you want to do? : ")
			fmt.Scanln(&query)
			// flag = false
			if query == "quit" {
				nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
				queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)
				mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
				cp.sendingToServer(nameByte, queryByte, conn, query, mylistenport)
				os.Exit(2)
			} else if query == "receive_file" {
				cp.RequestSomeFile(activeClient, name, myPeers)
			}
		
		}
		
	}
}

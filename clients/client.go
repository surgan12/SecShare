package main

import (
	// cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	// en "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/encryptionproperties"
	en "../src/encryptionproperties"
	cp "../src/clientproperties"
	"encoding/json"
	"fmt"
	"net"
	"os"
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

	// var myPeers []cp.MyPeers
	myfiles := make(map[string]cp.MyReceivedFiles)

	peers := []string{}

	// var myname string
	var name string
	var query string
	var listenPort string
	var flag = false

	_, PublicKey := en.GenerateKeyPair()

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	
	dial_count := 0
	for err != nil {
		fmt.Println("error in connecting to server, dialing again")
		conn1, err1 := net.Dial("tcp", "127.0.0.1:8081")
		conn = conn1
		err = err1
		dial_count++
		if (dial_count > 10){
			fmt.Println("Apparantly server's port is not open...")
			os.Exit(1)
		}
	}

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
			cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
			go cp.ListenOnSelfPort(ln, name, activeClient, myfiles)
			continue

		} else {

			fmt.Print("What do you want to do? : ")
			fmt.Scanln(&query)
			// flag = false
			if query == "quit" {

				nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
				queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)
				mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
				cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
				os.Exit(2)
			} else if query == "receive_file" {
				cp.RequestSomeFile(activeClient, name)
			}

		}

	}
}

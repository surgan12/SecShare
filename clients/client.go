package main

import (
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	en "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/encryptionproperties"
	// cp "../src/clientproperties"
	// en "../src/encryptionproperties"
	"encoding/json"
	"fmt"
	"bufio"
	"net"
	"os"
	"path/filepath"
)

// function to fetch peers which are currently active
func gettingPeersFromServer(c net.Conn, peers *[]string, msg *cp.ClientListen) {
	for {
		// json decoder to decode the struct obtained over the connection
		d := json.NewDecoder(c)
		d.Decode(msg)
	}
}

func main() {

	var activeClient cp.ClientListen  // to store information about peers
	var directoryFiles cp.ClientFiles // to store information about files in the directory
	fileDirectory := "../files"

	// map to keep account of the files received
	myfiles := make(map[string]cp.MyReceivedFiles)
	// struct - containing counter, which gives the index from which we have to read new messages
	// and an array of struct type - MessageRequest to store the attributes of messages
	mymessages := cp.MyReceivedMessages{Counter: 0}

	// to store current peers
	peers := []string{}

	// credentials of the client logging in
	var name string
	var query string
	var listenPort string // which port does he prefer to listen upon
	var flag = false

	// generating keys for connection with server
	_, PublicKey := en.GenerateKeyPair()

	conn, err := net.Dial("tcp", "127.0.0.1:8081")

	// limiting the dial count
	dialCount := 0
	for err != nil {
		// fmt.Println("error in connecting to server, dialing again")
		conn1, err1 := net.Dial("tcp", "127.0.0.1:8081")
		conn = conn1
		err = err1
		dialCount++
		if dialCount > 200 {
			fmt.Println("Apparently server's port is not open...")
			os.Exit(1)
		}
	}

	// performing handshake with server
	ServerKey := en.PerformHandshake(conn, PublicKey)

	// queries currently supported
	fmt.Println("The follwing queries are supported ->")
	fmt.Println("for quitting - quit")
	fmt.Println("for receiving file - receive_file")
	fmt.Println("for sending message - send_message")
	fmt.Println("for displaying recent messages - display_recent_messages")
	fmt.Println("for downloading a file - down")

	for {
		// getting others clients who are currenty active
		go gettingPeersFromServer(conn, &peers, &activeClient)

		// flag == false, signifies the client has to login
		if flag == false {

			fmt.Print("Enter your credentials : ")
			fmt.Scanln(&name)
			flag = true // has logged in
			query = "login"

			// encrypting the details with the PublicKey
			nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
			queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)

			fmt.Println("Which port do you want to listen upon ? : ")
			fmt.Scanln(&listenPort)
			ln, err := net.Listen("tcp", ":"+listenPort)
			fmt.Println("error on listening is ", err)
			for err != nil {
				fmt.Println("Cant listen on this port, choose another : ")
				fmt.Scanln(&listenPort)
				ln1, err1 := net.Listen("tcp", ":"+listenPort)
				ln = ln1
				err = err1
			}

			// adds files in directory to a clientFiles
			error1 := filepath.Walk(fileDirectory, func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				directoryFiles.FilesInDir = append(directoryFiles.FilesInDir, info.Name())
				return nil
			})

			if error1 != nil {
				panic(error1)
			}
			//for printing files in the directory
			// for _, file := range directoryFiles.FilesInDir {
			//     fmt.Println(file)
			// }

			mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
			// sending credentials to server
			cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
			// activating my own port for listening
			go cp.ListenOnSelfPort(ln, name, &activeClient, myfiles, &mymessages, &directoryFiles)
			continue

		} else {
			// accepting further queries, after login is done
			fmt.Print("What do you want to do? : ")
			fmt.Scanln(&query)

			if query == "quit" {
				// encrypting credentials to notify server, when a client wants to quit
				nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
				queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)
				mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
				// sending the information to server
				cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
				os.Exit(2)

			} else if query == "ask_for_file" {
				// broadcasting request for receiving some file
				_, fileName := cp.FileSenderCredentials(true)
				requestStatus := cp.RequestSomeFile(&activeClient, name, fileName)
				if requestStatus == "completed" {
					fmt.Println("Request broadcasted properly")
				} else {
					fmt.Println("Request not broadcasted properly")
				}
				// display the status of file existence from all clients
				cp.DisplayRecentUnseenMessages(&mymessages)

			} else if query == "receive_file" {
				// Receiving file from specific peer
				// getting credentials of the person from whom to receive the file
				fileSenderName, fileName := cp.FileSenderCredentials(false)
				// sending the request to receive some file
				// status of the request, whether or not the request is sent properly
				requestStatus := cp.GetRequestedFile(&activeClient, name, fileSenderName, fileName)
				if requestStatus == "completed" {
					fmt.Println("Request sent")
				} else {
					fmt.Println("Request not sent properly")
				}
			} else if query == "send_message" {
				// query to send messages to others
				// activating the message send mode
				messaging := true
				fmt.Println("Currently in messaging mode..")

				for messaging == true {
					// getting credentials of the person, to whom I want to send some message
					messageReceiverName, message := cp.MessageReceiverCredentials()
					// status of the message, whether or not sent properly
					message_status := cp.RequestMessage(&activeClient, name, messageReceiverName, message)
					if message_status == "sent" {
						fmt.Println("Message sent")
					} else {
						fmt.Println("Message not sent properly")
					}

					// whether I want to exit the message mode
					var queryMessage string
					fmt.Println("Do you want to send more messages? If Yes type Y, else N")
					fmt.Scanln(&queryMessage)
					if queryMessage == "N" {
						fmt.Println("Exiting messaging mode...")
						break
					}
				}

			} else if query == "display_recent_messages" {
				// to display recent messages, which haven't been seen yet
				fmt.Println("Display recent unseen messages - (type) 1")
				fmt.Println("Display recent Num messages - (type) 2")
				var queryMessage string
				fmt.Scanln(&queryMessage)

				// Display recently unseen messages
				if queryMessage == "1" {
					cp.DisplayRecentUnseenMessages(&mymessages)
				} else {
					// display N recent messages
					var num int
					fmt.Println("Number of recent messages you want to see : ")
					fmt.Scanln(&num)
					cp.DisplayNumRecentMessages(&mymessages, num)
				}

			} else if query == "down" {
				// to download files, support within file concurrency and can donwload muliple files simultaneously
				fmt.Print("URL for downloading: ") // url string
				var url string
				scanner := bufio.NewScanner(os.Stdin)
    			scanner.Scan() // use `for scanner.Scan()` to keep reading
    			url = scanner.Text()
				go cp.Download(url)

			}
		}
	}
}

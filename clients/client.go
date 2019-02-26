package main

import (
	// cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	// en "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/encryptionproperties"
	cp "../src/clientproperties"
	en "../src/encryptionproperties"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func gettingPeersFromServer(c net.Conn, peers *[]string, msg *cp.ClientListen) {
	for {
		d := json.NewDecoder(c)
		d.Decode(msg)
	}
}

func main() {

	var activeClient cp.ClientListen
	var directoryFiles cp.ClientFiles
	fileDirectory := "../files"

	myfiles := make(map[string]cp.MyReceivedFiles)
	mymessages := cp.MyReceivedMessages{Counter: -1} // No messages received yet

	peers := []string{}

	var name string
	var query string
	var listenPort string
	var flag = false

	_, PublicKey := en.GenerateKeyPair()

	conn, err := net.Dial("tcp", "127.0.0.1:8081")

	dialCount := 0
	for err != nil {
		fmt.Println("error in connecting to server, dialing again")
		conn1, err1 := net.Dial("tcp", "127.0.0.1:8081")
		conn = conn1
		err = err1
		dialCount++
		if dialCount > 10 {
			fmt.Println("Apparently server's port is not open...")
			os.Exit(1)
		}
	}

	ServerKey := en.PerformHandshake(conn, PublicKey)

	fmt.Println("The follwing queries are supported - quit, receive_file - <sender_name>, send_message, display_recent_messages, down for donwload")

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
				ln1, err1 := net.Listen("tcp", ":"+listenPort)
				ln = ln1
				err = err1
			}

			//adds files in directory to a clientFiles
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
			cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
			go cp.ListenOnSelfPort(ln, name, &activeClient, myfiles, &mymessages)
			continue

		} else {

			fmt.Print("What do you want to do? : ")
			fmt.Scanln(&query)

			if query == "quit" {

				nameByte := en.EncryptWithPublicKey([]byte(name), ServerKey)
				queryByte := en.EncryptWithPublicKey([]byte(query), ServerKey)
				mylistenport := en.EncryptWithPublicKey([]byte(listenPort), ServerKey)
				cp.SendingToServer(nameByte, queryByte, conn, query, mylistenport)
				os.Exit(2)

			} else if query == "receive_file" {

				fileSenderName, fileName := cp.FileSenderCredentials()
				request_status := cp.RequestSomeFile(&activeClient, name, &directoryFiles, fileSenderName, fileName)
				if request_status == "completed" {
					fmt.Println("Request sent")
				} else {
					fmt.Println("Request not sent properly")
				}

			} else if query == "send_message" {

				messaging := true
				fmt.Println("Currently in messaging mode..")

				for messaging == true {
					messageReceiverName, message := cp.MessageReceiverCredentials()
					message_status := cp.RequestMessage(&activeClient, name, messageReceiverName, message)
					if message_status == "sent" {
						fmt.Println("Message sent")
					} else {
						fmt.Println("Message not sent properly")
					}

					var query_message string
					fmt.Println("Do you want to send more messages? If Yes type Y, else N")
					fmt.Scanln(&query_message)
					if query_message == "N" {
						fmt.Println("Exiting messaging mode...")
						break
					}
				}

			} else if query == "display_recent_messages" {

        cp.DisplayRecentMessages(mymessages)
				messageReceiverName, message := cp.MessageReceiverCredentials()
				cp.RequestChatting(&activeClient, name, messageReceiverName, message)

			} else if query == "down" {
				cp.Download()
			}
		}
	}
}

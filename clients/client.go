package main

import (

	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	// cp "../src/clientproperties"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net"
	// "strings"
	// "sync"
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

func sendingToServer(name []byte, query []byte, conn net.Conn, 
					 queryType string, listenPort []byte) {
	objectToSend := cp.ClientQuery{Name: name, Query: query, ClientListenPort: listenPort}
	encoder := json.NewEncoder(conn)
	encoder.Encode(objectToSend)
	if queryType == "quit" {
		conn.Close()
	}
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return privkey, &privkey.PublicKey
}
// PerformHandshake performs handshake with encryption done
func PerformHandshake(conn net.Conn, pub *rsa.PublicKey) *rsa.PublicKey {
	encoder := json.NewEncoder(conn)
	encoder.Encode(pub)
	serverkeys := &rsa.PublicKey{}
	decoder := json.NewDecoder(conn)
	decoder.Decode(&serverkeys)
	return serverkeys
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, _ := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)

	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, _ := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	return plaintext
}

func main() {

	var activeClient cp.ClientListen

	peers := []string{}

	// var myname string
	var name string
	var query string
	var listenPort string
	var flag = false

	_, PublicKey := GenerateKeyPair()

	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	ServerKey := PerformHandshake(conn, PublicKey)
	// fmt.Print(ServerKey.N)
	fmt.Println("The follwing queries are supported - quit, receive_file - <sender_name>")
	for {

		go gettingPeersFromServer(conn, &peers, &activeClient)

		if flag == false {
			fmt.Print("Enter your credentials : ")
			fmt.Scanln(&name)
			flag = true
			query = "login"

			nameByte := EncryptWithPublicKey([]byte(name), ServerKey)
			queryByte := EncryptWithPublicKey([]byte(query), ServerKey)

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

			mylistenport := EncryptWithPublicKey([]byte(listenPort), ServerKey)
			sendingToServer(nameByte, queryByte, conn, query, mylistenport)
			go cp.ListenOnSelfPort(ln)
			continue

		} else {
	
			fmt.Print("What do you want to do? : ")
			fmt.Scanln(&query)
			// flag = false
			if query == "quit" {
				nameByte := EncryptWithPublicKey([]byte(name), ServerKey)
				queryByte := EncryptWithPublicKey([]byte(query), ServerKey)
				mylistenport := EncryptWithPublicKey([]byte(listenPort), ServerKey)
				sendingToServer(nameByte, queryByte, conn, query, mylistenport)
				os.Exit(2)
			} else if query == "receive_file" {
				cp.RequestSomeFile(activeClient, name)
			}
		
		}
		
	}
}

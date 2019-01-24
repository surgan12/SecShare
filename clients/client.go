package main

import (
	"net"
	"fmt"
	"encoding/json"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/rand"
	"strings"
	cp "../src/Client_properties"

)

// type  client_listen struct {
// 	List    []string
// 	Peer_IP map[string]string
// }

// type Client_Query struct {
// 	Name  []byte
// 	Query []byte
// }

func getting_peers_from_server(c net.Conn, peers *[]string, msg *cp.Client_listen) {
	for {
		d := json.NewDecoder(c)
		d.Decode(msg)
		// fmt.Print(msg)
	}
}

func sending_to_server(name []byte, query []byte, conn net.Conn) {
	object_to_send :=  cp.Client_Query{Name: name, Query: query}
	encoder := json.NewEncoder(conn)
	encoder.Encode(object_to_send)
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return privkey, &privkey.PublicKey
}
func PerformHandshake(conn net.Conn, pub *rsa.PublicKey) *rsa.PublicKey {
	encoder := json.NewEncoder(conn)
	encoder.Encode(pub)
	server_keys := &rsa.PublicKey{}
	decoder := json.NewDecoder(conn)
	decoder.Decode(&server_keys)
	return server_keys
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

	var active_client  cp.Client_listen

	peers := []string{}

	// var myname string
	var name string
	var query string
	var flag bool = false

	_, PublicKey := GenerateKeyPair()

	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	ServerKey := PerformHandshake(conn, PublicKey)
	// fmt.Print(ServerKey.N)
	fmt.Println("The follwing queries are supported - quit, receive_file - <sender_name>")
	for {

		if flag == false {
			fmt.Print("Enter your credentials : ")
			fmt.Scanln(&name)
			flag = true
			query = "login"
			name_byte := EncryptWithPublicKey([]byte(name), ServerKey)
			query_byte := EncryptWithPublicKey([]byte(query), ServerKey)
			go sending_to_server(name_byte, query_byte, conn)
			ln2, _ := net.Listen("tcp", strings.TrimLeft(active_client.Peer_IP[name], ":"))
			go cp.ListenOnSelfPort(ln2)
			continue

		} else {

			if flag == true {
				fmt.Print("What do you want to do? : ")
				fmt.Scanln(&query)
				flag = false
				if query == "quit" {
					name_byte := EncryptWithPublicKey([]byte(name), ServerKey)
					query_byte := EncryptWithPublicKey([]byte(query), ServerKey)
					go sending_to_server(name_byte, query_byte, conn)
					continue
				} else if query == "receive_file" {
					cp.Request_some_file(active_client, name)
				}
			}
			go getting_peers_from_server(conn, &peers, &active_client)
		}
	}
}

package main

import "net"
import "fmt"
import "encoding/json"
import "crypto/rsa"
import "crypto/sha512"
import "crypto/rand"
// import "os"


type client_list struct {
	List    []string
	Peer_IP map[string]string
}

type Client_Query struct {
	Name  [] byte
	Query [] byte
}

func getting_peers_from_server(c net.Conn, peers *[]string, msg *client_list) {
	for {
		d := json.NewDecoder(c)
		d.Decode(msg)
		// fmt.Print(msg)
	}
}

func sending_to_server(name [] byte, query [] byte, conn net.Conn) {
	object_to_send := Client_Query{Name: name, Query:query}
	encoder := json.NewEncoder(conn)
	encoder.Encode(object_to_send)
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return privkey, &privkey.PublicKey
}
func PerformHandshake(conn net.Conn, pub *rsa.PublicKey) *rsa.PublicKey{
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

	var active_client client_list

	peers := []string{}

	var myname string
	
	_, PublicKey := GenerateKeyPair()

	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	ServerKey := PerformHandshake(conn, PublicKey)
	// fmt.Print(ServerKey.N)
	for {

		fmt.Print("What do you want to do ")
		var name string
		var query string
		fmt.Scanln(&name)

		if name == "quit" {
			query_byte := EncryptWithPublicKey([]byte(myname), ServerKey)
			name_byte := EncryptWithPublicKey([]byte(name), ServerKey)
			go sending_to_server(name_byte, query_byte, conn)
			continue
		}
		fmt.Print("Enter the query ")
		fmt.Scanln(&query)
		query_byte := EncryptWithPublicKey([]byte(query), ServerKey)
		name_byte := EncryptWithPublicKey([]byte(name), ServerKey)
		
		go sending_to_server(name_byte, query_byte, conn)
		if name == "login" {
			myname = query
			go getting_peers_from_server(conn, &peers, &active_client)

		}
	}
}

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	sp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/serverproperties"
	"net"
	"sync"
	// cp "../src/clientproperties"
	// sp "../src/serverproperties"
)

/* Add relevant Print statements where confused , and comment print statements while pushing */

var mutex = &sync.Mutex{} // Lock and unlock (Mutex)

var clients []cp.Client

var cli cp.ClientListen

var jobs []cp.ClientJob

func pingAll(clients []cp.Client, cli cp.ClientListen) {
	for i := 0; i < len(clients); i++ {
		encode := json.NewEncoder(clients[i].ConnectionServer) //sending to each online peer
		encode.Encode(cli)
	}
}

func performJobs() { // storing each job in a queue in the server and executing it one by one

	for {
		if len(jobs) != 0 {
			mutex.Lock()
			// fmt.Println("number of jobs currently are ", len(jobs))
			getJob := jobs[0]
			jobs = jobs[1:]
			mutex.Unlock()
			// fmt.Println("number of jobs currently are ", len(jobs))
			handler(getJob.Conn, getJob.Name, getJob.Query, getJob.ClientListenPort)
		}
	}
}

func handler(c net.Conn, name string, query string, ClientListenPort string) { // handling each connection

	if query == "login" {

		remoteAddress := c.RemoteAddr().String()
		newClient := cp.Client{Address: remoteAddress, Name: name, ConnectionServer: c} //making struct
		cli.PeerIP[name] = remoteAddress
		cli.PeerListenPort[name] = ClientListenPort //creating the map
		clients = append(clients, newClient)        //append
		cli.List = append(cli.List, name)
		go pingAll(clients, cli)

	} else if query == "quit" {

		delete(cli.PeerIP, name)
		var j int
		for i := 0; i < len(cli.List); i++ {
			if cli.List[i] == name {
				j = i
				break
			}
		}
		cli.List = append(cli.List[:j], cli.List[j+1:]...)
		clients = sp.RemoveFromClient(clients, name)
	}

	fmt.Print("Active clients are -> ", cli.List, "\n")
	fmt.Print("Active clients IPs are -> ", cli.PeerIP, "\n")

}

func maintainConnection(conn net.Conn, PeerKeys map[net.Conn]*rsa.PublicKey,
	pub *rsa.PublicKey, pri *rsa.PrivateKey) { //maintaining the connection between client and server

	//performing handshake
	peerKey := &rsa.PublicKey{}
	decoder := json.NewDecoder(conn)
	decoder.Decode(&peerKey)
	encoder := json.NewEncoder(conn)
	encoder.Encode(pub)
	PeerKeys[conn] = peerKey
	// fmt.Print(peerKey.N)
	for {
		clientQuery := cp.ClientQuery{}

		decoder := json.NewDecoder(conn)
		decoder.Decode(&clientQuery)

		Name := string(DecryptWithPrivateKey(clientQuery.Name, pri))
		Query := string(DecryptWithPrivateKey(clientQuery.Query, pri))
		ClientListenPort := string(DecryptWithPrivateKey(clientQuery.ClientListenPort, pri))
		// fmt.Println("name and query are ", Name, Query)
		job := cp.ClientJob{Name: Name, Query: Query, Conn: conn, ClientListenPort: ClientListenPort}
		// fmt.Println("current job is ", job.Query)

		mutex.Lock()
		if job.Query != "" {
			jobs = append(jobs, job)
			fmt.Println("appended job is ", job)
		}
		mutex.Unlock()
		if Query == "quit" {
			break
		}
	}
	conn = nil
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 2048)

	return privkey, &privkey.PublicKey
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

	ln, _ := net.Listen("tcp", ":8081") // making a server
	fmt.Println(": SERVER STARTED ON PORT 8081  : ")

	PrivateKey, PublicKey := GenerateKeyPair()

	PeerKeys := make(map[net.Conn]*rsa.PublicKey)

	cli = cp.ClientListen{List: []string{}, PeerIP: make(map[string]string),
		PeerListenPort: make(map[string]string)}
	go performJobs()

	for {
		conn, _ := ln.Accept()
		go maintainConnection(conn, PeerKeys, PublicKey, PrivateKey) //accept a new connection and maintain it using the function above
	}

}

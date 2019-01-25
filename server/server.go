package main

import "net"
import "fmt"
import "sync"
import "encoding/json"
import "crypto/rsa"
import "crypto/sha512"
import "crypto/rand"

/* Add relevant Print statements where confused , and comment print statements while pushing */

var mutex = &sync.Mutex{} // Lock and unlock (Mutex)

type client struct {

	address           string
	name              string
	connectionServer net.Conn
}

var clients []client

type clientList struct {
	List    []string
	PeerIP map[string]string
}

var cli clientList

type ClientJob struct {
	name  string
	query string
	conn  net.Conn
}

var jobs []ClientJob

type ClientQuery struct {
	Name  []byte
	Query []byte
}

func removeFromClient(clients []client, name string) []client {
	temp_clients := []client{}
	for i := 0; i < len(clients); i++ {
		if clients[i].name != name {
			temp_clients = append(temp_clients, clients[i])
		}
	}
	return temp_clients
}

func pingAll(clients []client, cli clientList) {
	for i := 0; i < len(clients); i++ {
		encode := json.NewEncoder(clients[i].connectionServer) //sending to each online peer
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
			handler(getJob.conn, getJob.name, getJob.query)
		}
	}
}

func handler(c net.Conn, name string, query string) { // handling each connection

	if query == "login" {

		remoteAddress := c.RemoteAddr().String()
		newClient := client{address: remoteAddress, name: name, connectionServer: c} //making struct
		cli.PeerIP[name] = remoteAddress                                             //creating the map
		clients = append(clients, newClient)                                        //append
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
		clients = removeFromClient(clients, name)
	}

	fmt.Print("Active clients are -> ", cli.List, "\n")
	fmt.Print("Active clients IPs are -> ", cli.PeerIP, "\n")

}

func maintainConnection(conn net.Conn, Peer_Keys map[net.Conn]*rsa.PublicKey, 
			pub *rsa.PublicKey, pri *rsa.PrivateKey) { //maintaining the connection between client and server

	//performing handshake
	peer_key := &rsa.PublicKey{}
	decoder := json.NewDecoder(conn)
	decoder.Decode(&peer_key)
	encoder := json.NewEncoder(conn)
	encoder.Encode(pub)
	Peer_Keys[conn] = peer_key
	// fmt.Print(peer_key.N)
	for {
		clientQuery := ClientQuery{}

		decoder := json.NewDecoder(conn)

		decoder.Decode(&clientQuery)
		Name := string(DecryptWithPrivateKey(clientQuery.Name, pri))
		Query := string(DecryptWithPrivateKey(clientQuery.Query, pri))
		// fmt.Println("name and query are ", Name, Query)
		job := ClientJob{name: Name, query: Query, conn: conn}
		mutex.Lock()
		jobs = append(jobs, job)
		fmt.Println("appended job is ", job)
		mutex.Unlock()
		if Query == "quit"{
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
	fmt.Println("server started on port 8081")

	PrivateKey, PublicKey := GenerateKeyPair()

	PeerKeys := make(map[net.Conn]*rsa.PublicKey)

	cli = clientList{List: []string{}, PeerIP: make(map[string]string)}
	go performJobs()

	for {
		conn, _ := ln.Accept()
		go maintainConnection(conn, PeerKeys, PublicKey, PrivateKey) //accept a new connection and maintain it using the function above
	}

}

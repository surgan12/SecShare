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
	connection_server net.Conn
}

var clients []client

type client_list struct {
	List    []string
	Peer_IP map[string]string
}

var cli client_list

type ClientJob struct {
	name  string
	query string
	conn  net.Conn
}

var jobs []ClientJob

type Client_Query struct {
	Name  []byte
	Query []byte
}

func remove_from_client(clients []client, name string) []client {
	temp_clients := []client{}
	for i := 0; i < len(clients); i++ {
		if clients[i].name != name {
			temp_clients = append(temp_clients, clients[i])
		}
	}
	return temp_clients
}

func ping_all(clients []client, cli client_list) {
	for i := 0; i < len(clients); i += 1 {
		encode := json.NewEncoder(clients[i].connection_server) //sending to each online peer
		encode.Encode(cli)
	}
}

func perform_jobs() { // storing each job in a queue in the server and executing it one by one

	for {
		if len(jobs) != 0 {
			mutex.Lock()
			get_job := jobs[0]
			jobs = jobs[1:]
			mutex.Unlock()
			handler(get_job.conn, get_job.name, get_job.query)
		}
	}
}

func handler(c net.Conn, name string, query string) { // handling each connection

	if query == "login" {

		remote_addr := c.RemoteAddr().String()
		new_client := client{address: remote_addr, name: name, connection_server: c} //making struct
		cli.Peer_IP[name] = remote_addr                                             //creating the map
		clients = append(clients, new_client)                                        //append
		cli.List = append(cli.List, name)
		go ping_all(clients, cli)

	} else if query == "quit" {
		
		delete(cli.Peer_IP, name)
		var j int
		for i := 0; i < len(cli.List); i++ {
			if cli.List[i] == name {
				j = i
				break
			}
		}
		cli.List = append(cli.List[:j], cli.List[j+1:]...)
		clients = remove_from_client(clients, name)
	}

	fmt.Print("Active clients are -> ", cli.List, "\n")
	fmt.Print("Active clients IPs are -> ", cli.Peer_IP, "\n")

}

func maintain_connection(conn net.Conn, Peer_Keys map[net.Conn]*rsa.PublicKey, pub *rsa.PublicKey, pri *rsa.PrivateKey) { //maintaining the connection between client and server

	//performing handshake
	peer_key := &rsa.PublicKey{}
	decoder := json.NewDecoder(conn)
	decoder.Decode(&peer_key)
	encoder := json.NewEncoder(conn)
	encoder.Encode(pub)
	Peer_Keys[conn] = peer_key
	// fmt.Print(peer_key.N)
	for {
		client_query := Client_Query{}

		decoder := json.NewDecoder(conn)

		decoder.Decode(&client_query)
		Name := string(DecryptWithPrivateKey(client_query.Name, pri))
		Query := string(DecryptWithPrivateKey(client_query.Query, pri))
		// fmt.Print(Name, Query)
		job := ClientJob{name: Name, query: Query, conn: conn}
		mutex.Lock()
		jobs = append(jobs, job)
		mutex.Unlock()
	}
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

	Peer_Keys := make(map[net.Conn]*rsa.PublicKey)

	cli = client_list{List: []string{}, Peer_IP: make(map[string]string)}
	go perform_jobs()

	for {
		conn, _ := ln.Accept()
		go maintain_connection(conn, Peer_Keys, PublicKey, PrivateKey) //accept a new connection and maintain it using the function above
	}

}

package clientproperties

import (
	"encoding/json"
	// "fmt"
	// fp "../../fileproperties"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	fp "github.com/IITH-SBJoshi/concurrency-decentralized-network/fileproperties"
	"net"
)

// Client properties as stored in the server
type Client struct {
	Address          string
	Name             string
	ConnectionServer net.Conn
}

// ClientListen stores list of clients and map of their IP
type ClientListen struct {
	List           []string
	PeerIP         map[string]string
	PeerListenPort map[string]string
}

// ClientQuery stores name and query of clients
type ClientQuery struct {
	Name             []byte
	Query            []byte
	ClientListenPort []byte
}

// ClientJob stores the names, jobs and connection
type ClientJob struct {
	Name             string
	Query            string
	Conn             net.Conn
	ClientListenPort string
}

//MyPeers list of connections dialed by current client
type MyPeers struct {
	Conn     net.Conn
	PeerName string
}

type MyReceivedFiles struct {
	MyFileName string
	MyFile     []FilePartContents
}

type FilePartContents struct {
	Contents []byte
}

type BaseRequest struct {
	RequestType string
	FileRequest
	fp.FilePartInfo
}

// FileRequest stores the queries and information about requester
type FileRequest struct {
	query         string
	myAddress     string
	myName        string
	requestedFile string
}

// sendingToServer function to send queries to server
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

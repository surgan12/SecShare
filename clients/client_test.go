package main

import ("testing"
		"io"
		"os"
		"bytes"
		)

import cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientProperties"

import sp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/ServerProperties"

func TestRequestSomeFile(t *testing.T) {
	activeClient := cp.ClientList{List : {"Raju", "Bablu"}, PeerIP : {["Raju"]->"001", ["Bablu"]->"007"}}
	name := "Raju"
	old := os.Stdout // keep backup of the real stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    RequestSomeFile(activeClient, name)

    outC := make(chan string)
    // copy the output in a separate goroutine so printing can't block indefinitely
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outC <- buf.String()
    }()

    // back to normal state
    w.Close()
    os.Stdout = old // restoring the real stdout
    out := <-outC

    // reading our temp stdout
    //fmt.Println("previous output:")
    //fmt.Print(out)
	if out[69:len(out)]=="Please enter a valid person name - " {
		t.Fatal("Not Working")
	}
}
	
	
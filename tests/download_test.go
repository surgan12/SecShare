package main
	
import (
	cp "github.com/IITH-SBJoshi/concurrency-decentralized-network/src/clientproperties"
	"testing"
	"net/http"
	"os"
	"strconv"
	"sync"
	"fmt"
	"io/ioutil"
)

func TestAsyncDownloader(t *testing.T)	{
	
	client := &http.Client{}
	name := "test.mp3"
	var start int64
	start = 0
	url := "http://file-examples.com/wp-content/uploads/2017/11/file_example_MP3_1MG.mp3"
	resp, _ := http.Head(url)
	length := resp.Header.Get("content-length")
	lenh, _ := strconv.Atoi(length)
	end := lenh
	dummy := make([]byte, lenh)
	ioutil.WriteFile(fmt.Sprintf(name), dummy, 0644)
	f, _ := os.OpenFile(name, os.O_RDWR, 0666)
	var wg sync.WaitGroup
	wg.Add(1)
	result := cp.DummyAsync(&wg, client, start, end, 0, lenh, url, f)
	
	if result != nil{
		t.Fatal("Download not working correctly")
	}
	
}
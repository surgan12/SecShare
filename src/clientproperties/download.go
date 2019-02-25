package clientproperties

import "fmt"
import "bytes"
import "net/http"
import "io"
import "strings"
import "sync"
import "os"
import "strconv"
import "io/ioutil"

var wg sync.WaitGroup

func AsyncDownloader(client *http.Client, start int64 , end int, i int, size int, url string, f *os.File) {

if end > size {
	end = size
}
fmt.Println(end)
startString := strconv.FormatInt(start, 10)

endString := strconv.Itoa(end)
req, _ := http.NewRequest("GET", url, nil)

req.Header.Set("Range", "bytes="+startString+"-"+endString)
res, _ := client.Do(req)

f.Seek(start, 0)
var buf bytes.Buffer
io.Copy(&buf, res.Body)
var buffer []byte
buffer = buf.Bytes()
f.Write(buffer)
 
wg.Done()
}


func Download(){

var url string

fmt.Print("URL for downloading: ")
fmt.Scanln(&url)

s := strings.Split(url, "/")
flen := len(s)
name := s[flen - 1]

client := &http.Client{}
resp, _ := http.Head(url)
length := resp.Header.Get("content-length")
lenh ,_:= strconv.Atoi(length)
dummy := make([]byte, lenh)
ioutil.WriteFile(fmt.Sprintf(name), dummy, 0644)
f, _ := os.OpenFile(
        name,
        os.O_RDWR,
        0666,)
var start int64
start = 0
end := 0
part_length := int(lenh/4)
end = part_length

wg.Add(4)
fmt.Print("Downloading has begun")
for i := 0; i < 4; i++ {

	go AsyncDownloader(client, start, end, i, lenh, url, f)	
	start = int64(end) + 1
	end = end + part_length
	
}

wg.Wait() // waiting for the goroutines to finish

}
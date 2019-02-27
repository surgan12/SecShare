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
import (
  	"time" 
    "github.com/andlabs/ui"
    "github.com/andlabs/ui/winmanifest"
)

// var mainwin *ui.Window


type WriteCounter struct {
	Total uint64
}
var counter WriteCounter

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	// wc.PrintProgress()
	return n, nil
}

var help_map = map[string]*WriteCounter{}
var len_map = map[string]int{}
var mainwin = map[string]*ui.Window{}
var prog = map[string]*ui.ProgressBar{}
var name string

func AsyncDownloader(wg *sync.WaitGroup, client *http.Client, start int64 , end int, i int, size int, url string, f *os.File) {

if end > size {
	end = size
}
var fname string
fname = name
startString := strconv.FormatInt(start, 10)

endString := strconv.Itoa(end)
req, _ := http.NewRequest("GET", url, nil)


req.Header.Set("Range", "bytes="+startString+"-"+endString)
res, _ := client.Do(req)

f.Seek(start, 0)
var buf bytes.Buffer
io.Copy(&buf, io.TeeReader(res.Body,help_map[fname]))
// io.Copy(&buf, res.Body)
var buffer []byte
buffer = buf.Bytes()
f.Write(buffer)
 
wg.Done()
}

func set(ip *ui.ProgressBar) {
	var fname string
    fname = name
    lenth := len_map[fname]
for int(help_map[fname].Total) + 1 < lenth{
		
        g := int(help_map[fname].Total)*100
		val := int(g/lenth)
		time.Sleep(200 * time.Millisecond)
		// fmt.Print(val)
        if val > 90{
            break
        }
        ip.SetValue(val)
		
}

}
func Download(url string){

s := strings.Split(url, "/")
flen := len(s)
name = s[flen - 1]
var fname string
fname = name
help_map[fname] = &WriteCounter{}
client := &http.Client{}
resp, _ := http.Head(url)
length := resp.Header.Get("content-length")
lenh ,_:= strconv.Atoi(length)
len_map[fname] = lenh
dummy := make([]byte, lenh)
ioutil.WriteFile(fmt.Sprintf(name), dummy, 0644)
f, _ := os.OpenFile(name, os.O_RDWR, 0666, )
var start int64
start = 0
end := 0
part_length := int(lenh/4)
end = part_length
var wg sync.WaitGroup
wg.Add(4)
// for testing use http://file-examples.com/wp-content/uploads/2017/11/file_example_MP3_1MG.mp3
for i := 0; i < 4; i++ {

	go AsyncDownloader(&wg, client, start, end, i, lenh, url, f)	
	start = int64(end) + 1
	end = end + part_length
	
}
go ui.Main(setupUI)
wg.Wait()
mainwin[fname].Destroy() // waiting for the goroutines to finish
// prog[fname].Destroy
ui.Quit()
}

func makeBasicControlsPage() ui.Control {
    vbox := ui.NewVerticalBox()
    vbox.SetPadded(true)
    var fname string
    fname = name
    hbox := ui.NewHorizontalBox()
    hbox.SetPadded(true)
    vbox.Append(hbox, false)
    prog[fname] = ui.NewProgressBar()
    hbox.Append(prog[fname], false)
    go set(prog[fname])
    // ip.SetValue(10)
    
    return vbox
}

func setupUI() {
    var fname string
    fname = name
    mainwin[fname] = ui.NewWindow("Downloading" + fname, 300, 150, true)
    mainwin[fname].OnClosing(func(*ui.Window) bool {
        ui.Quit()
        return true
    })
    ui.OnShouldQuit(func() bool {
        mainwin[fname].Destroy()
        return true
    })

    tab := ui.NewTab()
    mainwin[fname].SetChild(tab)
    mainwin[fname].SetMargined(true)

    tab.Append("Downloading", makeBasicControlsPage())
    tab.SetMargined(0, true)
    // go set(tab)
    mainwin[fname].Show()
}

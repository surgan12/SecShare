package main

import (
	"time"
	"fmt"
	"os"
	"io/ioutil"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func concatFiles(i int, allFiles []byte) {
	defer wg.Done()

	filename := "part_" + strconv.Itoa(i)

	fileContents, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for j := i*len(fileContents); j < i*len(fileContents) + len(fileContents) ; j++ {
		allFiles[j] = fileContents[j-i*len(fileContents)]
	}

}

func main() {

	file, err := os.Open("./image.jpg")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	allFiles := make([]byte, size)
	
	startTme := time.Now()
	for i := int(0); i < 10; i++ {
		wg.Add(1)
		go concatFiles(i, allFiles)
	}
	wg.Wait()
	endTime := time.Now()
	fmt.Println("Time to concatenate files is ", endTime.Sub(startTme))

	currentfilename := "concatenated.jpg"
	ioutil.WriteFile(currentfilename, allFiles, os.ModeAppend)

}


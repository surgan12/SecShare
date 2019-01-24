package main

import (
	"time"
	"fmt"
	"os"
	"math"
	"io/ioutil"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func splitFile(partSize uint64, filesize int64, i uint64, fileContents []byte) {
	fmt.Println("Writing part ", i, " from file")
	currentSize := int(math.Min(float64(partSize), float64(int64(filesize) - int64(i*partSize))))
	currentBuffer := make([]byte, currentSize)

	for j := int(i*partSize); j < int(i*partSize) + int(currentSize); j++ {
		currentBuffer[j-int(i*partSize)] = fileContents[j]
	}

	currentfilename := "part_" + strconv.FormatUint(i, 10)

	defer wg.Done()

	_, err := os.Create(currentfilename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Writing form byte array to current file part
	ioutil.WriteFile(currentfilename, currentBuffer, os.ModeAppend)
	
}

func main() {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Enter file path - ")
	// filename, err := reader.ReadString('\n')

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Printf("File to be opened is -> %s, error is %s", filename, err)

	// Opening file
	file, err := os.Open("./image.jpg")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	//Fetching info about file
	fileInfo, _ := file.Stat()

	var filesize int64 = fileInfo.Size()
	fmt.Println("Size of file is -> ", filesize)

	var partSize = uint64(math.Ceil(float64(filesize) / float64(10)))

	fileContents, err := ioutil.ReadFile("./image.jpg")

	start_tme := time.Now()
	
	for i := uint64(0); i < 10; i++ {
		wg.Add(1)
		go splitFile(partSize, filesize, i, fileContents)
	}

	wg.Wait()	// waiting for all routines to finish
	end_time := time.Now()
	fmt.Println("Time taken to write the file ",end_time.Sub(start_tme))

}
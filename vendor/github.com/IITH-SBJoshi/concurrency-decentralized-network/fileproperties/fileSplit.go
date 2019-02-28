package fileproperties

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

//FilePartInfo for information regarding files
type FilePartInfo struct {
	FileName         string
	TotalParts       int
	PartName         string
	PartNumber       int
	FilePartContents []byte
}

// GetFileParts - goroutine to get all the parts of file to be sent over the network
func GetFileParts(completefilename string, partSize uint64, filesize int64, i uint64, fileContents []byte,
	allFileParts []FilePartInfo, wgSplit *sync.WaitGroup) {
	// size of the current part to be sent
	currentSize := int(math.Min(float64(partSize), float64((filesize)-int64(i*partSize))))
	// struct containing additional information alongside the byte contents
	currentpart := FilePartInfo{FileName: completefilename, TotalParts: 1, PartName: "part_" + strconv.FormatUint(i, 10),
		PartNumber: int(i), FilePartContents: make([]byte, currentSize)}

	for j := int(i * partSize); j < int(i*partSize)+int(currentSize); j++ {
		currentpart.FilePartContents[j-int(i*partSize)] = fileContents[j]
	}

	allFileParts[i] = currentpart
	// deferring the routine count by calling Done for the current goroutine
	defer wgSplit.Done()
}

//GetSplitFile fuction to return the splitted parts
func GetSplitFile(filename string, numberOfActiveClient int) []FilePartInfo {
	fileDirectory := "../files"
	file, err := os.Open(fileDirectory + "/image.jpg")
	// file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Fetching info about file
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var filesize = fileInfo.Size()
	// fmt.Println("Size of file is -> ", filesize)

	// Currently sending to one peer only
	var partSize = uint64(math.Ceil(float64(filesize) / float64(numberOfActiveClient-1)))

	// contents of the whole file as a byte array
	fileContents, err := ioutil.ReadFile(fileDirectory + "/image.jpg")

	startTime := time.Now()

	// slice of size numberOfActiveClient-1, to store all parts' structs
	allFileParts := make([]FilePartInfo, numberOfActiveClient-1)

	// waitgroup to wait for all goroutines to finish
	var wgSplit sync.WaitGroup

	for i := uint64(0); i < uint64(numberOfActiveClient-1); i++ {
		wgSplit.Add(1) // incrementing the count of waitgroup, for another goroutine is run
		go GetFileParts(filename, partSize, filesize, i, fileContents, allFileParts, &wgSplit)
	}

	wgSplit.Wait() // waiting for all routines to finish
	endTime := time.Now()
	fmt.Println("Time taken to split the file ", endTime.Sub(startTime))

	// closing the file
	defer file.Close()

	return allFileParts // returning the slice
}
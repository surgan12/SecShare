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

var wgSplit sync.WaitGroup

var mutex = &sync.Mutex{} // Lock and unlock (Mutex)
//FilePartInfo for information regarding files
type FilePartInfo struct {
	FileName         string
	TotalParts       int
	PartName         string
	PartNumber       int
	FilePartContents []byte
}

//getFileParts ..
func getFileParts(completefilename string, partSize uint64, filesize int64, i uint64, fileContents []byte,
	allFileParts []FilePartInfo) {
	// fmt.Println("Writing part ", i, " from file")
	currentSize := int(math.Min(float64(partSize), float64((filesize)-int64(i*partSize))))

	currentpart := FilePartInfo{FileName: completefilename, TotalParts: 1, PartName: "part_" + strconv.FormatUint(i, 10),
		PartNumber: int(i), FilePartContents: make([]byte, currentSize)}

	for j := int(i * partSize); j < int(i*partSize)+int(currentSize); j++ {
		currentpart.FilePartContents[j-int(i*partSize)] = fileContents[j]
	}

	mutex.Lock()
	allFileParts[i] = currentpart
	// fmt.Println("done")
	mutex.Unlock()

	defer wgSplit.Done()
}

//GetSplitFile fuction to split files 
func GetSplitFile(filename string, numberOfActiveClient int) []FilePartInfo {
	fileDirectory := "../files"
	file, err := os.Open(fileDirectory + "/image.jpg")
	// file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	//Fetching info about file
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var filesize = fileInfo.Size()
	// fmt.Println("Size of file is -> ", filesize)

	// Currently sending to one peer only
	var partSize = uint64(math.Ceil(float64(filesize) / float64(1)))

	fileContents, err := ioutil.ReadFile(fileDirectory + "/image.jpg")

	startTime := time.Now()

	allFileParts := make([]FilePartInfo, numberOfActiveClient - 1)

	for i := uint64(0); i < 1; i++ {
		wgSplit.Add(numberOfActiveClient - 1)
		go getFileParts(filename, partSize, filesize, i, fileContents, allFileParts)
	}

	// wgSplit.Wait() // waiting for all routines to finish
	endTime := time.Now()
	fmt.Println("Time taken to split the file ", endTime.Sub(startTime))

	return allFileParts

}

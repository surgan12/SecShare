package clientproperties

import (
	// "time"
	"fmt"
	"io/ioutil"
	"os"
	// "strconv"
	"sync"
)

// creating waitgroup to wait for all goroutines to finish
var wgConcat sync.WaitGroup

// func to write the filepart details to all Files byte slice
func concatFiles(i int, allFiles []byte, filePartContent FilePartContents) {
	defer wgConcat.Done()

	filePartContents := filePartContent.Contents

	for j := i * len(filePartContents); j < i*len(filePartContents)+len(filePartContents); j++ {
		allFiles[j] = filePartContents[j-i*len(filePartContents)]
	}

}

func concatenateFileParts(file MyReceivedFiles) {

	// getting total size of all parts
	var byteSizeLength int
	fileName := file.MyFileName
	fileParts := file.MyFile

	for i := 0; i < len(fileParts); i++ {
		byteSizeLength += len(fileParts[i].Contents)
	}

	// creating new byte slice for creating new file
	allFiles := make([]byte, byteSizeLength)

	// writing the received parts to allFiles slice
	for i := int(0); i < 1; i++ {
		wgConcat.Add(1)
		go concatFiles(i, allFiles, fileParts[i])
	}

	wgConcat.Wait()

	// writing the received file
	currentfilename := "Received_" + fileName
	ioutil.WriteFile(currentfilename, allFiles, os.ModeAppend)

	// Test File existence.
	_, err := os.Stat(currentfilename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file doesn't exist")
		}
	}

	// Change permissions in Linux.
	err = os.Chmod(currentfilename, 0777)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nReceived file : ", currentfilename, "\n")
}

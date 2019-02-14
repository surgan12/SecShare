package clientproperties

import (
	"time"
	"fmt"
	"os"
	"io/ioutil"
	// "strconv"
	"sync"
)

var wgConcat sync.WaitGroup

func concatFiles(i int, allFiles []byte, filePartContent FilePartContents) {
	defer wgConcat.Done()

	filePartContents := filePartContent.Contents

	for j := i*len(filePartContents); j < i*len(filePartContents) + len(filePartContents) ; j++ {
		allFiles[j] = filePartContents[j-i*len(filePartContents)]
	}

}

func concatenateFileParts (file MyReceivedFiles) {

	fileName := file.MyFileName
	fileParts := file.MyFile

	//TODO : calculate file size from file parts
	allFiles := make([]byte, 1)
	
	startTme := time.Now()
	
	for i := int(0); i < 1; i++ {
		wgConcat.Add(1)
		go concatFiles(i, allFiles, fileParts[i])
	}

	wgConcat.Wait()
	endTime := time.Now()
	fmt.Println("Time to concatenate all parts is ", endTime.Sub(startTme))

	currentfilename := fileName
	ioutil.WriteFile(currentfilename, allFiles, os.ModeAppend)
	fmt.Println("Received file : ", currentfilename)
}
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

func concatFiles(i int, allFiles *[]byte, filePartContents *[]byte) {
	defer wgConcat.Done()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for j := i*len(filePartContents); j < i*len(filePartContents) + len(filePartContents) ; j++ {
		allFiles[j] = filePartContents[j-i*len(filePartContents)]
	}

}

func concatenateFileParts (filename MyReceivedFiles.MyFileName, fileParts MyReceivedFiles.MyFile) {

	allFiles := make([]byte, size)
	
	startTme := time.Now()
	
	for i := int(0); i < 1; i++ {
		wgConcat.Add(1)
		go concatFiles(i, &allFiles, &fileParts[i])
	}

	wgConcat.Wait()
	endTime := time.Now()
	fmt.Println("Time to concatenate all parts is ", endTime.Sub(startTme))

	currentfilename := filename
	ioutil.WriteFile(currentfilename, allFiles, os.ModeAppend)
	fmt.Println("Received file : ", currentfilename)
}
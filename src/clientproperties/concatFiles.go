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
	var byteSizeLength int
	fileName := file.MyFileName
	fileParts := file.MyFile
	for i := 0; i < len(fileParts); i++ {
		byteSizeLength += len(fileParts[i].Contents)
	}
	//TODO : calculate file size from file parts
	fmt.Println(byteSizeLength)
	allFiles := make([]byte, byteSizeLength)
	
	startTme := time.Now()
	
	for i := int(0); i < 1; i++ {
		wgConcat.Add(1)
		go concatFiles(i, allFiles, fileParts[i])
	}

	wgConcat.Wait()
	endTime := time.Now()
	fmt.Println("Time to concatenate all parts is ", endTime.Sub(startTme))

	currentfilename := fileName + ".jpg"
	fmt.Println("check")
	ioutil.WriteFile(currentfilename, allFiles, os.ModeAppend)
	// currentfilename = os.Chmod(currentfilename, 0777)

	// Test File existence.
	_, err := os.Stat(currentfilename)
	if err != nil {
		if os.IsNotExist(err) {
			// log.Fatal("File does not exist.")
			fmt.Println("file doesn't exist")
		}
	}
	fmt.Println("File exist.")
 
	// Change permissions Linux.
	err = os.Chmod(currentfilename, 0777)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Received file : ", currentfilename)
}

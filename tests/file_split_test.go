package main

import (
	"testing"
	// fp "github.com/IITH-SBJoshi/concurrency-decentralized-network/fileproperties"
	fp "../fileproperties"
)

func TestFileSplit(t *testing.T) {

	filename := "SomeImage"
	allFileParts := fp.GetSplitFile(filename, 2)

	if len(allFileParts[0].FilePartContents) == 0 {
		t.Fatal("File not properly written to allFileParts slice")
	}
}

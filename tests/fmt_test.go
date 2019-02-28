package main

import (
	"testing"
	"ox/exec"
)

//TestGoFmt test for format errors
func TestGoFmt(t *testing.T) {
	cmd := exec.Command("gofmt", "-l", ".")

	if out, err := cmd.Output(); err != nil {
		if len(out) > 0 {
			t.Fatalf("Exit error: %v", err)
		}
	} else {
		if len(out) > 0 {
			t.Fatal("You need to run gofmt")
		}
	}
}

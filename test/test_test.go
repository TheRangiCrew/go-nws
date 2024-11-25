package main

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	_, err := os.Open("../assets/awips/test/1.txt")
	if err != nil {
		t.Fatal(err)
	}
	// stat, err := file.Stat()
	// data := make([]byte, file.Stat())
	// _, err := file.R
}

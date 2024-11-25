package awips

import (
	"fmt"
	"os"

	"github.com/TheRangiCrew/go-nws/pkg/awips"
)

func ParseAwips(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// TODO: Handle directories
	if stat.IsDir() {
		fmt.Println("provided path is a directory")
		return
	}

	bytes := make([]byte, stat.Size())
	file.Read(bytes)

	text := string(bytes)

	product, err := awips.New(text)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(product.Segments[0].HasVTEC())
}

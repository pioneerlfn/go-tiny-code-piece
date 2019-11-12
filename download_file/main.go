package main

import (
	"fmt"
	"go-tiny-code-piece/download_file/download"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: download url filename")
		os.Exit(1)
	}
	url := os.Args[1]
	filename := os.Args[2]

	err := download.DownloadFile(url, filename)
	if err != nil {
		panic(err)
	}
}

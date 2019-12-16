/*
@Time : 2019-12-16 13:55
@Author : lfn
@File : count_concurrent
*/

package word_count

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

func CountConcurrent() {
	if len(os.Args) < 2 {
		panic("no file path specified")
	}
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	b, w, l := countConcurrent(file)
	fmt.Println(b, w, l, file.Name())
}


func countConcurrent(file *os.File) (int, int, int) {
	num := runtime.NumCPU()

	chunks := make(chan Chunk)
	counts := make(chan Count)


	for i := 0; i < num; i++ {
		go GetCounts(chunks, counts)
	}
	readFile(file, chunks)

	var res Count
	for i :=0; i < num; i++ {
		count := <- counts
		res.LineCount += count.LineCount
		res.WordCount += count.WordCount
	}
	close(counts)

	byteCount := getBytes(file)
	return byteCount, res.WordCount, res.LineCount

}

func readFile(file *os.File, chunks chan<- Chunk) {
	const bufferSize  = 16 * 1024
	buffer := make([]byte, bufferSize)

	prevCharIsSpace := true
	for {
		b, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		chunks <- Chunk{prevCharIsSpace, buffer[:b]}
		prevCharIsSpace = IsSpace(buffer[b-1])
	}
	close(chunks)

}

func GetCounts(chunks <-chan Chunk, counts chan<- Count) {
	var total Count
	for {
		chunk, ok := <- chunks
		if !ok {
			break
		}
		cs := getCount(chunk)
		total.WordCount += cs.WordCount
		total.LineCount += cs.LineCount
	}

	counts <- total
}
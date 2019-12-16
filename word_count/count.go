/*
@Time : 2019-12-16 13:53
@Author : lfn
@File : word_count
*/

package word_count

import (
	"bufio"
	"io"
	"os"
)

func count(path string) (byteCount, wordCount, lineCount int) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	const bufferSize = 16 * 1024
	reader := bufio.NewReaderSize(file, bufferSize)

	preByteIsSpace := true
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		byteCount++

		switch b {
		case '\n':
			lineCount++
			preByteIsSpace = true
		case ' ', '\t', '\r', '\v', '\f':
			preByteIsSpace = true
		default:
			if preByteIsSpace {
				wordCount++
				preByteIsSpace = false
			}
		}
	}
	return
}

type Chunk struct {
	PrevCharIsSpace bool
	Buffer []byte
}

type Count struct {
	LineCount int
	WordCount int
}

func getCount(chunk Chunk) Count {
	count := Count{}
	prevCharIsSpace := chunk.PrevCharIsSpace
	for _, b := range chunk.Buffer {
		switch b {
		case '\n':
			count.LineCount++
			prevCharIsSpace = true
		case ' ', '\t', '\f', '\v', '\r':
			prevCharIsSpace = true
		default:
			if prevCharIsSpace {
				prevCharIsSpace = false
				count.WordCount++
			}
		}
	}
	return count
}

// count2自己开辟缓存，不利用bufio包.
func count2(file *os.File) (int, int, int) {
	totalCount := Count{}
	lastCharIsSpace := true
	const bufferSize = 16 * 1024
	buffer := make([]byte, bufferSize)

	for {
		bytes, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		count := getCount(Chunk{lastCharIsSpace,buffer[:bytes]})
		lastCharIsSpace = IsSpace(buffer[bytes-1])

		totalCount.LineCount += count.LineCount
		totalCount.WordCount += count.WordCount
	}


	byteCount := getBytes(file)
	return byteCount, totalCount.WordCount, totalCount.LineCount
}

func IsSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\v' || b == '\f' || b == '\r' || b == '\n'
}

func getBytes(file *os.File) int {
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return int(stat.Size())
}
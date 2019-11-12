package download

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(buf []byte) (int, error) {
	wc.Total += uint64(len(buf))
	wc.WriteProgress()
	return len(buf), nil
}

func (wc *WriteCounter) WriteProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 50))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

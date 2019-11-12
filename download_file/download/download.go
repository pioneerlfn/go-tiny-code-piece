package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url, filename string) error {
	out, err := os.Create(filename + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	fmt.Println()
	err = os.Rename(filename+".tmp", filename)
	if err != nil {
		return err
	}

	return nil
}

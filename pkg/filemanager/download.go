package filemanager

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, path string) (int64, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GET %s: %s", url, response.Status)
	}

	destination, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	bytes_written, err := io.Copy(destination, response.Body)
	if err != nil {
		return 0, err
	}

	return bytes_written, nil
}

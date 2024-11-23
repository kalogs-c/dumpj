package filemanager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

func GetFileSize(url string) (int64, error) {
	response, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("GET %s: %s", url, response.Status)
	}

	size := response.ContentLength
	if size == -1 {
		return size, fmt.Errorf("size is unknown")
	}

	return size, nil
}

func FileAlreadyExists(path string, originalSize int64) bool {
	stats, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	if originalSize != -1 && stats.Size() != originalSize {
		return false
	}

	if time.Since(stats.ModTime()).Hours() > 24*30 {
		return false
	}

	return true
}

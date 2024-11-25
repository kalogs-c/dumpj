package main

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/kalogs-c/dumpj/pkg/crawler"
	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

// TODO
// - [ ] Handle errors (retry, log, etc)
// - [ ] Parse CSV to sqlite3
// - [ ] Avoid duplicates and cache previous dump
// - [ ] Yaml config
// - [ ] CLI
// - [ ] Docker
// - [ ] Logging

func main() {
	path := "https://dados-abertos-rf-cnpj.casadosdados.com.br/arquivos/"

	content, err := crawler.Fetch(path)
	if err != nil {
		panic(err)
	}

	reg := regexp.MustCompile("^20[0-9]{2}-[0-1][0-9]-[0-3][0-9]")
	links := crawler.ScrapeLinks(content, reg)
	latest := links[len(links)-1]

	files_path := fmt.Sprintf("%s%s", path, latest)
	content, err = crawler.Fetch(files_path)
	if err != nil {
		panic(err)
	}

	reg = regexp.MustCompile(".zip")
	links = crawler.ScrapeLinks(content, reg)

	downloadWg := sync.WaitGroup{}
	zipch := make(chan string, len(links)/4)

	for _, link := range links {
		downloadWg.Add(1)
		go func(link string) {
			defer downloadWg.Done()

			httppath := fmt.Sprintf("%s%s", files_path, link)
			fpath := fmt.Sprintf("./_files/zips/%s", link)

			filesize, err := filemanager.GetFileSize(httppath)
			if err != nil {
				panic(err)
			}

			if filesize > int64(math.Pow(10, 8)) {
				fmt.Printf("File too big: %s - %d bytes\n", httppath, filesize)
				return
			}

			if filemanager.FileAlreadyExists(fpath, filesize) {
				fmt.Printf("File already downloaded: %s\n", fpath)
				zipch <- fpath
				return
			}

			fmt.Printf("Downloading: %s\n", httppath)
			written, err := filemanager.DownloadFile(httppath, fpath)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Downloaded: %s - %d bytes written\n", fpath, written)

			zipch <- fpath
		}(link)
	}

	go func() {
		downloadWg.Wait()
		close(zipch)
	}()

	for fpath := range zipch {
		csvName := strings.Replace(filepath.Base(fpath), ".zip", ".csv", 1)
		csvPath := fmt.Sprintf("./_files/unzipped/%s", csvName)

		if filemanager.FileAlreadyExists(csvPath, -1) {
			fmt.Printf("File already unzipped: %s\n", csvPath)
			continue
		}

		fmt.Printf("Unzipping: %s\n", fpath)
		err = filemanager.UnzipFile(fpath, "./_files/unzipped")
		if err != nil {
			panic(err)
		}
	}
}

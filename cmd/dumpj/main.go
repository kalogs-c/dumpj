package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kalogs-c/dumpj/pkg/connections"
	"github.com/kalogs-c/dumpj/pkg/crawler"
	"github.com/kalogs-c/dumpj/pkg/entitites"
	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

// TODO
// - [ ] Handle errors (retry, log, etc)
// - [ ] Avoid duplicates and cache previous dump
// - [ ] Yaml config
// - [ ] CLI
// - [ ] Docker
// - [ ] Logging

func downloadFiles(links []string, filesPath string) <-chan string {
	out := make(chan string, 3)

	go func() {
		wg := sync.WaitGroup{}

		for _, link := range links {
			wg.Add(1)
			go func(link string) {
				defer wg.Done()

				httppath := fmt.Sprintf("%s%s", filesPath, link)
				fpath := fmt.Sprintf("./_files/zips/%s", link)

				filesize, err := filemanager.GetFileSize(httppath)
				if err != nil {
					log.Println(err)
				}

				if filemanager.FileAlreadyExists(fpath, filesize) {
					fmt.Printf("File already downloaded: %s\n", fpath)
					out <- fpath
					return
				}

				fmt.Printf("Downloading: %s\n", httppath)
				written, err := filemanager.DownloadFile(httppath, fpath)
				if err != nil {
					log.Println(err)
				}

				fmt.Printf("Downloaded: %s - %d bytes written\n", fpath, written)

				out <- fpath
			}(link)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func extractFiles(zipch <-chan string) <-chan string {
	out := make(chan string, 3)

	go func() {
		wg := sync.WaitGroup{}

		for fpath := range zipch {
			wg.Add(1)
			go func(fpath string) {
				defer wg.Done()

				csvName := strings.Replace(filepath.Base(fpath), ".zip", ".csv", 1)
				csvPath := fmt.Sprintf("./_files/unzipped/%s", csvName)

				if filemanager.FileAlreadyExists(csvPath, -1) {
					fmt.Printf("File already unzipped: %s\n", csvPath)
					out <- csvPath
					return
				}

				fmt.Printf("Unzipping: %s\n", fpath)
				err := filemanager.UnzipFile(fpath, "./_files/unzipped")
				if err != nil {
					log.Println(err)
				}

				fmt.Printf("Extracting: %s\n", fpath)
				out <- csvPath
			}(fpath)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func parseCsv(csvch <-chan string) <-chan entitites.Entity {
	out := make(chan entitites.Entity, 1024)

	go func() {
		wg := sync.WaitGroup{}

		for csvPath := range csvch {
			wg.Add(1)

			go func(csvPath string) {
				defer wg.Done()

				fmt.Printf("Parsing: %s\n", csvPath)
				csvFile, err := os.Open(csvPath)
				if err != nil {
					log.Println(err)
				}
				defer csvFile.Close()

				csvname := strings.Replace(filepath.Base(csvPath), ".csv", "", 1)
				for row, err := range filemanager.StreamCSV(csvFile, ';') {
					if err != nil {
						log.Println(err)
					}

					entity := entitites.PickEntity(csvname)
					err = entitites.NewEntityFromCSV(entity, row)
					if err != nil {
						log.Println(err)
					}

					if entity.IsValid() {
						out <- entity
					}
				}
			}(csvPath)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	q, closeConn, err := connections.NewSQLite(context.Background(), "file:./_files/dumpj.db?_fk=true&_journal_mode=WAL")
	if err != nil {
		panic(err)
	}
	defer closeConn()

	path := "https://dados-abertos-rf-cnpj.casadosdados.com.br/arquivos/"

	content, err := crawler.Fetch(path)
	if err != nil {
		log.Println(err)
	}

	reg := regexp.MustCompile("^20[0-9]{2}-[0-1][0-9]-[0-3][0-9]")
	links := crawler.ScrapeLinks(content, reg)
	latest := links[len(links)-1]

	files_path := fmt.Sprintf("%s%s", path, latest)
	content, err = crawler.Fetch(files_path)
	if err != nil {
		log.Println(err)
	}

	reg = regexp.MustCompile("Cnaes|Estabelecimentos|Municipios|Naturezas")
	links = crawler.ScrapeLinks(content, reg)

	zipch := downloadFiles(links, files_path)
	csvch := extractFiles(zipch)
	entitiesch := parseCsv(csvch)

	validEstabelecimentos := make(map[string]bool)
	globalCtx := context.WithValue(context.Background(), "validEstabelecimentos", validEstabelecimentos)
	for entity := range entitiesch {
		ctx, cancel := context.WithTimeout(globalCtx, 5*time.Second)
		defer cancel()

		err := entity.Save(ctx, q)
		if err != nil {
			log.Println(err)
		}
	}

	log.Printf("Fetching empresas...\n\tTotal Estabelecimentos: %d\n", len(validEstabelecimentos))

	reg = regexp.MustCompile("Empresas")
	links = crawler.ScrapeLinks(content, reg)

	zipch = downloadFiles(links, files_path)
	csvch = extractFiles(zipch)
	entitiesch = parseCsv(csvch)
	for entity := range entitiesch {
		ctx, cancel := context.WithTimeout(globalCtx, 5*time.Second)
		defer cancel()

		err := entity.Save(ctx, q)
		if err != nil {
			log.Println(err)
		}
	}
}

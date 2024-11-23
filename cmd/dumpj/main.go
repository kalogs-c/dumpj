package main

import (
	"fmt"
	"regexp"

	"github.com/kalogs-c/dumpj/pkg/crawler"
	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

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

	for _, link := range links {
		httppath := fmt.Sprintf("%s%s", files_path, link)
		fpath := fmt.Sprintf("./_files/zips/%s", link)

		fmt.Printf("Downloading: %s\n", httppath)
		_, err := filemanager.DownloadFile(httppath, fpath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Unzipping: %s\n", fpath)
		err = filemanager.UnzipFile(fpath, "./_files/unzipped")
		if err != nil {
			panic(err)
		}
	}
}

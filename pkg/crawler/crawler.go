package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

func Fetch(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, response.Status)
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func ScrapeLinks(content []byte, reg *regexp.Regexp) []string {
	body := bytes.NewReader(content)
	tokenizer := html.NewTokenizer(body)

	var links []string

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links
		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if token.Data != "a" {
				continue
			}

			for _, attr := range token.Attr {
				if attr.Key == "href" && reg.MatchString(attr.Val) {
					links = append(links, attr.Val)
				}
			}
		}
	}
}

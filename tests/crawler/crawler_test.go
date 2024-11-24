package crawler_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/kalogs-c/dumpj/pkg/crawler"
)

func TestFetch_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer mockServer.Close()

	content, err := crawler.Fetch(mockServer.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "Hello, World!"
	if string(content) != expected {
		t.Errorf("Expected %q, got %q", expected, string(content))
	}
}

func TestFetch_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	_, err := crawler.Fetch(mockServer.URL)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	expectedError := "GET " + mockServer.URL + ": 500 Internal Server Error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %q, got %q", expectedError, err.Error())
	}
}

func TestFetch_NetworkError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	mockServer.Close()

	_, err := crawler.Fetch(mockServer.URL)
	if err == nil {
		t.Fatal("Expected a network error, got nil")
	}
}

func TestScrapeLinks_ValidLinks(t *testing.T) {
	htmlContent := []byte(`
		<html>
			<body>
				<a href="https://example.com/page1">Page 1</a>
				<a href="https://example.org/page2">Page 2</a>
				<a href="mailto:someone@example.com">Email</a>
			</body>
		</html>`)

	reg := regexp.MustCompile(`^https?://`)

	links := crawler.ScrapeLinks(htmlContent, reg)
	expected := []string{
		"https://example.com/page1",
		"https://example.org/page2",
	}

	if len(links) != len(expected) {
		t.Fatalf("Expected %d links, got %d", len(expected), len(links))
	}

	for i, link := range links {
		if link != expected[i] {
			t.Errorf("Expected link %q, got %q", expected[i], link)
		}
	}
}

func TestScrapeLinks_NoMatches(t *testing.T) {
	htmlContent := []byte(`
		<html>
			<body>
				<a href="ftp://example.com/file">File</a>
				<a href="mailto:someone@example.com">Email</a>
			</body>
		</html>`)

	reg := regexp.MustCompile(`^https?://`)

	links := crawler.ScrapeLinks(htmlContent, reg)
	if len(links) != 0 {
		t.Fatalf("Expected no links, got %d", len(links))
	}
}

func TestScrapeLinks_EmptyContent(t *testing.T) {
	htmlContent := []byte("")
	reg := regexp.MustCompile(`.*`)

	links := crawler.ScrapeLinks(htmlContent, reg)
	if len(links) != 0 {
		t.Fatalf("Expected no links, got %d", len(links))
	}
}

func TestScrapeLinks_MalformedHTML(t *testing.T) {
	htmlContent := []byte(`<html><body><a href="https://example.com/page1">Page 1<a href="https://example.com/page2">Page 2</body></html>`)

	reg := regexp.MustCompile(`^https?://`)

	links := crawler.ScrapeLinks(htmlContent, reg)
	expected := []string{
		"https://example.com/page1",
		"https://example.com/page2",
	}

	if len(links) != len(expected) {
		t.Fatalf("Expected %d links, got %d", len(expected), len(links))
	}

	for i, link := range links {
		if link != expected[i] {
			t.Errorf("Expected link %q, got %q", expected[i], link)
		}
	}
}

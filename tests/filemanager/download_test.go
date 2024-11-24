package filemanager_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

func TestDownloadFile_Success(t *testing.T) {
	mockContent := "Hello, World!"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockContent))
	}))
	defer server.Close()

	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFilePath := tmpFile.Name()

	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	bytesWritten, err := filemanager.DownloadFile(server.URL, tmpFilePath)
	if err != nil {
		t.Fatalf("DownloadFile returned an error: %v", err)
	}

	expectedBytes := int64(len(mockContent))
	if bytesWritten != expectedBytes {
		t.Errorf("Expected %d bytes written, got %d", expectedBytes, bytesWritten)
	}

	data, err := os.ReadFile(tmpFilePath)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	if string(data) != mockContent {
		t.Errorf("Expected file content %q, got %q", mockContent, string(data))
	}
}

func TestDownloadFile_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	_, err := filemanager.DownloadFile(server.URL, "dummy_path")
	if err == nil {
		t.Fatal("Expected an error for HTTP 404 response, but got none")
	}
}

func TestDownloadFile_FileCreationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	}))
	defer server.Close()

	invalidPath := "/invalid_path/testfile.txt"

	_, err := filemanager.DownloadFile(server.URL, invalidPath)
	if err == nil {
		t.Fatal("Expected an error for file creation, but got none")
	}
}

package filemanager_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

func createTestZip(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "testzip-*.zip")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	zipWriter := zip.NewWriter(tmpFile)
	writer, err := zipWriter.Create(content)
	if err != nil {
		t.Fatalf("Failed to add file to zip: %v", err)
	}
	defer zipWriter.Close()

	_, err = writer.Write([]byte(content))
	if err != nil {
		t.Fatalf("Failed to write data to zip: %v", err)
	}

	return tmpFile.Name()
}

func TestUnzipFile_Success(t *testing.T) {
	zipPath := createTestZip(t, "Hello, World!")
	defer os.Remove(zipPath)

	destDir, err := os.MkdirTemp("", "unziptest-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	err = filemanager.UnzipFile(zipPath, destDir)
	if err != nil {
		t.Fatalf("UnzipFile returned an error: %v", err)
	}

	extractedFilename := strings.Replace(filepath.Base(zipPath), ".zip", ".csv", 1)
	extractedPath := filepath.Join(destDir, extractedFilename)
	data, err := os.ReadFile(extractedPath)
	if err != nil {
		t.Fatalf("Failed to read extracted file %s: %v", filepath.Base(extractedPath), err)
	}

	expectedData := "Hello, World!"
	if string(data) != expectedData {
		t.Errorf("File %s content mismatch: expected %q, got %q", filepath.Base(extractedPath), expectedData, string(data))
	}
}

func TestUnzipFile_InvalidZip(t *testing.T) {
	invalidZip := []byte("this is not a zip file")
	tmpFile, err := os.CreateTemp("", "invalidzip-*.zip")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write(invalidZip)
	tmpFile.Close()

	destDir, err := os.MkdirTemp("", "unziptest-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	err = filemanager.UnzipFile(tmpFile.Name(), destDir)
	if err == nil {
		t.Fatal("Expected an error for invalid zip file, but got none")
	}
}

func TestUnzipFile_EmptyZip(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "emptyzip-*.zip")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	zipWriter := zip.NewWriter(tmpFile)
	zipWriter.Close()

	destDir, err := os.MkdirTemp("", "unziptest-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	err = filemanager.UnzipFile(tmpFile.Name(), destDir)
	if err != nil {
		t.Fatalf("UnzipFile returned an error: %v", err)
	}

	files, err := os.ReadDir(destDir)
	if err != nil {
		t.Fatalf("Failed to read extraction directory: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected no files in extraction directory, found %d", len(files))
	}
}

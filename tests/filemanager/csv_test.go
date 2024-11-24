package filemanager_test

import (
	"strings"
	"testing"

	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

func TestStreamCSV_Success(t *testing.T) {
	csvContent := `name,age,city
Alice,30,New York
Bob,25,Los Angeles`

	reader := strings.NewReader(csvContent)
	notifier := make(chan filemanager.CSVRow)

	go filemanager.StreamCSV(notifier, reader, ',')

	var results []filemanager.CSVRow
	for row := range notifier {
		if row.Error != nil {
			t.Fatalf("Expected no errors, got %v", row.Error)
		}
		results = append(results, row)
	}

	expected := [][]string{
		{"name", "age", "city"},
		{"Alice", "30", "New York"},
		{"Bob", "25", "Los Angeles"},
	}

	if len(results) != len(expected) {
		t.Fatalf("Expected %d rows, got %d", len(expected), len(results))
	}

	for i, row := range results {
		if len(row.Data) != len(expected[i]) {
			t.Fatalf("Row %d length mismatch: expected %d, got %d", i, len(expected[i]), len(row.Data))
		}
		for j, cell := range row.Data {
			if cell != expected[i][j] {
				t.Errorf("Row %d, column %d: expected %q, got %q", i, j, expected[i][j], cell)
			}
		}
	}
}

func TestStreamCSV_ErrorHandling(t *testing.T) {
	csvContent := `name,age,city
Alice,30,New York
Bob,25,"Los Angeles`

	reader := strings.NewReader(csvContent)
	notifier := make(chan filemanager.CSVRow)

	go filemanager.StreamCSV(notifier, reader, ',')

	var errorOccurred bool
	for row := range notifier {
		if row.Error != nil {
			errorOccurred = true
			break
		}
	}

	if !errorOccurred {
		t.Fatal("Expected an error for malformed CSV row, but no error occurred")
	}
}

func TestStreamCSV_EmptyCSV(t *testing.T) {
	reader := strings.NewReader("")
	notifier := make(chan filemanager.CSVRow)

	go filemanager.StreamCSV(notifier, reader, ',')

	count := 0
	for range notifier {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 rows, got %d", count)
	}
}

func TestStreamCSV_ChannelClosing(t *testing.T) {
	reader := strings.NewReader("")
	notifier := make(chan filemanager.CSVRow)

	go filemanager.StreamCSV(notifier, reader, ',')

	_, open := <-notifier
	if open {
		t.Fatal("Expected channel to be closed after reading all data")
	}
}

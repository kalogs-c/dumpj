package filemanager_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

func TestStreamCSV_Success(t *testing.T) {
	csvContent := `name,age,city
Alice,30,New York
Bob,25,Los Angeles`

	reader := strings.NewReader(csvContent)

	var results [][]string
	for row, err := range filemanager.StreamCSV(reader, ',') {
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
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
		if len(row) != len(expected[i]) {
			t.Fatalf("Row %d length mismatch: expected %d, got %d", i, len(expected[i]), len(row))
		}
		for j, cell := range row {
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
	errorOccurred := false

	for _, err := range filemanager.StreamCSV(reader, ',') {
		if err != nil {
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

	count := 0
	for range filemanager.StreamCSV(reader, ',') {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 rows, got %d", count)
	}
}

type TestStruct struct {
	Name   string  `csv_column:"1"`
	Age    int     `csv_column:"2"`
	Salary float64 `csv_column:"3"`
}

func TestBindFields_Success(t *testing.T) {
	row := filemanager.CSVRow{
		Data: []string{"Alice", "30", "45000.75"},
	}

	var entity TestStruct
	err := filemanager.BindFields(row, &entity)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := TestStruct{Name: "Alice", Age: 30, Salary: 45000.75}
	if !reflect.DeepEqual(entity, expected) {
		t.Errorf("Expected %+v, got %+v", expected, entity)
	}
}

func TestBindFields_InvalidInteger(t *testing.T) {
	row := filemanager.CSVRow{
		Data: []string{"Alice", "invalid", "45000.75"},
	}

	var entity TestStruct
	err := filemanager.BindFields(row, &entity)
	if err == nil {
		t.Fatal("Expected an error for invalid integer, but got none")
	}
}

func TestBindFields_InvalidFloat(t *testing.T) {
	row := filemanager.CSVRow{
		Data: []string{"Alice", "30", "invalid"},
	}

	var entity TestStruct
	err := filemanager.BindFields(row, &entity)
	if err == nil {
		t.Fatal("Expected an error for invalid float, but got none")
	}
}

func TestBindFields_RowError(t *testing.T) {
	row := filemanager.CSVRow{
		Data:  []string{},
		Error: errors.New("mock error"),
	}

	var entity TestStruct
	err := filemanager.BindFields(row, &entity)
	if err == nil || err.Error() != "the row has an error: mock error" {
		t.Errorf("Expected row error, but got %v", err)
	}
}

func TestBindFields_InvalidColumnIndex(t *testing.T) {
	type InvalidStruct struct {
		InvalidField string `csv_column:"5"` // Out-of-range index
	}

	row := filemanager.CSVRow{
		Data: []string{"Alice", "30", "45000.75"},
	}

	var entity InvalidStruct
	err := filemanager.BindFields(row, &entity)
	if err == nil {
		t.Fatal("Expected an error for invalid column index, but got none")
	}
}

func TestBindFields_MissingTag(t *testing.T) {
	type PartialStruct struct {
		Name string `csv_column:"1"`
		City string // No tag
	}

	row := filemanager.CSVRow{
		Data: []string{"Alice", "30"},
	}

	var entity PartialStruct
	err := filemanager.BindFields(row, &entity)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := PartialStruct{Name: "Alice", City: ""} // City should remain zero value
	if !reflect.DeepEqual(entity, expected) {
		t.Errorf("Expected %+v, got %+v", expected, entity)
	}
}

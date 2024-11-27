package filemanager

import (
	"encoding/csv"
	"fmt"
	"io"
	"iter"
	"reflect"
	"strconv"
)

type CSVRow struct {
	Data  []string
	Error error
}

func StreamCSV(reader io.Reader, separator rune) iter.Seq2[[]string, error] {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = separator
	csvReader.TrimLeadingSpace = true

	return func(yield func([]string, error) bool) {
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				return
			}

			yield(record, err)
		}
	}
}

func BindFields[T any](row CSVRow, entity *T) error {
	if row.Error != nil {
		return fmt.Errorf("the row has an error: %v", row.Error)
	}

	entityValue := reflect.ValueOf(entity).Elem()
	entityType := entityValue.Type()
	if entityType.Kind() != reflect.Struct {
		return fmt.Errorf("entity must be a struct, got %T", entityType.Kind())
	}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		value := entityValue.Field(i)

		tag := field.Tag.Get("csv_column")
		if tag == "" {
			continue
		}

		columnIndex, err := strconv.Atoi(tag)
		if err != nil || columnIndex < 1 || columnIndex > len(row.Data) {
			return fmt.Errorf("invalid column index %q for field %s", tag, field.Name)
		}

		column := row.Data[columnIndex-1]
		if column == "" {
			value.Set(reflect.Zero(field.Type))
			continue
		}

		switch value.Kind() {
		case reflect.String:
			value.SetString(column)
		case reflect.Int:
			integer, err := strconv.Atoi(column)
			if err != nil {
				return fmt.Errorf("error converting %q to int for field %s: %v", column, field.Name, err)
			}
			value.SetInt(int64(integer))
		case reflect.Float64:
			float, err := strconv.ParseFloat(column, 64)
			if err != nil {
				return fmt.Errorf("error converting %q to float for field %s: %v", column, field.Name, err)
			}
			value.SetFloat(float)
		default:
			return fmt.Errorf("unsupported type %s for field %s", value.Kind(), field.Name)
		}
	}

	return nil
}

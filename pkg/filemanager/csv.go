package filemanager

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"iter"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func StreamCSV(reader io.Reader, separator rune) iter.Seq2[[]string, error] {
	csvReader := csv.NewReader(transform.NewReader(reader, charmap.ISO8859_1.NewDecoder()))
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

func BindFields(row []string, entity any) error {
	entityValue := reflect.ValueOf(entity).Elem()
	entityType := entityValue.Type()
	if entityType.Kind() != reflect.Struct {
		return fmt.Errorf("entity must be a struct, got %T", entityType.Kind())
	}

	numFields := entityType.NumField()
	for i := 0; i < numFields; i++ {
		field := entityType.Field(i)
		value := entityValue.Field(i)

		tag := field.Tag.Get("csv_column")
		if tag == "" {
			continue
		}

		columnIndex, err := strconv.Atoi(tag)
		if err != nil || columnIndex < 1 || columnIndex > len(row) {
			return fmt.Errorf("invalid column index %q for field %s", tag, field.Name)
		}

		if !value.CanSet() {
			continue
		}

		column := row[columnIndex-1]

		switch value.Type() {
		case reflect.TypeOf(sql.NullString{}):
			value.Set(reflect.ValueOf(sql.NullString{String: column, Valid: column != ""}))
			continue

		case reflect.TypeOf(time.Time{}):
			t, err := time.Parse("20060102", column)
			if err != nil {
				return fmt.Errorf("error parsing %s to time for field %s: %v", column, field.Name, err)
			}
			value.Set(reflect.ValueOf(t))
			continue
		}

		if column == "" {
			value.Set(reflect.Zero(field.Type))
			continue
		}

		switch value.Kind() {
		case reflect.String:
			value.SetString(column)

		case reflect.Int:
		case reflect.Int8:
		case reflect.Int16:
		case reflect.Int32:
		case reflect.Int64:
			if strings.Contains(column, ",") {
				column = strings.Split(column, ",")[0]
			}

			integer, err := strconv.Atoi(column)
			if err != nil {
				return fmt.Errorf("error converting %s to int for field %s: %v", column, field.Name, err)
			}
			value.SetInt(int64(integer))

		case reflect.Float64:
			float, err := strconv.ParseFloat(strings.Replace(column, ",", ".", 1), 64)
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

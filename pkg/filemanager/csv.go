package filemanager

import (
	"encoding/csv"
	"io"
)

type CSVRow struct {
	Data  []string
	Error error
}

func StreamCSV(notifier chan<- CSVRow, reader io.Reader, separator rune) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = separator
	csvReader.TrimLeadingSpace = true

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			close(notifier)
			return
		}

		row := CSVRow{Data: record, Error: nil}
		if err != nil {
			row.Error = err
		}

		notifier <- row
	}
}

package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
)

//counterfeiter:generate . CSVFileWriter
type CSVFileWriter interface {
	Write(csvFilePath string, csvHeader []string, csvRecords [][]string) error
}

type csvWriter struct {}

func NewCSVFileWriter() CSVFileWriter {
	return csvWriter{}
}

func (w csvWriter) Write(csvFilePath string, csvHeader []string, csvRecords [][]string) error {
	file, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to create csv file: %v", err)
	}
	defer file.Close()

	utf16Writer := transform.NewWriter(file, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder())
	csvWriter := csv.NewWriter(utf16Writer)
	// Set the delimiter to semicolon
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	if err := csvWriter.Write(csvHeader); err != nil {
		return fmt.Errorf("failed to write csv header: %v", err)
	}

	for _, record := range csvRecords {
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write csv records: %v", err)
		}
	}

	return nil
}

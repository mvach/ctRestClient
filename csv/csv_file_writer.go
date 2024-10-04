package csv

import (
	"encoding/csv"
	"fmt"
	"os"
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

	// Write the UTF-8 BOM for Excel on Windows compatibility
	_, err = file.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		return fmt.Errorf("failed to write UTF-8 BOM to csv file: %v", err)
	}

	writer := csv.NewWriter(file)
	// Set the delimiter to semicolon
	writer.Comma = ';'
	defer writer.Flush()

	if err := writer.Write(csvHeader); err != nil {
		return fmt.Errorf("failed to write csv header: %v", err)
	}

	for _, record := range csvRecords {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write csv records: %v", err)
		}
	}

	return nil
}

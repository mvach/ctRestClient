package app

import (
	"encoding/csv"
	"fmt"
	"os"
)

//counterfeiter:generate . CSVWriter
type CSVWriter interface {
	Write(csvFilePath string, csvHeader []string, csvRecords [][]string) error
}

type csvWriter struct {}

func NewCSVWriter() CSVWriter {
	return csvWriter{}
}

func (w csvWriter) Write(csvFilePath string, csvHeader []string, csvRecords [][]string) error {
	file, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to create csv file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
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

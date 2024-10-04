package csv

import (
	"ctRestClient/logger"
	"encoding/json"
	"fmt"
)

type personData struct {
	header  []string
	records [][]string
}

func NewPersonData(persons []json.RawMessage, fields []string, logger logger.Logger) (CsvData, error) {
	csvRecords := make([][]string, 0)

	for _, person := range persons {
		var data map[string]interface{}
		err := json.Unmarshal([]byte(person), &data)
		if err != nil {
			return nil, fmt.Errorf("failed to read person information raw json: %v", err)
		}

		record := make([]string, len(fields))

		for i, field := range fields {
			if value, ok := data[field].(string); ok {
				// get string values
				record[i] = value
			} else if value, ok := data[field].(float64); ok {
				// get int values
				record[i] = fmt.Sprintf("%d", int(value))
			} else {
				logger.Warn(fmt.Sprintf("    Field '%s' is not a string or int, or not found", field))
				record[i] = ""
			}
		}
		csvRecords = append(csvRecords, record)
	}

	return &personData{
		header:  fields,
		records: csvRecords,
	}, nil
}

func (p *personData) Records() [][]string {
	return p.records
}

func (p *personData) Header() []string {
	return p.header
}

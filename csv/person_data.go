package csv

import (
	"ctRestClient/config"
	"ctRestClient/logger"
	"encoding/json"
	"fmt"
)

type personData struct {
	header  []string
	records [][]string
}

func NewPersonData(persons []json.RawMessage, fields []config.Field, personDataProvider FileDataProvider, logger logger.Logger) (CsvData, error) {
	csvRecords := make([][]string, 0)

	for _, person := range persons {
		var jsonData map[string]json.RawMessage
		err := json.Unmarshal([]byte(person), &jsonData)
		if err != nil {
			return nil, fmt.Errorf("failed to read person information raw json: %v", err)
		}

		record := make([]string, len(fields))

		for i, field := range fields {
			fieldName := field.GetFieldName()
			rawValue, exists := jsonData[fieldName]

			value := ""

			if !exists {
				logger.Warn(fmt.Sprintf("    Field '%s' does not exist", fieldName))
				record[i] = ""
			} else if rawValue == nil {
				record[i] = ""
			} else {
				if !field.IsMappedData() {
					value = convertToString(rawValue)
				} else {
					value, err = personDataProvider.GetData(fieldName, rawValue)
					if err != nil {
						logger.Error(fmt.Sprintf("     failed to get data for field '%s': %v", fieldName, err))
						value = ""
					}
				}
			}
			record[i] = value

		}
		csvRecords = append(csvRecords, record)
	}

	// Extract field names for the header
	csvHeader := make([]string, len(fields))
	for i, field := range fields {
		csvHeader[i] = field.GetColumnName()
	}

	return &personData{
		header:  csvHeader,
		records: csvRecords,
	}, nil
}

// Helper function to convert JSON values to strings
func convertToString(value json.RawMessage) string {
	// Parse the raw message to get the actual value
	var parsedValue interface{}
	if err := json.Unmarshal(value, &parsedValue); err != nil {
		// If parsing fails, return the raw string
		return string(value)
	}

	switch v := parsedValue.(type) {
	case string:
		// For strings, return the value without quotes
		return v
	case float64:
		// For numbers, use the original format from raw message
		return string(value)
	case int:
		return string(value)
	case bool:
		return string(value)
	case nil:
		return ""
	default:
		// For other types, marshal to get string representation
		if jsonBytes, err := json.Marshal(v); err == nil {
			return string(jsonBytes)
		}
		return fmt.Sprintf("%v", v)
	}
}

func (p *personData) Records() [][]string {
	return p.records
}

func (p *personData) Header() []string {
	return p.header
}

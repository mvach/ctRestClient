package csv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//counterfeiter:generate . FileDataProvider
type FileDataProvider interface {
	GetData(ctFieldName string, ctFieldValue json.RawMessage) (string, error)
}

type typedValue struct {
	Tag   string
	Value string
}

type fileDataProvider struct {
	dataDir   string
	dataCache map[string]map[typedValue]string
}

func NewFileDataProvider(dataDir string) FileDataProvider {
	return &fileDataProvider{
		dataDir:   dataDir,
		dataCache: make(map[string]map[typedValue]string),
	}
}

func (dp *fileDataProvider) GetData(ctFieldName string, ctFieldValue json.RawMessage) (string, error) {
	dataFilePath := filepath.Join(dp.dataDir, ctFieldName+".yml")

	data, exists := dp.dataCache[ctFieldName]

	if !exists {
		yamlData, err := os.ReadFile(dataFilePath)
		if err != nil {
			return "", err
		}

		var yamlNode yaml.Node
		if err := yaml.Unmarshal(yamlData, &yamlNode); err != nil {
			return "", err
		}

		dataMap := make(map[typedValue]string)

		if yamlNode.Kind == yaml.DocumentNode && len(yamlNode.Content) > 0 {
			mapNode := yamlNode.Content[0]
			if mapNode.Kind == yaml.MappingNode {
				for i := 0; i < len(mapNode.Content); i += 2 {
					keyNode := mapNode.Content[i]
					valueNode := mapNode.Content[i+1]

					dataMap[typedValue{
						Tag:   keyNode.Tag,
						Value: keyNode.Value,
					}] = valueNode.Value
				}
			}
		}

		// Fill the cache with the YAML data
		dp.dataCache[ctFieldName] = dataMap
		data = dataMap
	}

	typedValue := dp.createYamlKeyFromJSON(ctFieldValue)

	if mappedValue, exists := data[typedValue]; exists {
		return mappedValue, nil
	}

	return "", fmt.Errorf("the value %s is not in '%s'", typedValue.Value, dataFilePath)
}

func (dp *fileDataProvider) createYamlKeyFromJSON(value json.RawMessage) typedValue {
	// Parse to determine type
	var parsedValue interface{}
	json.Unmarshal(value, &parsedValue)

	switch parsedValue.(type) {
	case string:
		// Remove quotes from the raw message
		trimmed := strings.Trim(string(value), "\"")
		return typedValue{Tag: "!!str", Value: trimmed}
	case float64:
		// Use the original raw message to preserve format
		originalValue := string(value)
		if strings.Contains(originalValue, ".") {
			return typedValue{Tag: "!!float", Value: originalValue}
		} else {
			return typedValue{Tag: "!!int", Value: originalValue}
		}
	case int:
		return typedValue{Tag: "!!int", Value: string(value)}
	case bool:
		return typedValue{Tag: "!!bool", Value: string(value)}
	default:
		return typedValue{Tag: "!!str", Value: string(value)}
	}
}

package config

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Group struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

func (g Group) CSVFileName() string {
	return g.sanitizedGroupName() + ".csv"
}

func (g Group) BlocklistFileName() string {
	return g.sanitizedGroupName() + ".yml"
}

func (g Group) sanitizedGroupName() string {
	fileName := g.Name
	fileName = strings.ReplaceAll(fileName, " ", "_")
	fileName = strings.ReplaceAll(fileName, ",", ".")
	fileName = strings.ReplaceAll(fileName, "ä", "ae")
	fileName = strings.ReplaceAll(fileName, "ö", "oe")
	fileName = strings.ReplaceAll(fileName, "ü", "ue")
	fileName = strings.ReplaceAll(fileName, "Ä", "Ae")
	fileName = strings.ReplaceAll(fileName, "Ö", "Oe")
	fileName = strings.ReplaceAll(fileName, "Ü", "Ue")

	re := regexp.MustCompile(`[^\w\-.]`)
	fileName = re.ReplaceAllString(fileName, "")
	fileName = strings.ReplaceAll(fileName, "__", "_")

	return fileName
}

func (g *Group) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Define a temporary type to capture the raw YAML structure
	type rawGroup struct {
		Name   string                   `yaml:"name"`
		Fields []map[string]interface{} `yaml:"fields"`
	}

	raw := rawGroup{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	g.Name = raw.Name
	g.Fields = make([]Field, len(raw.Fields))

	for i, rawField := range raw.Fields {
		// Check if the field has an "object" key (ExportGroupField)
		if _, hasObject := rawField["object"]; hasObject {
			var exportField ExportGroupField
			// Re-marshal the raw field into ExportGroupField
			fieldYAML, err := yaml.Marshal(rawField)
			if err != nil {
				return err
			}
			if err := yaml.Unmarshal(fieldYAML, &exportField); err != nil {
				return err
			}
			g.Fields[i] = &exportField
		} else {
			// Otherwise, treat it as a CreateGroupField
			var createField CreateGroupField
			// Re-marshal the raw field into CreateGroupField
			fieldYAML, err := yaml.Marshal(rawField)
			if err != nil {
				return err
			}
			if err := yaml.Unmarshal(fieldYAML, &createField); err != nil {
				return err
			}
			g.Fields[i] = &createField
		}
	}

	return nil
}

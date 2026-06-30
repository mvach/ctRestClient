package config

import (
	"regexp"
	"strings"
	"errors"
)

type ExportGroupGroup struct {
	Name   string  `yaml:"name"`
	Fields []ExportGroupField `yaml:"fields"`
}

func (g ExportGroupGroup) CSVFileName() string {
	return g.sanitizedGroupName() + ".csv"
}

func (g ExportGroupGroup) BlocklistFileName() string {
	return g.sanitizedGroupName() + ".yml"
}

func (g ExportGroupGroup) sanitizedGroupName() string {
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

func (g ExportGroupGroup) Validate() error {
	if g.Name == "" {
		return errors.New("property name is not set")
	}
	if len(g.Fields) == 0 {
		return errors.New("property fields is not set")
	}
	return nil
}

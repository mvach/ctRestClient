package config

import (
	"regexp"
	"strings"
)

type Group struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

func (g Group) SanitizedGroupName() string {
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

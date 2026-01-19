package testutil

import (
	"strings"

	"github.com/lithammer/dedent"
)

func YamlToString(yamlString string) string {
	yamlString = dedent.Dedent(yamlString)
	yamlString = strings.ReplaceAll(yamlString, "\t", "  ")
	yamlString = strings.TrimSpace(yamlString)
	return yamlString
}

func YamlToByteArray(yamlString string) []byte {
	return []byte(YamlToString(yamlString))
}

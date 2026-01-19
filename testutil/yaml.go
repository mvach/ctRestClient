package testutil

import (
	"strings"

	"github.com/lithammer/dedent"
)

func YamlToByteArray(yamlString string) []byte {
	yamlString = dedent.Dedent(yamlString)
	yamlString = strings.ReplaceAll(yamlString, "\t", "  ")
	yamlString = strings.TrimSpace(yamlString)

	return []byte(yamlString)
}

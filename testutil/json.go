package testutil

import (
	"bytes"
	"strings"

	"github.com/lithammer/dedent"
)

func JsonToBufferString(jsonString string)  *bytes.Buffer {
	jsonString = dedent.Dedent(jsonString)
	jsonString = strings.ReplaceAll(jsonString, "\t", "  ")
	jsonString = strings.TrimSpace(jsonString)

	return bytes.NewBuffer([]byte(jsonString))
}

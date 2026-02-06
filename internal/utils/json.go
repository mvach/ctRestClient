package utils

import (
	"encoding/json"
	"fmt"
)

func UnescapeUnicodeCharacters(jsonRaw json.RawMessage) (json.RawMessage, error) {
	var temp interface{}
	if err := json.Unmarshal(jsonRaw, &temp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json raw message: %v", err)
	}

	result, err := json.Marshal(temp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json to raw message: %v", err)
	}
	return result, nil
}
package restendpoints

import (
	"encoding/json"
)

type PersonResponseJson struct {
	Data json.RawMessage `json:"data"`
}

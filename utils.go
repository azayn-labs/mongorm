package mongorm

import (
	"encoding/json"
)

func clonePtr[T any](src *T, reset bool) *T {
	b, err := json.Marshal(src)
	if err != nil {
		return nil
	}

	var dst T
	if err := json.Unmarshal(b, &dst); err != nil {
		return nil
	}

	if reset {
		var zero T
		return &zero
	}

	return &dst
}

// Function to check if the JSON contains a field.
// Returns the content of the field and bool if the field exists and is not nil.
func jsonContainsField(jsonData []byte, field string) (any, bool) {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, false
	}

	value, exists := data[field]
	if !exists || value == nil {
		return nil, false
	}

	return value, true
}

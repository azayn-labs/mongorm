package mongorm

import (
	"encoding/json"
)

// clonePtr creates a deep copy of the provided struct pointer using JSON marshaling and
// unmarshaling. If the reset parameter is true, it returns a pointer to a zero-value instance
// of the struct. This function is useful for creating new instances of the schema or MongORM
// struct without copying the existing data.
//
// > NOTE: This function is internal only.
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

// MongORMOptions holds the configuration options for a MongORM instance, including settings for
// timestamps, collection and database names, and the MongoDB client. This struct is used to
// customize the behavior of the MongORM instance when connecting to the database and performing
// operations.
//
// > NOTE: This function is internal only.
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

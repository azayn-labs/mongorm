package mongorm

import (
	"encoding/json"
	"reflect"
	"strings"
)

func convertToStruct[T any](input any) (T, error) {
	var result T

	switch v := input.(type) {
	case string:
		err := json.Unmarshal([]byte(v), &result)
		return result, err
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return result, err
		}
		err = json.Unmarshal(bytes, &result)
		return result, err
	}
}

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

func hasModelFlag(tag reflect.StructTag, flag string) bool {
	v := tag.Get("mongorm")
	if v == "" {
		return false
	}

	for _, f := range strings.Split(v, ",") {
		if strings.TrimSpace(f) == flag {
			return true
		}
	}

	return false
}

func getBSONName(f reflect.StructField) (name string, inline bool, ok bool) {
	tag := f.Tag.Get("bson")

	if tag == "-" {
		return "", false, false
	}

	if tag == "" {
		return f.Name, false, true
	}

	parts := strings.Split(tag, ",")

	name = parts[0]

	if name == "" {
		name = f.Name
	}

	for _, p := range parts[1:] {
		if p == "inline" {
			return "", true, true
		}
	}

	return name, false, true
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

func fieldName(sf reflect.StructField) string {
	tag := sf.Tag.Get("bson")
	if tag == "" {
		return sf.Name
	}

	return strings.Split(tag, ",")[0]
}

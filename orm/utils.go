package orm

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

func getBSONName(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("bson")
	if tag == "-" {
		return "", false
	}

	if tag == "" {
		return f.Name, true
	}

	name := tag
	if idx := strings.Index(tag, ","); idx >= 0 {
		name = tag[:idx]
	}

	if name == "" {
		name = f.Name
	}

	return name, true
}

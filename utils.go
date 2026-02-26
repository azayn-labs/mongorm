package mongorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func castDistinctValue[V any](value any, targetType reflect.Type) (V, error) {
	var zero V

	if value == nil {
		return zero, fmt.Errorf("value is nil")
	}

	if direct, ok := value.(V); ok {
		return direct, nil
	}

	switch any(zero).(type) {
	case int64:
		switch typed := value.(type) {
		case int64:
			return any(typed).(V), nil
		case int32:
			return any(int64(typed)).(V), nil
		case int16:
			return any(int64(typed)).(V), nil
		case int8:
			return any(int64(typed)).(V), nil
		case int:
			return any(int64(typed)).(V), nil
		}
	case float64:
		switch typed := value.(type) {
		case float64:
			return any(typed).(V), nil
		case float32:
			return any(float64(typed)).(V), nil
		case int64:
			return any(float64(typed)).(V), nil
		case int32:
			return any(float64(typed)).(V), nil
		case int16:
			return any(float64(typed)).(V), nil
		case int8:
			return any(float64(typed)).(V), nil
		case int:
			return any(float64(typed)).(V), nil
		}
	case bson.ObjectID:
		switch typed := value.(type) {
		case bson.ObjectID:
			return any(typed).(V), nil
		case string:
			id, err := bson.ObjectIDFromHex(typed)
			if err != nil {
				return zero, fmt.Errorf("invalid objectid hex: %w", err)
			}
			return any(id).(V), nil
		}
	case time.Time:
		switch typed := value.(type) {
		case time.Time:
			return any(typed).(V), nil
		case interface{ Time() time.Time }:
			return any(typed.Time()).(V), nil
		case int64:
			return any(time.UnixMilli(typed)).(V), nil
		}
	}

	if targetType != nil {
		raw := reflect.ValueOf(value)
		if raw.IsValid() && raw.Type().ConvertibleTo(targetType) {
			converted := raw.Convert(targetType).Interface()
			if typed, ok := converted.(V); ok {
				return typed, nil
			}
		}
	}

	return zero, fmt.Errorf("value type %T cannot be converted", value)
}

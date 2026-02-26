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
			return assertDistinctType[V](typed)
		case int32:
			return assertDistinctType[V](int64(typed))
		case int16:
			return assertDistinctType[V](int64(typed))
		case int8:
			return assertDistinctType[V](int64(typed))
		case int:
			return assertDistinctType[V](int64(typed))
		}
	case float64:
		switch typed := value.(type) {
		case float64:
			return assertDistinctType[V](typed)
		case float32:
			return assertDistinctType[V](float64(typed))
		case int64:
			return assertDistinctType[V](float64(typed))
		case int32:
			return assertDistinctType[V](float64(typed))
		case int16:
			return assertDistinctType[V](float64(typed))
		case int8:
			return assertDistinctType[V](float64(typed))
		case int:
			return assertDistinctType[V](float64(typed))
		}
	case bson.ObjectID:
		switch typed := value.(type) {
		case bson.ObjectID:
			return assertDistinctType[V](typed)
		case string:
			id, err := bson.ObjectIDFromHex(typed)
			if err != nil {
				return zero, fmt.Errorf("invalid objectid hex: %w", err)
			}
			return assertDistinctType[V](id)
		}
	case time.Time:
		switch typed := value.(type) {
		case time.Time:
			return assertDistinctType[V](typed)
		case interface{ Time() time.Time }:
			return assertDistinctType[V](typed.Time())
		case int64:
			return assertDistinctType[V](time.UnixMilli(typed))
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

func assertDistinctType[V any](value any) (V, error) {
	var zero V

	typed, ok := value.(V)
	if !ok {
		return zero, fmt.Errorf("value type %T cannot be asserted", value)
	}

	return typed, nil
}

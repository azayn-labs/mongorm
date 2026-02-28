package mongorm

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// clonePtr creates a copy of the provided pointer target by value. If reset is true,
// it returns a pointer to a zero-value instance of T.
//
// > NOTE: This function is internal only.
func clonePtr[T any](src *T, reset bool) *T {
	if reset {
		var zero T
		return &zero
	}

	if src == nil {
		return nil
	}

	dst := *src

	return &dst
}

func castDistinctValue[V any](value any, targetType reflect.Type) (V, error) {
	var zero V

	if value == nil {
		return zero, configErrorf("value is nil")
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
				return zero, configErrorf("invalid objectid hex: %v", err)
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

	return zero, configErrorf("value type %T cannot be converted", value)
}

func assertDistinctType[V any](value any) (V, error) {
	var zero V

	typed, ok := value.(V)
	if !ok {
		return zero, configErrorf("value type %T cannot be asserted", value)
	}

	return typed, nil
}

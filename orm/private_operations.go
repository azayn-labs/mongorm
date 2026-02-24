package orm

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *Model[T]) applySchema(doc *T) error {
	m.schema = clonePtr(doc, false)
	m.clone = clonePtr(doc, false)

	return nil
}

func (m *Model[T]) updateSchema(ctx context.Context, id bson.ObjectID) error {
	if err := m.collection.FindOne(ctx, BaseModel{
		ID: &id,
	}).Decode(m.schema); err != nil {
		return err
	}

	// Refresh the clone with the updated schema
	m.clone = clonePtr(m.schema, false)

	return nil
}

func (m *Model[T]) getInformationToSet() bson.M {
	out := bson.M{}

	if m == nil || m.schema == nil || m.clone == nil {
		return out
	}

	src := reflect.ValueOf(m.schema).Elem()
	dst := reflect.ValueOf(m.clone).Elem()
	t := dst.Type()

	for i := 0; i < src.NumField(); i++ {
		fieldType := t.Field(i)

		// Skip unexported
		if fieldType.PkgPath != "" {
			continue
		}

		// Skip protected
		if hasModelFlag(fieldType.Tag, string(ModelTagPrimary)) ||
			hasModelFlag(fieldType.Tag, string(ModelTagReadonly)) {
			continue
		}

		key, ok := getBSONName(fieldType)
		if !ok {
			continue
		}

		sf := src.Field(i)
		df := dst.Field(i)

		if !sf.IsValid() || !df.IsValid() {
			continue
		}

		// ---- Pointer ↔ Pointer ----
		if sf.Kind() == reflect.Pointer && df.Kind() == reflect.Pointer {

			switch {
			case sf.IsNil() && df.IsNil():
				continue

			case sf.IsNil() && !df.IsNil():
				out[key] = df.Interface()
				continue

			case !sf.IsNil() && df.IsNil():
				continue

			default:
				if !reflect.DeepEqual(sf.Elem().Interface(), df.Elem().Interface()) {
					out[key] = df.Interface()
				}
				continue
			}
		}

		// ---- Pointer ↔ Value mismatch ----
		if sf.Kind() == reflect.Pointer && df.Kind() != reflect.Pointer {
			if sf.IsNil() || !reflect.DeepEqual(sf.Elem().Interface(), df.Interface()) {
				out[key] = df.Interface()
			}
			continue
		}

		if sf.Kind() != reflect.Pointer && df.Kind() == reflect.Pointer {
			if !df.IsNil() && !reflect.DeepEqual(sf.Interface(), df.Elem().Interface()) {
				out[key] = df.Interface()
			}
			continue
		}

		// ---- Slice / Map / Interface ----
		switch sf.Kind() {
		case reflect.Slice, reflect.Map, reflect.Interface:
			if !reflect.DeepEqual(sf.Interface(), df.Interface()) {
				out[key] = df.Interface()
			}
			continue
		}

		// ---- Value types ----
		if !reflect.DeepEqual(sf.Interface(), df.Interface()) {
			out[key] = df.Interface()
		}
	}

	return out
}

func (m *Model[T]) getInformationToUnset() bson.M {
	out := bson.M{}

	if m == nil || m.schema == nil || m.clone == nil {
		return out
	}

	src := reflect.ValueOf(m.schema).Elem()
	dst := reflect.ValueOf(m.clone).Elem()
	t := dst.Type()

	for i := 0; i < src.NumField(); i++ {
		fieldType := t.Field(i)

		if fieldType.PkgPath != "" {
			continue
		}

		if hasModelFlag(fieldType.Tag, string(ModelTagPrimary)) ||
			hasModelFlag(fieldType.Tag, string(ModelTagReadonly)) {
			continue
		}

		key, ok := getBSONName(fieldType)
		if !ok {
			continue
		}

		sf := src.Field(i)
		df := dst.Field(i)

		if !sf.IsValid() || !df.IsValid() {
			continue
		}

		// ---- Pointer removed ----
		if sf.Kind() == reflect.Pointer && df.Kind() == reflect.Pointer {
			if !sf.IsNil() && df.IsNil() {
				out[key] = 1
			}
			continue
		}

		// ---- Pointer ↔ Value mismatch ----
		if sf.Kind() == reflect.Pointer && df.Kind() != reflect.Pointer {
			if !sf.IsNil() {
				out[key] = 1
			}
			continue
		}

		if sf.Kind() != reflect.Pointer && df.Kind() == reflect.Pointer {
			if df.IsNil() && !sf.IsZero() {
				out[key] = 1
			}
			continue
		}

		// ---- Slice / Map ----
		switch sf.Kind() {
		case reflect.Slice, reflect.Map:
			if !sf.IsNil() && (df.IsNil() || df.Len() == 0) {
				out[key] = 1
			}
			continue

		case reflect.Interface:
			if !sf.IsNil() && df.IsNil() {
				out[key] = 1
			}
			continue
		}

		// ---- Value → zero ----
		if !sf.IsZero() && df.IsZero() {
			out[key] = 1
		}
	}

	return out
}

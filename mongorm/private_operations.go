package mongorm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *Model[T]) applySchema(doc *T) error {
	*m.schema = *doc
	m.clone = clonePtr(doc, false)

	return nil
}

func (m *Model[T]) updateSchema(ctx context.Context, id *bson.ObjectID) error {
	_, fieldID, err := m.getPrimaryField()
	if err != nil {
		return err
	}

	if err := m.collection.FindOne(ctx, bson.M{
		fieldID: id,
	}).Decode(m.schema); err != nil {
		return err
	}

	// Refresh the clone with the updated schema
	m.clone = clonePtr(m.schema, false)

	return nil
}

func (m *Model[T]) processEmbeddedSet(sf, df reflect.Value, out bson.M) {
	if sf.Kind() == reflect.Pointer {
		if sf.IsNil() {
			return
		}
		sf = sf.Elem()
	}

	if df.Kind() == reflect.Pointer {
		if df.IsNil() {
			return
		}
		df = df.Elem()
	}

	for i := 0; i < sf.NumField(); i++ {
		fType := df.Type().Field(i)

		if fType.PkgPath != "" {
			continue
		}

		key, inline, ok := getBSONName(fType)
		if !ok {
			continue
		}

		subSF := sf.Field(i)
		subDF := df.Field(i)

		if inline {
			m.processEmbeddedSet(subSF, subDF, out)
			continue
		}

		if !reflect.DeepEqual(subSF.Interface(), subDF.Interface()) {
			out[key] = subDF.Interface()
		}
	}
}

func (m *Model[T]) processEmbeddedUnset(sf, df reflect.Value, out bson.M) {
	if sf.Kind() == reflect.Pointer {
		if sf.IsNil() {
			return
		}
		sf = sf.Elem()
	}

	if df.Kind() == reflect.Pointer {
		if df.IsNil() {
			// whole embedded struct removed → unset all fields
			for i := 0; i < sf.NumField(); i++ {
				fType := sf.Type().Field(i)
				key, inline, ok := getBSONName(fType)
				if ok && !inline {
					out[key] = 1
				}
			}
			return
		}
		df = df.Elem()
	}

	for i := 0; i < sf.NumField(); i++ {
		fType := df.Type().Field(i)

		if fType.PkgPath != "" {
			continue
		}

		key, inline, ok := getBSONName(fType)
		if !ok {
			continue
		}

		subSF := sf.Field(i)
		subDF := df.Field(i)

		if inline {
			m.processEmbeddedUnset(subSF, subDF, out)
			continue
		}

		// ---- Pointer removed ----
		if subSF.Kind() == reflect.Pointer && subDF.Kind() == reflect.Pointer {
			if !subSF.IsNil() && subDF.IsNil() {
				out[key] = 1
			}
			continue
		}

		// ---- Value → zero ----
		if !subSF.IsZero() && subDF.IsZero() {
			out[key] = 1
		}
	}
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

		key, inline, ok := getBSONName(fieldType)
		if !ok {
			continue
		}

		sf := src.Field(i)
		df := dst.Field(i)

		if !sf.IsValid() || !df.IsValid() {
			continue
		}

		if inline {
			m.processEmbeddedSet(sf, df, out)
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

		key, inline, ok := getBSONName(fieldType)
		if !ok {
			continue
		}

		sf := src.Field(i)
		df := dst.Field(i)

		if !sf.IsValid() || !df.IsValid() {
			continue
		}

		if inline {
			m.processEmbeddedUnset(sf, df, out)
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

func (m *Model[T]) insertOne(ctx context.Context) error {
	ins, err := m.collection.InsertOne(ctx, m.clone)
	if err != nil {
		return err
	}

	id, ok := ins.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("Invalid document from database: missing identifier")
	}

	return m.updateSchema(ctx, &id)
}

func (m *Model[T]) updateOne(ctx context.Context, id *bson.ObjectID) error {
	var doc T
	set := m.getInformationToSet()
	unset := m.getInformationToUnset()

	delete(set, "createdAt") // Never update createdAt

	toDo := bson.M{}
	if len(set) > 0 {
		toDo["$set"] = set
	}
	if len(unset) > 0 {
		toDo["$unset"] = unset
	}

	if len(toDo) == 0 {
		return nil
	}

	_, fieldID, err := m.getPrimaryField()
	if err != nil {
		return err
	}

	if err := m.collection.FindOneAndUpdate(
		ctx,
		bson.M{
			fieldID: id,
		},
		toDo,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc); err != nil {
		return err
	}

	return m.applySchema(&doc)
}

func (m *Model[T]) getPrimaryField() (string, string, error) {
	v := reflect.ValueOf(m.schema).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)

		if hasModelFlag(fieldType.Tag, string(ModelTagPrimary)) {
			return fieldType.Name, strings.Split(fieldType.Tag.Get("bson"), ",")[0], nil
		}
	}

	return "", "", fmt.Errorf("No primary field found")
}

func (m *Model[T]) findOne(ctx context.Context) error {
	var doc T
	if err := m.collection.FindOne(ctx, m.query).Decode(&doc); err != nil {
		return err
	}

	return m.applySchema(&doc)
}

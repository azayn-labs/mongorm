package mongorm

import (
	"reflect"
	"strings"

	"github.com/CdTgr/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Field represents a field in the schema of a MongORM instance. It provides methods to
// retrieve the BSON name of the field, which is used for constructing queries and updates.
// The Field interface is implemented by specific field types (e.g., StringField, Int64Field)
// that correspond to the Go types used in the schema struct.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	}
//	var ToDoFields = mongorm.FieldsOf[ToDo, struct {
//	  Text *primitives.StringField
//	}]()
type Field interface {
	BSONName() string
}

// Builds a map of field names to Field objects for the given struct type T.
//
// > NOTE: This function is internal only.
func buildFields[T any]() map[string]Field {
	t := reflect.TypeOf((*T)(nil)).Elem()

	fields := make(map[string]Field)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Skip embedded structs (handle separately if needed)
		if f.Anonymous {
			continue
		}

		bsonTag := f.Tag.Get("bson")
		name := parseBSONName(bsonTag, f.Name)

		fields[f.Name] = NewFieldFromType(f.Type, name)
	}

	return fields
}

// Parses the BSON tag to extract the field name, falling back to the struct field
// name if the tag is empty.
//
// > NOTE: This function is internal only.
func parseBSONName(tag, fallback string) string {
	if tag == "" {
		return fallback
	}

	name := strings.Split(tag, ",")[0]
	if name == "" {
		return fallback
	}

	return name
}

// Creates a Field object based on the Go type. This is used to build the schema
// information for the MongORM instance.
//
// > NOTE: This function is internal only.
func NewFieldFromType(t reflect.Type, name string) Field {
	switch t.Kind() {
	case reflect.Pointer:
		return NewFieldFromType(t.Elem(), name)

	case reflect.String:
		return primitives.StringType(name)

	case reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int:
		return primitives.Int64Type(name)

	case reflect.Float64, reflect.Float32:
		return primitives.Float64Type(name)

	case reflect.Bool:
		return primitives.BoolType(name)

	case reflect.Struct:
		if t == reflect.TypeOf(bson.ObjectID{}) {
			return primitives.ObjectIDType(name)
		}
	}

	switch t.Name() {
	case "ObjectID":
		return primitives.ObjectIDType(name)

	case "Time":
		return primitives.TimestampType(name)

	case "GeoPoint", "GeoLineString", "GeoPolygon":
		return primitives.GeoType(name)
	}

	return primitives.GenericType(name)
}

// Generates a struct of Field objects corresponding to the fields of the struct type T.
// This is used to create a schema struct that can be used for type-safe queries and updates.
//
// Example usage:
//
//		type ToDo struct {
//		  Text *string `bson:"text"`
//		}
//		type ToDoSchema struct {
//		  Text *primitives.StringField
//		}
//	 var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()
func FieldsOf[T any, F any]() F {
	var out F

	t := reflect.TypeOf((*T)(nil)).Elem()
	v := reflect.ValueOf(&out).Elem()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		// Handle embedded structs (inline)
		if sf.Anonymous {
			continue
		}

		bsonName := parseBSONName(sf.Tag.Get("bson"), sf.Name)

		// Find matching field in F by name
		fField := v.FieldByName(sf.Name)
		if !fField.IsValid() || !fField.CanSet() {
			continue
		}

		fieldObj := NewFieldFromType(sf.Type, bsonName)

		fField.Set(reflect.ValueOf(fieldObj))
	}

	return out
}

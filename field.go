package mongorm

import (
	"reflect"
	"strings"

	"github.com/azayn-labs/mongorm/primitives"
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

	case reflect.Slice, reflect.Array:
		elementType := dereferenceType(t.Elem())
		if elementType.Kind() == reflect.String {
			return primitives.StringArrayType(name)
		}

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
	return fieldsOfWithPrefix[T, F]("")
}

// Generates nested fields for a struct type T under a parent field path.
// This is useful for user-defined struct fields mapped as GenericField where
// deep, typed primitives are still desired.
//
// Example:
//
//	type Profile struct {
//	  Provider *string `bson:"provider,omitempty"`
//	}
//	type ProfileSchema struct {
//	  Provider *primitives.StringField
//	}
//	var UserFields = mongorm.FieldsOf[User, UserSchema]()
//	var ProfileFields = mongorm.NestedFieldsOf[Profile, ProfileSchema](UserFields.Goth)
func NestedFieldsOf[T any, F any](parent Field) F {
	if parent == nil {
		var out F
		return out
	}

	return fieldsOfWithPrefix[T, F](parent.BSONName())
}

func fieldsOfWithPrefix[T any, F any](prefix string) F {
	var out F

	modelType := reflect.TypeOf((*T)(nil)).Elem()
	schemaValue := reflect.ValueOf(&out).Elem()
	populateSchemaFields(modelType, schemaValue, prefix)

	return out
}

func populateSchemaFields(modelType reflect.Type, schemaValue reflect.Value, prefix string) {
	modelType = dereferenceType(modelType)
	if modelType.Kind() != reflect.Struct {
		return
	}

	schemaValue = dereferenceValue(schemaValue)
	if !schemaValue.IsValid() || schemaValue.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < modelType.NumField(); i++ {
		sf := modelType.Field(i)
		if sf.Anonymous {
			continue
		}

		fField := schemaValue.FieldByName(sf.Name)
		if !fField.IsValid() || !fField.CanSet() {
			continue
		}

		bsonName := parseBSONName(sf.Tag.Get("bson"), sf.Name)
		if prefix != "" {
			bsonName = prefix + "." + bsonName
		}

		if populateNestedSchemaField(sf.Type, fField, bsonName) {
			continue
		}

		fieldObj := NewFieldFromType(sf.Type, bsonName)
		if setSchemaField(fField, fieldObj) {
			continue
		}

		schemaFieldObj, ok := NewFieldFromSchemaType(fField.Type(), bsonName)
		if ok {
			setSchemaField(fField, schemaFieldObj)
		}
	}
}

// Creates a Field object based on a schema field type. This lets schema definitions
// explicitly select a primitive wrapper, even when the model field Go type differs.
func NewFieldFromSchemaType(t reflect.Type, name string) (Field, bool) {
	switch t {
	case reflect.TypeOf((*primitives.StringField)(nil)):
		return primitives.StringType(name), true
	case reflect.TypeOf((*primitives.Int64Field)(nil)):
		return primitives.Int64Type(name), true
	case reflect.TypeOf((*primitives.Float64Field)(nil)):
		return primitives.Float64Type(name), true
	case reflect.TypeOf((*primitives.BoolField)(nil)):
		return primitives.BoolType(name), true
	case reflect.TypeOf((*primitives.ObjectIDField)(nil)):
		return primitives.ObjectIDType(name), true
	case reflect.TypeOf((*primitives.TimestampField)(nil)):
		return primitives.TimestampType(name), true
	case reflect.TypeOf((*primitives.GeoField)(nil)):
		return primitives.GeoType(name), true
	case reflect.TypeOf((*primitives.GenericField)(nil)):
		return primitives.GenericType(name), true
	case reflect.TypeOf((*primitives.StringArrayField)(nil)):
		return primitives.StringArrayType(name), true
	}

	return nil, false
}

func populateNestedSchemaField(modelFieldType reflect.Type, schemaField reflect.Value, prefix string) bool {
	modelStruct, ok := nestedModelStructType(modelFieldType)
	if !ok {
		return false
	}

	if isNativePrimitiveStruct(modelStruct) || isFieldType(schemaField.Type()) {
		return false
	}

	targetType := schemaField.Type()
	if targetType.Kind() == reflect.Pointer {
		if targetType.Elem().Kind() != reflect.Struct {
			return false
		}

		if schemaField.IsNil() {
			schemaField.Set(reflect.New(targetType.Elem()))
		}

		populateSchemaFields(modelStruct, schemaField.Elem(), prefix)
		return true
	}

	if targetType.Kind() != reflect.Struct {
		return false
	}

	populateSchemaFields(modelStruct, schemaField, prefix)
	return true
}

func nestedModelStructType(t reflect.Type) (reflect.Type, bool) {
	t = dereferenceType(t)

	if t.Kind() == reflect.Struct {
		return t, true
	}

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		elementType := dereferenceType(t.Elem())
		if elementType.Kind() == reflect.Struct {
			return elementType, true
		}
	}

	return nil, false
}

func isNativePrimitiveStruct(t reflect.Type) bool {
	if t == reflect.TypeOf(bson.ObjectID{}) {
		return true
	}

	switch t.Name() {
	case "Time", "GeoPoint", "GeoLineString", "GeoPolygon":
		return true
	}

	return false
}

func isFieldType(t reflect.Type) bool {
	fieldType := reflect.TypeOf((*Field)(nil)).Elem()

	if t.Implements(fieldType) {
		return true
	}

	if t.Kind() == reflect.Pointer {
		return t.Implements(fieldType)
	}

	if t.Kind() == reflect.Struct {
		return reflect.PointerTo(t).Implements(fieldType)
	}

	return false
}

func dereferenceType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t
}

func dereferenceValue(v reflect.Value) reflect.Value {
	for v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}

	return v
}

// Assigns a generated Field into schema field slots while safely supporting:
// - exact concrete primitive pointers (e.g. *primitives.StringField)
// - interface fields (e.g. mongorm.Field, any)
// - pointer-to-interface fields (e.g. *any)
//
// > NOTE: This function is internal only.
func setSchemaField(target reflect.Value, fieldObj Field) bool {
	if !target.IsValid() || !target.CanSet() || fieldObj == nil {
		return false
	}

	value := reflect.ValueOf(fieldObj)
	targetType := target.Type()

	if value.Type().AssignableTo(targetType) {
		target.Set(value)
		return true
	}

	if value.Type().ConvertibleTo(targetType) {
		target.Set(value.Convert(targetType))
		return true
	}

	if target.Kind() == reflect.Interface {
		if targetType.NumMethod() == 0 || value.Type().Implements(targetType) {
			target.Set(value)
			return true
		}
		return false
	}

	if target.Kind() == reflect.Pointer && targetType.Elem().Kind() == reflect.Interface {
		ifaceType := targetType.Elem()
		if ifaceType.NumMethod() == 0 || value.Type().Implements(ifaceType) {
			ptr := reflect.New(ifaceType)
			ptr.Elem().Set(value)
			target.Set(ptr)
			return true
		}
	}

	return false
}

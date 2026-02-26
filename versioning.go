package mongorm

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongORM[T]) getVersionField() (reflect.Value, string, bool, error) {
	if m == nil || m.schema == nil {
		return reflect.Value{}, "", false, configErrorf("schema is nil")
	}

	v := reflect.ValueOf(m.schema)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, "", false, configErrorf("schema must be a pointer to struct")
	}

	structValue := v.Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		fieldType := structType.Field(i)
		if !doesModelIncludeAnyModelFlags(fieldType.Tag, string(ModelTagVersion)) {
			continue
		}

		if fieldType.PkgPath != "" {
			return reflect.Value{}, "", false, configErrorf("version field must be exported")
		}

		bsonName := strings.Split(fieldType.Tag.Get("bson"), ",")[0]
		if bsonName == "" {
			bsonName = fieldType.Name
		}

		sv := structValue.Field(i)
		if !sv.CanSet() || !sv.CanAddr() {
			return reflect.Value{}, "", false, configErrorf("version field must be settable and addressable")
		}

		return sv, bsonName, true, nil
	}

	return reflect.Value{}, "", false, nil
}

func readVersionValue(field reflect.Value) (int64, bool, error) {
	if !field.IsValid() {
		return 0, false, nil
	}

	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return 0, false, nil
		}

		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int(), true, nil
	default:
		return 0, false, configErrorf("version field must be an integer type")
	}
}

func setVersionValue(field reflect.Value, version int64) error {
	if !field.IsValid() {
		return nil
	}

	if !field.CanAddr() {
		return configErrorf("version field must be addressable")
	}

	if field.Kind() == reflect.Pointer {
		if !field.CanSet() {
			return configErrorf("version field must be settable")
		}

		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	if !field.CanSet() || !field.CanAddr() {
		return configErrorf("version field must be settable and addressable")
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.OverflowInt(version) {
			return configErrorf("version value overflows field type")
		}

		field.SetInt(version)
		return nil
	default:
		return configErrorf("version field must be an integer type")
	}
}

func (m *MongORM[T]) initializeVersionForInsert() error {
	field, _, exists, err := m.getVersionField()
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	current, hasValue, err := readVersionValue(field)
	if err != nil {
		return err
	}

	if !hasValue || current <= 0 {
		return setVersionValue(field, 1)
	}

	return nil
}

func (m *MongORM[T]) applyOptimisticLock(filter *bson.M, update *bson.M) (bool, error) {
	field, versionKey, exists, err := m.getVersionField()
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	current, hasValue, err := readVersionValue(field)
	if err != nil {
		return true, err
	}
	if !hasValue || current <= 0 {
		return false, nil
	}

	(*filter)[versionKey] = current

	if *update == nil {
		*update = bson.M{}
	}

	incDoc, ok := (*update)["$inc"].(bson.M)
	if !ok || incDoc == nil {
		incDoc = bson.M{}
	}
	incDoc[versionKey] = int64(1)
	(*update)["$inc"] = incDoc

	return true, nil
}

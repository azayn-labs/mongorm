package mongorm

import (
	"maps"
	"reflect"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongORM[T]) IsModified(field any) bool {
	if m == nil || len(m.modified) == 0 {
		return false
	}

	normalized := strings.TrimSpace(m.resolveFieldBSONName(field))
	if normalized == "" {
		return false
	}

	if _, ok := m.modified[normalized]; ok {
		return true
	}

	for changedPath := range m.modified {
		if strings.HasPrefix(changedPath, normalized+".") || strings.HasPrefix(normalized, changedPath+".") {
			return true
		}
	}

	return false
}

func (m *MongORM[T]) ModifiedFields() []Field {
	if m == nil || len(m.modified) == 0 {
		return []Field{}
	}

	fieldNames := slices.Collect(maps.Keys(m.modified))
	slices.Sort(fieldNames)

	fields := make([]Field, 0, len(fieldNames))
	for _, fieldName := range fieldNames {
		field := RawField(fieldName)
		if field != nil {
			fields = append(fields, field)
		}
	}

	return fields
}

func (m *MongORM[T]) ModifiedValue(field any) (oldValue any, newValue any, ok bool) {
	if m == nil {
		return nil, nil, false
	}

	path := strings.TrimSpace(m.resolveFieldBSONName(field))
	if path == "" {
		return nil, nil, false
	}

	if _, changed := m.modified[path]; !changed {
		return nil, nil, false
	}

	oldValue, _ = m.schemaValueByPath(path)
	newValue, hasNew := m.modifiedNewValue(path, oldValue)

	if !hasNew {
		newValue = oldValue
	}

	return unwrapPointers(oldValue), unwrapPointers(newValue), true
}

func (m *MongORM[T]) clearModified() {
	if m == nil {
		return
	}

	m.modified = map[string]struct{}{}
}

func (m *MongORM[T]) markModified(fieldPath string) {
	if m == nil {
		return
	}

	normalized := strings.TrimSpace(fieldPath)
	if normalized == "" {
		return
	}

	if m.modified == nil {
		m.modified = map[string]struct{}{}
	}

	m.modified[normalized] = struct{}{}
}

func (m *MongORM[T]) rebuildModifiedFromUpdate(update bson.M) {
	m.clearModified()
	m.appendModifiedFromUpdate(update)
}

func (m *MongORM[T]) appendModifiedFromUpdate(update bson.M) {
	if update == nil {
		return
	}

	for key, value := range update {
		if strings.HasPrefix(key, "$") {
			for _, fieldPath := range extractFieldPaths(value) {
				m.markModified(fieldPath)
			}
			continue
		}

		m.markModified(key)
	}
}

func extractFieldPaths(value any) []string {
	result := []string{}

	switch typed := value.(type) {
	case bson.M:
		for key := range typed {
			result = append(result, key)
		}
	case map[string]any:
		for key := range typed {
			result = append(result, key)
		}
	case bson.D:
		for _, entry := range typed {
			if entry.Key != "" {
				result = append(result, entry.Key)
			}
		}
	}

	return result
}

func (m *MongORM[T]) rebuildModifiedFromSchema() {
	m.clearModified()

	if m == nil || m.schema == nil {
		return
	}

	collectModifiedSchemaFields(reflect.ValueOf(m.schema), "", m.markModified)
}

func (m *MongORM[T]) schemaValueByPath(path string) (any, bool) {
	if m == nil || m.schema == nil {
		return nil, false
	}

	parts := strings.Split(path, ".")
	current := reflect.ValueOf(m.schema)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, false
		}

		if part == "$" || part == "$[]" || strings.HasPrefix(part, "$[") {
			return nil, false
		}

		current = dereferenceValue(current)
		if !current.IsValid() {
			return nil, false
		}

		switch current.Kind() {
		case reflect.Struct:
			fieldValue, found := findStructFieldByBSONName(current, part)
			if !found {
				return nil, false
			}
			current = fieldValue
		case reflect.Map:
			if current.Type().Key().Kind() != reflect.String {
				return nil, false
			}
			entry := current.MapIndex(reflect.ValueOf(part))
			if !entry.IsValid() {
				return nil, false
			}
			current = entry
		default:
			return nil, false
		}
	}

	current = dereferenceValue(current)
	if !current.IsValid() {
		return nil, false
	}

	return current.Interface(), true
}

func findStructFieldByBSONName(value reflect.Value, bsonName string) (reflect.Value, bool) {
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		if fieldType.PkgPath != "" {
			continue
		}

		name := parseBSONName(fieldType.Tag.Get("bson"), fieldType.Name)
		if name != bsonName {
			continue
		}

		return value.Field(i), true
	}

	return reflect.Value{}, false
}

func (m *MongORM[T]) modifiedNewValue(path string, oldValue any) (any, bool) {
	if m == nil || m.operations == nil || m.operations.update == nil {
		return nil, false
	}

	update := m.operations.update

	if set, ok := update["$set"].(bson.M); ok {
		if value, exists := set[path]; exists {
			return value, true
		}
	}

	if unset, ok := update["$unset"].(bson.M); ok {
		if _, exists := unset[path]; exists {
			return nil, true
		}
	}

	if inc, ok := update["$inc"].(bson.M); ok {
		if delta, exists := inc[path]; exists {
			if value, computed := computeIncrementValue(oldValue, delta); computed {
				return value, true
			}
			return delta, true
		}
	}

	for _, op := range []string{"$push", "$addToSet", "$pull", "$pop"} {
		doc, ok := update[op].(bson.M)
		if !ok {
			continue
		}

		if value, exists := doc[path]; exists {
			return value, true
		}
	}

	if value, exists := update[path]; exists {
		return value, true
	}

	return nil, false
}

func computeIncrementValue(oldValue any, delta any) (any, bool) {
	oldFloat, oldOk := asFloat64(oldValue)
	deltaFloat, deltaOk := asFloat64(delta)
	if !oldOk || !deltaOk {
		return nil, false
	}

	return oldFloat + deltaFloat, true
}

func asFloat64(value any) (float64, bool) {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return 0, false
	}

	v = dereferenceValue(v)
	if !v.IsValid() {
		return 0, false
	}

	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	default:
		return 0, false
	}
}

func unwrapPointers(value any) any {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return nil
	}

	v = dereferenceValue(v)
	if !v.IsValid() {
		return nil
	}

	return v.Interface()
}

func collectModifiedSchemaFields(value reflect.Value, prefix string, mark func(string)) {
	if !value.IsValid() {
		return
	}

	value = dereferenceValue(value)
	if !value.IsValid() {
		return
	}

	valueType := value.Type()
	if valueType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < valueType.NumField(); i++ {
		fieldType := valueType.Field(i)
		if fieldType.PkgPath != "" {
			continue
		}

		fieldValue := value.Field(i)
		bsonName := parseBSONName(fieldType.Tag.Get("bson"), fieldType.Name)
		fieldPath := bsonName
		if prefix != "" {
			fieldPath = prefix + "." + bsonName
		}

		fieldKind := fieldValue.Kind()
		if fieldKind == reflect.Pointer {
			if fieldValue.IsNil() {
				continue
			}

			elemType := dereferenceType(fieldValue.Type())
			if elemType.Kind() == reflect.Struct && !isNativePrimitiveStruct(elemType) {
				collectModifiedSchemaFields(fieldValue.Elem(), fieldPath, mark)
				continue
			}

			mark(fieldPath)
			continue
		}

		if fieldKind == reflect.Struct && !isNativePrimitiveStruct(fieldValue.Type()) {
			if fieldValue.IsZero() {
				continue
			}
			collectModifiedSchemaFields(fieldValue, fieldPath, mark)
			continue
		}

		if fieldValue.IsZero() {
			continue
		}

		mark(fieldPath)
	}
}

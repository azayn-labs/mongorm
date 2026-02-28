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

	normalized := strings.TrimSpace(resolveFieldBSONName(field))
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

package mongorm

import (
	"strconv"
	"strings"
)

type dynamicField struct {
	name string
}

func (f *dynamicField) BSONName() string {
	return f.name
}

// RawField builds a Field from a raw BSON path.
func RawField(path string) Field {
	normalized := strings.Trim(strings.TrimSpace(path), ".")
	if normalized == "" {
		return nil
	}

	return &dynamicField{name: normalized}
}

// FieldPath appends a dotted suffix path to an existing field path.
func FieldPath(base Field, suffix string) Field {
	suffixNormalized := strings.Trim(strings.TrimSpace(suffix), ".")

	if base == nil {
		return RawField(suffixNormalized)
	}

	baseName := strings.Trim(strings.TrimSpace(base.BSONName()), ".")
	if baseName == "" {
		return RawField(suffixNormalized)
	}

	if suffixNormalized == "" {
		return RawField(baseName)
	}

	return RawField(baseName + "." + suffixNormalized)
}

// Positional builds a field path using MongoDB positional operator `$`.
// Example: Positional(ToDoFields.Items) => "items.$"
func Positional(field Field) Field {
	return FieldPath(field, "$")
}

// PositionalAll builds a field path using MongoDB all positional operator `$[]`.
// Example: PositionalAll(ToDoFields.Items) => "items.$[]"
func PositionalAll(field Field) Field {
	return FieldPath(field, "$[]")
}

// PositionalFiltered builds a field path using MongoDB filtered positional operator `$[identifier]`.
// Example: PositionalFiltered(ToDoFields.Items, "item") => "items.$[item]"
func PositionalFiltered(field Field, identifier string) Field {
	name := strings.TrimSpace(identifier)
	if name == "" {
		return PositionalAll(field)
	}

	return FieldPath(field, "$["+name+"]")
}

// Indexed builds a field path for a specific array index.
// Example: Indexed(ToDoFields.Items, 2) => "items.2"
func Indexed(field Field, index int) Field {
	if index < 0 {
		return nil
	}

	return FieldPath(field, strconv.Itoa(index))
}

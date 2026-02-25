package mongorm

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type ModelTags string

// Connection tags
const (
	ModelTagDatabase         ModelTags = "connection:database"
	ModelTagCollection       ModelTags = "connection:collection"
	ModelTagConnectionString ModelTags = "connection:url"
)

// Field tags
const (
	ModelTagPrimary            ModelTags = "primary"
	ModelTagReadonly           ModelTags = "readonly"
	ModelTagTimestampCreatedAt ModelTags = "timestamp:created_at"
	ModelTagTimestampUpdatedAt ModelTags = "timestamp:updated_at"
)

func doesModelIncludeAnyModelFlags(tag reflect.StructTag, flags ...string) bool {
	tags := getModelTags(tag)

	for _, flag := range flags {
		if slices.Contains(tags, flag) {
			return true
		}
	}

	return false
}

func getModelTags(tag reflect.StructTag) []string {
	v := tag.Get("mongorm")
	if v == "" {
		return nil
	}

	var tags []string
	for f := range strings.SplitSeq(v, ",") {
		tags = append(tags, strings.TrimSpace(f))
	}

	return tags
}

func (m *MongORM[T]) getFieldByTag(tag ModelTags) (string, string, error) {
	v := reflect.ValueOf(m.schema).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)

		if doesModelIncludeAnyModelFlags(fieldType.Tag, string(tag)) {
			return fieldType.Name, strings.Split(fieldType.Tag.Get("bson"), ",")[0], nil
		}
	}

	return "", "", fmt.Errorf("No field found with tag: %s", tag)
}

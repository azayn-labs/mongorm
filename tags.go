package mongorm

import (
	"reflect"
	"slices"
	"strings"
)

type ModelTags string

// Connection tags
// These tags are used to specify connection information for the MongORM instance, such as
// the database name, collection name, and connection string. These tags can be included
// in the struct fields of the schema to provide the necessary information for connecting
// to the MongoDB database. The MongORM instance will use these tags to configure its connection
// settings when performing operations.
// Example usage:
//
//	type ToDo struct {
//	   Text *string `bson:"text"`
//	   // MongORM options
//	   connectionString *string `mongorm:"[connection_string],connection:url"`
//	}
const (
	ModelTagDatabase         ModelTags = "connection:database"
	ModelTagCollection       ModelTags = "connection:collection"
	ModelTagConnectionString ModelTags = "connection:url"
)

// Field tags
// These tags are used to specify field-level options for the MongORM instance, such as
// primary key fields, read-only fields, and timestamp fields. These tags can be included
// in the struct fields of the schema to provide additional information about how the fields
// should be treated during operations. The MongORM instance will use these tags to determine
// how to handle the fields when performing operations such as updates and saves.
// Example usage:
//
//	type ToDo struct {
//	   Text *string `bson:"text"`
//	   // MongORM Tags
//	   ID *string `bson:"_id" mongorm:"primary"`
//	   CreatedAt *time.Time `bson:"created_at" mongorm:"timestamp:created_at"`
//	   UpdatedAt *time.Time `bson:"updated_at" mongorm:"timestamp:updated_at"`
//	}
const (
	ModelTagPrimary            ModelTags = "primary"
	ModelTagVersion            ModelTags = "version"
	ModelTagReadonly           ModelTags = "readonly"
	ModelTagTimestampCreatedAt ModelTags = "timestamp:created_at"
	ModelTagTimestampUpdatedAt ModelTags = "timestamp:updated_at"
)

// getFieldNameFromTag extracts the field name from the provided struct tag. If the tag is empty
// or does not contain a valid field name, it returns the fallback field name. This function is
// used internally to determine the field name for a struct field based on its tags.
//
// > NOTE: This function is internal only.
func doesModelIncludeAnyModelFlags(tag reflect.StructTag, flags ...string) bool {
	tags := getModelTags(tag)

	for _, flag := range flags {
		if slices.Contains(tags, flag) {
			return true
		}
	}

	return false
}

// getModelTags parses the "mongorm" struct tag and returns a slice of individual tags. It splits
// the tag value by commas and trims any whitespace from each tag. If the "mongorm" tag is not
// present or is empty, it returns nil. This function is used internally to extract the relevant
// tags from a struct field for processing by the MongORM instance.
//
// > NOTE: This function is internal only.
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

// MongORMOptions holds the configuration options for a MongORM instance, including settings for
// timestamps, collection and database names, and the MongoDB client. This struct is used to
// customize the behavior of the MongORM instance when connecting to the database and performing
// operations.
//
// > NOTE: This struct is not intended for public use.
func (m *MongORM[T]) getFieldByTag(tag ModelTags) (string, string, error) {
	v := reflect.ValueOf(m.schema).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)

		if doesModelIncludeAnyModelFlags(fieldType.Tag, string(tag)) {
			return fieldType.Name, strings.Split(fieldType.Tag.Get("bson"), ",")[0], nil
		}
	}

	return "", "", configErrorf("no field found with tag: %s", tag)
}

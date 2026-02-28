package mongorm

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// applyTimestampsToUpdateDoc injects timestamp fields only into the outgoing update
// document right before update execution. It does not mutate schema fields.
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) applyTimestampsToUpdateDoc(update *bson.M) {
	if !m.options.Timestamps {
		return
	}
	if update == nil {
		return
	}

	now := time.Now()

	if _, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt); err == nil {
		if (*update)["$set"] == nil {
			(*update)["$set"] = bson.M{}
		}

		set, ok := (*update)["$set"].(bson.M)
		if !ok || set == nil {
			set = bson.M{}
		}

		set[updatedFieldName] = now
		(*update)["$set"] = set
	}
}

// documentForInsertWithTimestamps builds the outgoing insert document and injects
// timestamp fields right before insert execution. It does not mutate schema fields.
func (m *MongORM[T]) documentForInsertWithTimestamps() (bson.M, error) {
	raw, err := bson.Marshal(m.schema)
	if err != nil {
		return nil, err
	}

	doc := bson.M{}
	if err := bson.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}

	if !m.options.Timestamps {
		return doc, nil
	}

	now := time.Now()

	if _, createdFieldName, err := m.getFieldByTag(ModelTagTimestampCreatedAt); err == nil {
		if !timestampFieldHasValue(doc[createdFieldName]) {
			doc[createdFieldName] = now
		}
	}

	if _, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt); err == nil {
		doc[updatedFieldName] = now
	}

	return doc, nil
}

func timestampFieldHasValue(value any) bool {
	if value == nil {
		return false
	}

	switch typed := value.(type) {
	case time.Time:
		return !typed.IsZero()
	case *time.Time:
		return typed != nil && !typed.IsZero()
	case bson.DateTime:
		return typed != 0
	default:
		return true
	}
}

// setTimestampRequirementsFromSchema checks the schema for any fields that are tagged with
// the timestamp tags (created_at and updated_at). If it finds any fields with these tags,
// it sets the Timestamps option to true, indicating that the MongORM instance should
// manage timestamps for this schema.
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) setTimestampRequirementsFromSchema() error {
	ref := reflect.ValueOf(m.schema).Elem()
	t := ref.Type()

	counter := 0

	for i := 0; i < ref.NumField(); i++ {
		fieldType := t.Field(i)

		// Skip nonexported fields
		if fieldType.PkgPath != "" {
			continue
		}

		if doesModelIncludeAnyModelFlags(
			fieldType.Tag,
			string(ModelTagTimestampCreatedAt),
			string(ModelTagTimestampUpdatedAt),
		) {
			tags := getModelTags(fieldType.Tag)
			if len(tags) <= 1 {
				return configErrorf("field %s is missing the timestamps tag value", fieldType.Name)
			}

			counter++
		}
	}

	if counter > 0 {
		m.options.Timestamps = true
	}

	return nil
}

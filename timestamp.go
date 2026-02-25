package mongorm

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongORM[T]) applyTimestamps() {
	if !m.options.Timestamps {
		return
	}

	v := reflect.ValueOf(m.schema).Elem()
	now := time.Now()

	createdField, _, err := m.getFieldByTag(ModelTagTimestampCreatedAt)
	if err != nil {
		return
	}

	updatedField, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
	if err != nil {
		return
	}

	if f := v.FieldByName(createdField); f.IsValid() && f.CanSet() {
		if tI, ok := f.Interface().(*time.Time); ok {
			if tI == nil || tI.IsZero() {
				f.Set(reflect.ValueOf(&now))
			}
		}
	}

	if f := v.FieldByName(updatedField); f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(&now))
		if m.operations.update["$set"] == nil {
			if m.operations.update == nil {
				m.operations.update = bson.M{}
			}

			m.operations.update["$set"] = bson.M{updatedFieldName: now}
		} else {
			set, ok := m.operations.update["$set"].(bson.M)
			if ok {
				set[updatedFieldName] = now
			}
		}
	}
}

func (m *MongORM[T]) setTimestampRequirementsFromSchema() {
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
				panic(fmt.Sprintf("Field %s is missing the timestamps tag value", fieldType.Name))
			}

			counter++
		}
	}

	if counter > 1 {
		m.options.Timestamps = true
	}
}

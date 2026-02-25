package mongorm

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func mergeMaps(one bson.M, two bson.M) bson.M {
	for k, v := range two {
		one[k] = v
	}

	return one
}

func (m *MongORM[T]) Set(value *T) *MongORM[T] {
	if value == nil {
		return m
	}

	set := bson.M{}
	v := reflect.ValueOf(value).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		if doesModelIncludeAnyModelFlags(
			t.Field(i).Tag,
			string(ModelTagPrimary),
			string(ModelTagReadonly),
		) {
			// These fields cannot be updated
			continue
		}

		if fieldValue.Kind() == reflect.Pointer {
			if !fieldValue.IsNil() {
				fieldName := m.info.fields[fieldType.Name].BSONName()
				set[fieldName] = fieldValue.Interface()
			}
			continue
		}

		if !fieldValue.IsZero() {
			fieldName := m.info.fields[fieldType.Name].BSONName()
			set[fieldName] = fieldValue.Interface()
		}
	}

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err != nil {
			return m
		}
		set[updatedFieldName] = time.Now()
	}

	if len(set) > 0 {
		if m.options.Timestamps {
			_, createdFieldName, err := m.getFieldByTag(ModelTagTimestampCreatedAt)
			if err != nil {
				return m
			}

			if set[createdFieldName] != nil {
				delete(set, createdFieldName) // Never update createdAt
			}
		}

		if m.operations.update["$set"] == nil {
			if m.operations.update == nil {
				m.operations.update = bson.M{}
			}
			m.operations.update["$set"] = set
		} else {
			currentSet, ok := m.operations.update["$set"].(bson.M)
			if ok {
				m.operations.update["$set"] = mergeMaps(currentSet, set)
			}
		}
	}

	return m
}

func (m *MongORM[T]) Unset(value *T) *MongORM[T] {
	if value == nil {
		return m
	}
	_, primarykeyName, err := m.getFieldByTag(ModelTagPrimary)
	if err != nil {
		return m
	}

	unset := bson.M{}
	v := reflect.ValueOf(value).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		if doesModelIncludeAnyModelFlags(
			t.Field(i).Tag,
			string(ModelTagPrimary),
			string(ModelTagReadonly),
		) {
			// These fields cannot be updated
			continue
		}

		if fieldValue.Kind() == reflect.Pointer {
			if !fieldValue.IsNil() {
				unset[m.info.fields[fieldType.Name].BSONName()] = 1
			}
			continue
		}

		if !fieldValue.IsZero() {
			unset[m.info.fields[fieldType.Name].BSONName()] = 1
		}
	}

	if len(unset) > 0 {
		if m.options.Timestamps {
			_, createdFieldName, err := m.getFieldByTag(ModelTagTimestampCreatedAt)
			if err != nil {
				return m
			}
			_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
			if err != nil {
				return m
			}

			if unset[createdFieldName] != nil {
				delete(unset, createdFieldName) // Never unset createdAt
			}
			if unset[updatedFieldName] != nil {
				delete(unset, updatedFieldName) // Never unset updatedAt
			}
		}

		if unset[primarykeyName] != nil {
			delete(unset, primarykeyName) // Never unset primary key
		}

		if m.operations.update["$unset"] == nil {
			if m.operations.update == nil {
				m.operations.update = bson.M{}
			}
			m.operations.update["$unset"] = unset
		} else {
			currentUnset, ok := m.operations.update["$unset"].(bson.M)
			if ok {
				m.operations.update["$unset"] = mergeMaps(currentUnset, unset)
			}
		}

		if m.options.Timestamps {
			_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
			if err != nil {
				return m
			}

			if m.operations.update["$set"] == nil {
				if m.operations.update == nil {
					m.operations.update = bson.M{}
				}
				m.operations.update["$set"] = bson.M{updatedFieldName: time.Now()}
			} else {
				set, ok := m.operations.update["$set"].(bson.M)
				if ok {
					set[updatedFieldName] = time.Now()
				}
			}
		}
	}

	return m
}

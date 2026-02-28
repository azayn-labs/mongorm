package mongorm

import (
	"maps"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func resolveFieldBSONName(field any) string {
	if field == nil {
		return ""
	}

	if typedField, ok := field.(Field); ok && typedField != nil {
		return strings.TrimSpace(typedField.BSONName())
	}

	v := reflect.ValueOf(field)
	if !v.IsValid() {
		return ""
	}

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return ""
	}

	if v.Len() == 0 {
		return ""
	}

	item := v.Index(0)
	if item.IsValid() {
		if item.CanAddr() {
			if typedField, ok := item.Addr().Interface().(Field); ok && typedField != nil {
				return strings.TrimSpace(typedField.BSONName())
			}
		}

		if typedField, ok := item.Interface().(Field); ok && typedField != nil {
			return strings.TrimSpace(typedField.BSONName())
		}
	}

	return ""
}

func (m *MongORM[T]) resolveFieldBSONName(field any) string {
	return resolveFieldBSONName(field)
}

// Set adds the specified fields and values to the update document for the current operation.
// It takes a pointer to a struct of type T, which represents the fields to be updated.
// The method iterates through the fields of the struct, checking for non-zero values and
// adding them to the $set operator in the update document. It also handles timestamp fields
// if the Timestamps option is enabled. The method returns the MongORM instance, allowing for
// method chaining.
//
// Example usage:
//
//	type ToDo struct {
//	   Text *string `bson:"text"`
//	   // MongORM options
//	}
//
//	toDo := &ToDo{Text: mongorm.String("Buy milk")}
//	orm := mongorm.New(&ToDo{})
//	orm.Set(&ToDo{Text: mongorm.String("Canceled Buy milk")})
//	err := orm.Save(ctx)
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
		for fieldName := range set {
			m.markModified(fieldName)
		}

		if m.options.Timestamps {
			_, createdFieldName, err := m.getFieldByTag(ModelTagTimestampCreatedAt)
			if err == nil && set[createdFieldName] != nil {
				delete(set, createdFieldName) // Never update createdAt
				delete(m.modified, createdFieldName)
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
				maps.Copy(currentSet, set)
				m.operations.update["$set"] = currentSet
			}
		}
	}

	return m
}

// SetData adds or overrides a single field/value in the current $set update document.
// It accepts a schema Field, so nested fields are supported via field paths (for example:
// `ToDoFields.User.Email` => `user.email`).
func (m *MongORM[T]) SetData(field any, value any) *MongORM[T] {
	fieldName := m.resolveFieldBSONName(field)
	if fieldName == "" {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagPrimary, ModelTagReadonly) {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagTimestampCreatedAt) {
		return m
	}

	if m.operations.update == nil {
		m.operations.update = bson.M{}
	}

	set, ok := m.operations.update["$set"].(bson.M)
	if !ok || set == nil {
		set = bson.M{}
	}

	set[fieldName] = value
	m.markModified(fieldName)

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err == nil {
			set[updatedFieldName] = time.Now()
			m.markModified(updatedFieldName)
		}
	}

	m.operations.update["$set"] = set

	return m
}

// UnsetData adds or overrides a single field in the current $unset update document.
// It accepts a schema Field, so nested fields are supported via field paths (for example:
// `ToDoFields.User.Email` => `user.email`).
func (m *MongORM[T]) UnsetData(field any) *MongORM[T] {
	fieldName := m.resolveFieldBSONName(field)
	if fieldName == "" {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagPrimary, ModelTagReadonly, ModelTagTimestampCreatedAt, ModelTagTimestampUpdatedAt) {
		return m
	}

	if m.operations.update == nil {
		m.operations.update = bson.M{}
	}

	unset, ok := m.operations.update["$unset"].(bson.M)
	if !ok || unset == nil {
		unset = bson.M{}
	}

	unset[fieldName] = 1
	m.markModified(fieldName)

	if len(unset) > 0 {
		m.operations.update["$unset"] = unset
	} else {
		delete(m.operations.update, "$unset")
	}

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err == nil {
			set, ok := m.operations.update["$set"].(bson.M)
			if !ok || set == nil {
				set = bson.M{}
			}
			set[updatedFieldName] = time.Now()
			m.markModified(updatedFieldName)
			m.operations.update["$set"] = set
		}
	}

	return m
}

// IncData adds or overrides a single field/value in the current $inc update document.
// It accepts a schema Field, so nested fields are supported via field paths.
//
// Example usage:
//
//	orm.Where(ToDoFields.ID.Eq(id)).IncData(ToDoFields.Count, int64(2)).Save(ctx)
//	orm.Where(ToDoFields.ID.Eq(id)).IncData(ToDoFields.User.Score, 1).Save(ctx)
func (m *MongORM[T]) IncData(field any, value any) *MongORM[T] {
	if value == nil {
		return m
	}

	fieldName := m.resolveFieldBSONName(field)
	if fieldName == "" {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagPrimary, ModelTagReadonly, ModelTagTimestampCreatedAt, ModelTagTimestampUpdatedAt) {
		return m
	}

	if m.operations.update == nil {
		m.operations.update = bson.M{}
	}

	inc, ok := m.operations.update["$inc"].(bson.M)
	if !ok || inc == nil {
		inc = bson.M{}
	}

	inc[fieldName] = value
	m.markModified(fieldName)
	m.operations.update["$inc"] = inc

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err == nil {
			set, ok := m.operations.update["$set"].(bson.M)
			if !ok || set == nil {
				set = bson.M{}
			}
			set[updatedFieldName] = time.Now()
			m.markModified(updatedFieldName)
			m.operations.update["$set"] = set
		}
	}

	return m
}

// DecData decrements a field using MongoDB's $inc with a negative delta.
//
// Example usage:
//
//	orm.Where(ToDoFields.ID.Eq(id)).DecData(ToDoFields.Count, 1).Save(ctx)
//	// equivalent raw update: {"$inc": {"count": -1}}
func (m *MongORM[T]) DecData(field any, value int64) *MongORM[T] {
	if value < 0 {
		value = -value
	}

	return m.IncData(field, -value)
}

// PushData adds or overrides a single field/value in the current $push update document.
// It accepts a schema Field, so nested fields are supported via field paths.
func (m *MongORM[T]) PushData(field any, value any) *MongORM[T] {
	if value == nil {
		return m
	}

	fieldName := m.resolveFieldBSONName(field)
	if fieldName == "" {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagPrimary, ModelTagReadonly, ModelTagTimestampCreatedAt, ModelTagTimestampUpdatedAt) {
		return m
	}

	if m.operations.update == nil {
		m.operations.update = bson.M{}
	}

	push, ok := m.operations.update["$push"].(bson.M)
	if !ok || push == nil {
		push = bson.M{}
	}

	push[fieldName] = value
	m.markModified(fieldName)
	m.operations.update["$push"] = push

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err == nil {
			set, ok := m.operations.update["$set"].(bson.M)
			if !ok || set == nil {
				set = bson.M{}
			}
			set[updatedFieldName] = time.Now()
			m.markModified(updatedFieldName)
			m.operations.update["$set"] = set
		}
	}

	return m
}

// PushEachData appends multiple values using MongoDB's $push + $each syntax.
func (m *MongORM[T]) PushEachData(field any, values []any) *MongORM[T] {
	if len(values) == 0 {
		return m
	}

	return m.PushData(field, bson.M{"$each": values})
}

// AddToSetData adds or overrides a single field/value in the current $addToSet update document.
// It accepts a schema Field, so nested fields are supported via field paths.
func (m *MongORM[T]) AddToSetData(field any, value any) *MongORM[T] {
	if value == nil {
		return m
	}

	fieldName := m.resolveFieldBSONName(field)
	if fieldName == "" {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagPrimary, ModelTagReadonly, ModelTagTimestampCreatedAt, ModelTagTimestampUpdatedAt) {
		return m
	}

	if m.operations.update == nil {
		m.operations.update = bson.M{}
	}

	addToSet, ok := m.operations.update["$addToSet"].(bson.M)
	if !ok || addToSet == nil {
		addToSet = bson.M{}
	}

	addToSet[fieldName] = value
	m.markModified(fieldName)
	m.operations.update["$addToSet"] = addToSet

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err == nil {
			set, ok := m.operations.update["$set"].(bson.M)
			if !ok || set == nil {
				set = bson.M{}
			}
			set[updatedFieldName] = time.Now()
			m.markModified(updatedFieldName)
			m.operations.update["$set"] = set
		}
	}

	return m
}

// AddToSetEachData appends multiple values uniquely using MongoDB's $addToSet + $each syntax.
func (m *MongORM[T]) AddToSetEachData(field any, values []any) *MongORM[T] {
	if len(values) == 0 {
		return m
	}

	return m.AddToSetData(field, bson.M{"$each": values})
}

// PullData adds or overrides a single field/value in the current $pull update document.
// It accepts a schema Field, so nested fields are supported via field paths.
func (m *MongORM[T]) PullData(field any, value any) *MongORM[T] {
	if value == nil {
		return m
	}

	fieldName := m.resolveFieldBSONName(field)
	if fieldName == "" {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagPrimary, ModelTagReadonly, ModelTagTimestampCreatedAt, ModelTagTimestampUpdatedAt) {
		return m
	}

	if m.operations.update == nil {
		m.operations.update = bson.M{}
	}

	pull, ok := m.operations.update["$pull"].(bson.M)
	if !ok || pull == nil {
		pull = bson.M{}
	}

	pull[fieldName] = value
	m.markModified(fieldName)
	m.operations.update["$pull"] = pull

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err == nil {
			set, ok := m.operations.update["$set"].(bson.M)
			if !ok || set == nil {
				set = bson.M{}
			}
			set[updatedFieldName] = time.Now()
			m.markModified(updatedFieldName)
			m.operations.update["$set"] = set
		}
	}

	return m
}

// PopData adds or overrides a single field/value in the current $pop update document.
// value must be either -1 (remove first) or 1 (remove last).
func (m *MongORM[T]) PopData(field any, value int) *MongORM[T] {
	if value != -1 && value != 1 {
		return m
	}

	fieldName := m.resolveFieldBSONName(field)
	if fieldName == "" {
		return m
	}

	if m.pathHasAnyModelTag(fieldName, ModelTagPrimary, ModelTagReadonly, ModelTagTimestampCreatedAt, ModelTagTimestampUpdatedAt) {
		return m
	}

	if m.operations.update == nil {
		m.operations.update = bson.M{}
	}

	pop, ok := m.operations.update["$pop"].(bson.M)
	if !ok || pop == nil {
		pop = bson.M{}
	}

	pop[fieldName] = value
	m.markModified(fieldName)
	m.operations.update["$pop"] = pop

	if m.options.Timestamps {
		_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
		if err == nil {
			set, ok := m.operations.update["$set"].(bson.M)
			if !ok || set == nil {
				set = bson.M{}
			}
			set[updatedFieldName] = time.Now()
			m.markModified(updatedFieldName)
			m.operations.update["$set"] = set
		}
	}

	return m
}

// PopFirstData removes the first element from an array field.
func (m *MongORM[T]) PopFirstData(field any) *MongORM[T] {
	return m.PopData(field, -1)
}

// PopLastData removes the last element from an array field.
func (m *MongORM[T]) PopLastData(field any) *MongORM[T] {
	return m.PopData(field, 1)
}

func (m *MongORM[T]) pathHasAnyModelTag(path string, tags ...ModelTags) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}

	typeOfModel := reflect.TypeFor[T]()

	for i := 0; i < typeOfModel.NumField(); i++ {
		field := typeOfModel.Field(i)

		hasAny := false
		for _, tag := range tags {
			if doesModelIncludeAnyModelFlags(field.Tag, string(tag)) {
				hasAny = true
				break
			}
		}

		if !hasAny {
			continue
		}

		fieldName := parseBSONName(field.Tag.Get("bson"), field.Name)
		if fieldName == path || strings.HasPrefix(path, fieldName+".") {
			return true
		}
	}

	return false
}

// Save performs an upsert operation, updating an existing document if it exists or inserting
// a new one if it does not. The method applies any necessary timestamps and executes any
// defined hooks before and after the save operation. It returns an error if the operation
// fails.
//
// Example usage:
//
//	type ToDo struct {
//	   Text *string `bson:"text"`
//	   // MongORM options
//	}
//
//	toDo := &ToDo{Text: mongorm.String("Buy milk")}
//	orm := mongorm.New(&ToDo{})
//	orm.Unset(&ToDo{Text: nil})
//	err := orm.Save(ctx)
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
		for fieldName := range unset {
			m.markModified(fieldName)
		}

		if m.options.Timestamps {
			_, createdFieldName, createdErr := m.getFieldByTag(ModelTagTimestampCreatedAt)
			if createdErr == nil && unset[createdFieldName] != nil {
				delete(unset, createdFieldName) // Never unset createdAt
				delete(m.modified, createdFieldName)
			}

			_, updatedFieldName, updatedErr := m.getFieldByTag(ModelTagTimestampUpdatedAt)
			if updatedErr == nil && unset[updatedFieldName] != nil {
				delete(unset, updatedFieldName) // Never unset updatedAt
				delete(m.modified, updatedFieldName)
			}
		}

		if unset[primarykeyName] != nil {
			delete(unset, primarykeyName) // Never unset primary key
			delete(m.modified, primarykeyName)
		}

		if m.operations.update["$unset"] == nil {
			if m.operations.update == nil {
				m.operations.update = bson.M{}
			}
			m.operations.update["$unset"] = unset
		} else {
			currentUnset, ok := m.operations.update["$unset"].(bson.M)
			if ok {
				maps.Copy(currentUnset, unset)
				m.operations.update["$unset"] = currentUnset
			}
		}

		if m.options.Timestamps {
			_, updatedFieldName, err := m.getFieldByTag(ModelTagTimestampUpdatedAt)
			if err == nil {
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
	}

	return m
}

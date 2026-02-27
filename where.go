package mongorm

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Where adds a query filter to the MongORM instance. It takes a bson.M expression as an
// argument and appends it to the existing query filters using the $and operator.
// This allows you to chain multiple query filters together using the $and operator.
//
// Example usage:
//
//	orm.Where(bson.M{"age": bson.M{"$gt": 30}}).Where(bson.M{"name": "John"})
//	// OR
//	orm.Where(fielType.Age.Gt(30)).Where(fieldType.Name.Eq("John"))
func (m *MongORM[T]) Where(expr bson.M) *MongORM[T] {
	if m.operations.query == nil {
		m.operations.query = bson.M{}
	}

	if m.operations.query["$and"] == nil {
		m.operations.query["$and"] = bson.A{}
	}

	if and, ok := m.operations.query["$and"].(bson.A); ok {
		m.operations.query["$and"] = append(
			and,
			expr,
		)
	}

	return m
}

// Where adds a query filter for a specific field and value to the MongORM instance.
// It constructs a bson.M expression for the given field and value and appends it to the
// existing query filters using the $and operator. This allows you to chain multiple query
// filters together using the $and operator.
//
// Example usage:
//
//	orm.WhereBy(fieldType.Age, 30).WhereBy(fieldType.Name, "John")
func (m *MongORM[T]) WhereBy(field Field, value any) *MongORM[T] {
	name := field.BSONName()
	if m.operations.query == nil {
		m.operations.query = bson.M{}
	}

	if m.operations.query["$and"] == nil {
		m.operations.query["$and"] = bson.A{}
	}

	if and, ok := m.operations.query["$and"].(bson.A); ok {
		m.operations.query["$and"] = append(
			and,
			bson.M{name: value},
		)
	}

	return m
}

// Sort sets the sort order for find operations. It accepts the same values supported by
// MongoDB options, such as bson.D{{"field", 1}} or bson.M{"field": -1}.
func (m *MongORM[T]) Sort(value any) *MongORM[T] {
	m.operations.sort = value
	return m
}

// SortBy sets sort using a schema field and direction.
// Use 1 for ascending and -1 for descending.
func (m *MongORM[T]) SortBy(field Field, direction int) *MongORM[T] {
	if field == nil {
		return m
	}

	m.operations.sort = bson.D{{Key: field.BSONName(), Value: direction}}
	return m
}

// SortAsc sets ascending sort using a schema field.
func (m *MongORM[T]) SortAsc(field Field) *MongORM[T] {
	return m.SortBy(field, 1)
}

// SortDesc sets descending sort using a schema field.
func (m *MongORM[T]) SortDesc(field Field) *MongORM[T] {
	return m.SortBy(field, -1)
}

// Limit sets the maximum number of documents returned by find operations.
// For First()/Find(), this value is ignored because the operation always returns one document.
func (m *MongORM[T]) Limit(value int64) *MongORM[T] {
	m.operations.limit = &value
	return m
}

// Skip sets the number of documents to skip before returning results for find operations.
func (m *MongORM[T]) Skip(value int64) *MongORM[T] {
	m.operations.skip = &value
	return m
}

// Projection sets the fields returned by find operations.
// Example: bson.M{"text": 1, "count": 1}
func (m *MongORM[T]) Projection(value any) *MongORM[T] {
	m.operations.projection = value
	return m
}

// ProjectionInclude sets projection to include only the given schema fields.
func (m *MongORM[T]) ProjectionInclude(fields ...Field) *MongORM[T] {
	projection := bson.M{}

	for _, field := range fields {
		if field == nil {
			continue
		}
		projection[field.BSONName()] = 1
	}

	if len(projection) == 0 {
		return m
	}

	m.operations.projection = projection
	return m
}

// ProjectionExclude sets projection to exclude the given schema fields.
func (m *MongORM[T]) ProjectionExclude(fields ...Field) *MongORM[T] {
	projection := bson.M{}

	for _, field := range fields {
		if field == nil {
			continue
		}
		projection[field.BSONName()] = 0
	}

	if len(projection) == 0 {
		return m
	}

	m.operations.projection = projection
	return m
}

// After adds a keyset-style pagination filter: field > cursor.
func (m *MongORM[T]) After(field Field, cursor any) *MongORM[T] {
	return m.Where(bson.M{field.BSONName(): bson.M{"$gt": cursor}})
}

// Before adds a keyset-style pagination filter: field < cursor.
func (m *MongORM[T]) Before(field Field, cursor any) *MongORM[T] {
	return m.Where(bson.M{field.BSONName(): bson.M{"$lt": cursor}})
}

// PageSize is an alias for Limit and is intended for pagination use-cases.
func (m *MongORM[T]) PageSize(size int64) *MongORM[T] {
	return m.Limit(size)
}

// PaginateAfter applies keyset pagination in ascending order for the given field.
func (m *MongORM[T]) PaginateAfter(field Field, cursor any, size int64) *MongORM[T] {
	name := field.BSONName()

	return m.
		After(field, cursor).
		Sort(bson.D{{Key: name, Value: 1}}).
		PageSize(size)
}

// PaginateBefore applies keyset pagination in descending order for the given field.
func (m *MongORM[T]) PaginateBefore(field Field, cursor any, size int64) *MongORM[T] {
	name := field.BSONName()

	return m.
		Before(field, cursor).
		Sort(bson.D{{Key: name, Value: -1}}).
		PageSize(size)
}

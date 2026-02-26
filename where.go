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

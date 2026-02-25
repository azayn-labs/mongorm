package mongorm

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

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

func (m *MongORM[T]) where(field Field, value any) *MongORM[T] {
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

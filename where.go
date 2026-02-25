package mongorm

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *MongORM[T]) Where(expr bson.M) *MongORM[T] {
	if m.operations.query["$and"] == nil {
		m.operations.query["$and"] = bson.A{}
	}

	m.operations.query["$and"] = append(
		m.operations.query["$and"].(bson.A),
		expr,
	)

	return m
}

func (m *MongORM[T]) where(field Field, value any) *MongORM[T] {
	name := field.BSONName()
	if m.operations.query["$and"] == nil {
		m.operations.query["$and"] = bson.A{}
	}
	m.operations.query["$and"] = append(m.operations.query["$and"].(bson.A), bson.M{name: value})

	return m
}

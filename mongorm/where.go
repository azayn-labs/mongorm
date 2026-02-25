package mongorm

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *Model[T]) Where(expr bson.M) *Model[T] {
	if m.query["$and"] == nil {
		m.query["$and"] = bson.A{}
	}

	m.query["$and"] = append(
		m.query["$and"].(bson.A),
		expr,
	)

	fmt.Println("Where", m.query)

	return m
}

func (m *Model[T]) where(field Field, value any) *Model[T] {
	name := field.BSONName()
	if m.query["$and"] == nil {
		m.query["$and"] = bson.A{}
	}
	m.query["$and"] = append(m.query["$and"].(bson.A), bson.M{name: value})

	return m
}

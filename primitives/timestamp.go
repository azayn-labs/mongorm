package primitives

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TimestampField struct {
	name string
}

func TimestampType(name string) *TimestampField {
	return &TimestampField{name: name}
}

func (f *TimestampField) BSONName() string {
	return f.name
}

// ########## Query methods ###########

// This method generates a query for equality, e.g., {field: value}
func (f *TimestampField) Eq(v time.Time) bson.M {
	return bson.M{f.name: v}
}

// This method generates a query for inequality, e.g., {field: {$ne: value}}
func (f *TimestampField) Ne(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

// This method generates a query for multiple values, e.g., {field: {$in: [value1, value2, ...]}}
func (f *TimestampField) In(v []time.Time) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

// This method generates a query for values not in the given list, e.g., {field: {$nin: [value1, value2, ...]}}
func (f *TimestampField) Nin(v []time.Time) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

// This method generates a query to check if the field exists, e.g., {field: {$exists: true}}
func (f *TimestampField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

// This method generates a query to check if the field does not exist, e.g., {field: {$exists: false}}
func (f *TimestampField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

// This method generates a query to check if the field is null, e.g., {field: null}
func (f *TimestampField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

// This method generates a query to check if the field is not null, e.g., {field: {$ne: null}}
func (f *TimestampField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

// This method generates a query for greater than, e.g., {field: {$gt: value}}
func (f *TimestampField) Gt(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

// This method generates a query for greater than or equal to, e.g., {field: {$gte: value}}
func (f *TimestampField) Gte(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

// This method generates a query for less than, e.g., {field: {$lt: value}}
func (f *TimestampField) Lt(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

// This method generates a query for less than or equal to, e.g., {field: {$lte: value}}
func (f *TimestampField) Lte(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

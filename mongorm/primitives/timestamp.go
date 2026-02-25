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

func (f *TimestampField) Eq(v time.Time) bson.M {
	return bson.M{f.name: v}
}

func (f *TimestampField) Ne(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *TimestampField) In(v []time.Time) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *TimestampField) Nin(v []time.Time) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *TimestampField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *TimestampField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *TimestampField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *TimestampField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *TimestampField) Gt(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

func (f *TimestampField) Gte(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

func (f *TimestampField) Lt(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

func (f *TimestampField) Lte(v time.Time) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

package primitives

import "go.mongodb.org/mongo-driver/v2/bson"

type Int64Field struct {
	name string
}

func Int64Type(name string) *Int64Field {
	return &Int64Field{name: name}
}

func (f *Int64Field) BSONName() string {
	return f.name
}

func (f *Int64Field) Eq(v int64) bson.M {
	return bson.M{f.name: v}
}

func (f *Int64Field) Ne(v int64) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *Int64Field) In(v []int64) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *Int64Field) Nin(v []int64) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *Int64Field) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *Int64Field) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *Int64Field) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *Int64Field) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *Int64Field) Gt(v int64) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

func (f *Int64Field) Gte(v int64) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

func (f *Int64Field) Lt(v int64) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

func (f *Int64Field) Lte(v int64) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

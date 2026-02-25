package primitives

import "go.mongodb.org/mongo-driver/v2/bson"

type Float64Field struct {
	name string
}

func Float64Type(name string) *Float64Field {
	return &Float64Field{name: name}
}

func (f *Float64Field) BSONName() string {
	return f.name
}

func (f *Float64Field) Eq(v float64) bson.M {
	return bson.M{f.name: v}
}

func (f *Float64Field) Ne(v float64) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *Float64Field) In(v []float64) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *Float64Field) Nin(v []float64) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *Float64Field) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *Float64Field) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *Float64Field) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *Float64Field) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *Float64Field) Gt(v float64) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

func (f *Float64Field) Gte(v float64) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

func (f *Float64Field) Lt(v float64) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

func (f *Float64Field) Lte(v float64) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

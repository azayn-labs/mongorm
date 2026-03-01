package primitives

import "go.mongodb.org/mongo-driver/v2/bson"

type Decimal128Field struct {
	name string
}

func Decimal128Type(name string) *Decimal128Field {
	return &Decimal128Field{name: name}
}

func (f *Decimal128Field) BSONName() string {
	return f.name
}

func (f *Decimal128Field) Eq(v bson.Decimal128) bson.M {
	return bson.M{f.name: v}
}

func (f *Decimal128Field) Ne(v bson.Decimal128) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *Decimal128Field) In(v []bson.Decimal128) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *Decimal128Field) Nin(v []bson.Decimal128) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *Decimal128Field) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *Decimal128Field) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *Decimal128Field) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *Decimal128Field) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *Decimal128Field) Gt(v bson.Decimal128) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

func (f *Decimal128Field) Gte(v bson.Decimal128) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

func (f *Decimal128Field) Lt(v bson.Decimal128) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

func (f *Decimal128Field) Lte(v bson.Decimal128) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

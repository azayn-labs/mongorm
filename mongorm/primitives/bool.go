package primitives

import "go.mongodb.org/mongo-driver/v2/bson"

type BoolField struct {
	name string
}

func BoolType(name string) *BoolField {
	return &BoolField{name: name}
}

func (f *BoolField) BSONName() string {
	return f.name
}

func (f *BoolField) Eq(v bool) bson.M {
	return bson.M{f.name: v}
}

func (f *BoolField) Ne(v bool) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *BoolField) In(v []bool) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *BoolField) Nin(v []bool) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *BoolField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *BoolField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *BoolField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *BoolField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

package primitives

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StringField struct {
	name string
}

func StringType(name string) *StringField {
	return &StringField{name: name}
}

func (f *StringField) BSONName() string {
	return f.name
}

func (f *StringField) Eq(v string) bson.M {
	return bson.M{f.name: v}
}

func (f *StringField) Ne(v string) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *StringField) In(v []string) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *StringField) Nin(v []string) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *StringField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *StringField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *StringField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *StringField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

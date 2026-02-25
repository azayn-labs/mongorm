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

// ########## Query methods ###########

// This method generates a query for equality, e.g., {field: value}
func (f *StringField) Eq(v string) bson.M {
	return bson.M{f.name: v}
}

// This method generates a query for regular expression matching, e.g., {field: {$regex: value}}
func (f *StringField) Reg(v string) bson.M {
	return bson.M{f.name: bson.M{"$regex": v}}
}

// This method generates a query for inequality, e.g., {field: {$ne: value}}
func (f *StringField) Ne(v string) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

// This method generates a query for multiple values, e.g., {field: {$in: [value1, value2, ...]}}
func (f *StringField) In(v []string) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

// This method generates a query for values not in the given list, e.g., {field: {$nin: [value1, value2, ...]}}
func (f *StringField) Nin(v []string) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

// This method generates a query to check if the field exists, e.g., {field: {$exists: true}}
func (f *StringField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

// This method generates a query to check if the field does not exist, e.g., {field: {$exists: false}}
func (f *StringField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

// This method generates a query to check if the field is null, e.g., {field: null}
func (f *StringField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

// This method generates a query to check if the field is not null, e.g., {field: {$ne: null}}
func (f *StringField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

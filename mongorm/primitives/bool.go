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

// ########## Query methods ###########

// This method generates a query for equality, e.g., {field: value}
func (f *BoolField) Eq(v bool) bson.M {
	return bson.M{f.name: v}
}

// This method generates a query for inequality, e.g., {field: {$ne: value}}
func (f *BoolField) Ne(v bool) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

// This method generates a query for multiple values, e.g., {field: {$in: [value1, value2, ...]}}
func (f *BoolField) In(v []bool) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

// This method generates a query for values not in the given list, e.g., {field: {$nin: [value1, value2, ...]}}
func (f *BoolField) Nin(v []bool) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

// This method generates a query to check if the field exists, e.g., {field: {$exists: true}}
func (f *BoolField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

// This method generates a query to check if the field does not exist, e.g., {field: {$exists: false}}
func (f *BoolField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

// This method generates a query to check if the field is null, e.g., {field: null}
func (f *BoolField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

// This method generates a query to check if the field is not null, e.g., {field: {$ne: null}}
func (f *BoolField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

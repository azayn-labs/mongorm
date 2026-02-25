package primitives

import "go.mongodb.org/mongo-driver/v2/bson"

type ObjectIDField struct {
	name string
}

func ObjectIDType(name string) *ObjectIDField {
	return &ObjectIDField{name: name}
}

func (f *ObjectIDField) BSONName() string {
	return f.name
}

// ########## Query methods ###########

// This method generates a query for equality, e.g., {field: value}
func (f *ObjectIDField) Eq(v bson.ObjectID) bson.M {
	return bson.M{f.name: v}
}

// This method generates a query for inequality, e.g., {field: {$ne: value}}
func (f *ObjectIDField) Ne(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

// This method generates a query for multiple values, e.g., {field: {$in: [value1, value2, ...]}}
func (f *ObjectIDField) In(v []bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

// This method generates a query for values not in the given list, e.g., {field: {$nin: [value1, value2, ...]}}
func (f *ObjectIDField) Nin(v []bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

// This method generates a query to check if the field exists, e.g., {field: {$exists: true}}
func (f *ObjectIDField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

// This method generates a query to check if the field does not exist, e.g., {field: {$exists: false}}
func (f *ObjectIDField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

// This method generates a query to check if the field is null, e.g., {field: null}
func (f *ObjectIDField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

// This method generates a query to check if the field is not null, e.g., {field: {$ne: null}}
func (f *ObjectIDField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

// This method generates a query for greater than, e.g., {field: {$gt: value}}
func (f *ObjectIDField) Gt(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

// This method generates a query for greater than or equal, e.g., {field: {$gte: value}}
func (f *ObjectIDField) Gte(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

// This method generates a query for less than, e.g., {field: {$lt: value}}
func (f *ObjectIDField) Lt(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

// This method generates a query for less than or equal, e.g., {field: {$lte: value}}
func (f *ObjectIDField) Lte(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

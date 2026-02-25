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

// ########## Query methods ###########

// This method generates a query for equality, e.g., {field: value}
func (f *Float64Field) Eq(v float64) bson.M {
	return bson.M{f.name: v}
}

// This method generates a query for inequality, e.g., {field: {$ne: value}}
func (f *Float64Field) Ne(v float64) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

// This method generates a query for multiple values, e.g., {field: {$in: [value1, value2, ...]}}
func (f *Float64Field) In(v []float64) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

// This method generates a query for values not in the given list, e.g., {field: {$nin: [value1, value2, ...]}}
func (f *Float64Field) Nin(v []float64) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

// This method generates a query to check if the field is not null, e.g., {field: {$ne: null}}
func (f *Float64Field) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

// This method generates a query to check if the field exists, e.g., {field: {$exists: true}}
func (f *Float64Field) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

// This method generates a query to check if the field does not exist, e.g., {field: {$exists: false}}
func (f *Float64Field) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

// This method generates a query to check if the field is null, e.g., {field: null}
func (f *Float64Field) IsNull() bson.M {
	return bson.M{f.name: nil}
}

// This method generates a query for greater than, e.g., {field: {$gt: value}}
func (f *Float64Field) Gt(v float64) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

// This method generates a query for greater than or equal to, e.g., {field: {$gte: value}}
func (f *Float64Field) Gte(v float64) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

// This method generates a query for less than, e.g., {field: {$lt: value}}
func (f *Float64Field) Lt(v float64) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

// This method generates a query for less than or equal to, e.g., {field: {$lte: value}}
func (f *Float64Field) Lte(v float64) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

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

func (f *ObjectIDField) Eq(v bson.ObjectID) bson.M {
	return bson.M{f.name: v}
}

func (f *ObjectIDField) Ne(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *ObjectIDField) In(v []bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *ObjectIDField) Nin(v []bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *ObjectIDField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *ObjectIDField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *ObjectIDField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *ObjectIDField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *ObjectIDField) Gt(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$gt": v}}
}

func (f *ObjectIDField) Gte(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$gte": v}}
}

func (f *ObjectIDField) Lt(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$lt": v}}
}

func (f *ObjectIDField) Lte(v bson.ObjectID) bson.M {
	return bson.M{f.name: bson.M{"$lte": v}}
}

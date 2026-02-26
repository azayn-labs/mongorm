package primitives

import "go.mongodb.org/mongo-driver/v2/bson"

type StringArrayField struct {
	name string
}

func StringArrayType(name string) *StringArrayField {
	return &StringArrayField{name: name}
}

func (f *StringArrayField) BSONName() string {
	return f.name
}

func (f *StringArrayField) Eq(v []string) bson.M {
	return bson.M{f.name: v}
}

func (f *StringArrayField) Ne(v []string) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *StringArrayField) In(v []string) bson.M {
	values := make([]any, len(v))
	for i, value := range v {
		values[i] = value
	}

	return bson.M{f.name: bson.M{"$in": values}}
}

func (f *StringArrayField) Nin(v []string) bson.M {
	values := make([]any, len(v))
	for i, value := range v {
		values[i] = value
	}

	return bson.M{f.name: bson.M{"$nin": values}}
}

func (f *StringArrayField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *StringArrayField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *StringArrayField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *StringArrayField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *StringArrayField) Contains(v string) bson.M {
	return bson.M{f.name: bson.M{"$in": []string{v}}}
}

func (f *StringArrayField) ContainsAll(v []string) bson.M {
	return bson.M{f.name: bson.M{"$all": v}}
}

func (f *StringArrayField) Size(v int) bson.M {
	return bson.M{f.name: bson.M{"$size": v}}
}

func (f *StringArrayField) ElemMatch(v bson.M) bson.M {
	return bson.M{f.name: bson.M{"$elemMatch": v}}
}

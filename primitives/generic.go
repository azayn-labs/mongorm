package primitives

import (
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GenericField struct {
	name string
}

func GenericType(name string) *GenericField {
	return &GenericField{name: name}
}

func (f *GenericField) BSONName() string {
	return f.name
}

func (f *GenericField) Eq(v any) bson.M {
	return bson.M{f.name: v}
}

func (f *GenericField) Ne(v any) bson.M {
	return bson.M{f.name: bson.M{"$ne": v}}
}

func (f *GenericField) In(v []any) bson.M {
	return bson.M{f.name: bson.M{"$in": v}}
}

func (f *GenericField) Nin(v []any) bson.M {
	return bson.M{f.name: bson.M{"$nin": v}}
}

func (f *GenericField) Exists() bson.M {
	return bson.M{f.name: bson.M{"$exists": true}}
}

func (f *GenericField) NotExists() bson.M {
	return bson.M{f.name: bson.M{"$exists": false}}
}

func (f *GenericField) IsNull() bson.M {
	return bson.M{f.name: nil}
}

func (f *GenericField) IsNotNull() bson.M {
	return bson.M{f.name: bson.M{"$ne": nil}}
}

func (f *GenericField) Contains(v any) bson.M {
	return bson.M{f.name: bson.M{"$in": []any{v}}}
}

func (f *GenericField) ContainsAll(v []any) bson.M {
	return bson.M{f.name: bson.M{"$all": v}}
}

func (f *GenericField) Size(v int) bson.M {
	return bson.M{f.name: bson.M{"$size": v}}
}

func (f *GenericField) ElemMatch(v bson.M) bson.M {
	return bson.M{f.name: bson.M{"$elemMatch": v}}
}

func (f *GenericField) Path(path string) *GenericField {
	cleanedPath := strings.TrimSpace(path)
	if cleanedPath == "" {
		return GenericType(f.name)
	}

	return GenericType(f.name + "." + cleanedPath)
}

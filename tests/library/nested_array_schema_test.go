package main

import (
	"reflect"
	"testing"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type userDefinedStruct struct {
	Name  *string `bson:"name,omitempty"`
	Score *int64  `bson:"score,omitempty"`
}

type userDefinedStructSchema struct {
	Name  *primitives.StringField
	Score *primitives.Int64Field
}

type containerWithStructArray struct {
	Key *[]userDefinedStruct `bson:"key,omitempty"`
}

type containerWithStructArraySchema struct {
	Key *userDefinedStructSchema
}

func TestNestedSchemaForPointerSliceStruct(t *testing.T) {
	fields := mongorm.FieldsOf[containerWithStructArray, containerWithStructArraySchema]()

	if fields.Key == nil {
		t.Fatal("expected nested schema for Key to be initialized")
	}

	if fields.Key.Name == nil || fields.Key.Name.BSONName() != "key.name" {
		t.Fatal("expected key.name nested field")
	}

	if fields.Key.Score == nil || fields.Key.Score.BSONName() != "key.score" {
		t.Fatal("expected key.score nested field")
	}
}

func TestNestedSchemaQueryForPointerSliceStruct(t *testing.T) {
	fields := mongorm.FieldsOf[containerWithStructArray, containerWithStructArraySchema]()

	nameQuery := fields.Key.Name.Eq("alice")
	if !reflect.DeepEqual(nameQuery, bson.M{"key.name": "alice"}) {
		t.Fatalf("unexpected key.name query: %#v", nameQuery)
	}

	scoreQuery := fields.Key.Score.Gte(10)
	if !reflect.DeepEqual(scoreQuery, bson.M{"key.score": bson.M{"$gte": int64(10)}}) {
		t.Fatalf("unexpected key.score query: %#v", scoreQuery)
	}
}

package main

import (
	"reflect"
	"testing"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestStrictFilterAndUpdateBuildersFromFields(t *testing.T) {
	filter := mongorm.FilterBy(ToDoFields.User.Email, "john@example.com")
	if !reflect.DeepEqual(filter, bson.M{"user.email": "john@example.com"}) {
		t.Fatalf("unexpected filter: %#v", filter)
	}

	update := mongorm.SetUpdateFromPairs(
		mongorm.FieldValuePair{Field: ToDoFields.Text, Value: "x@example.com"},
		mongorm.FieldValuePair{Field: ToDoFields.User.Auth.Provider, Value: "google"},
	)
	if !reflect.DeepEqual(update, bson.M{"$set": bson.M{"text": "x@example.com", "user.auth.provider": "google"}}) {
		t.Fatalf("unexpected set update: %#v", update)
	}

	unset := mongorm.UnsetUpdateFromFields(ToDoFields.User.Auth.Provider)
	if !reflect.DeepEqual(unset, bson.M{"$unset": bson.M{"user.auth.provider": 1}}) {
		t.Fatalf("unexpected unset update: %#v", unset)
	}

	inc := mongorm.IncUpdateFromPairs(
		mongorm.FieldValuePair{Field: ToDoFields.Count, Value: int64(2)},
		mongorm.FieldValuePair{Field: ToDoFields.User.Auth.Provider, Value: 1},
	)
	if !reflect.DeepEqual(inc, bson.M{"$inc": bson.M{"count": int64(2), "user.auth.provider": 1}}) {
		t.Fatalf("unexpected inc update: %#v", inc)
	}

	push := mongorm.PushUpdateFromPairs(
		mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: "urgent"},
	)
	if !reflect.DeepEqual(push, bson.M{"$push": bson.M{"tags": "urgent"}}) {
		t.Fatalf("unexpected push update: %#v", push)
	}

	addToSet := mongorm.AddToSetUpdateFromPairs(
		mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: bson.M{"$each": []any{"urgent", "backend"}}},
	)
	if !reflect.DeepEqual(addToSet, bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": []any{"urgent", "backend"}}}}) {
		t.Fatalf("unexpected addToSet update: %#v", addToSet)
	}

	pull := mongorm.PullUpdateFromPairs(
		mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: "deprecated"},
	)
	if !reflect.DeepEqual(pull, bson.M{"$pull": bson.M{"tags": "deprecated"}}) {
		t.Fatalf("unexpected pull update: %#v", pull)
	}

	pop := mongorm.PopUpdateFromPairs(
		mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: 1},
	)
	if !reflect.DeepEqual(pop, bson.M{"$pop": bson.M{"tags": 1}}) {
		t.Fatalf("unexpected pop update: %#v", pop)
	}
}

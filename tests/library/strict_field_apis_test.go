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
}

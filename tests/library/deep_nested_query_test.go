package main

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestDeepNestedSchemaMapping(t *testing.T) {
	if ToDoFields.User == nil {
		t.Fatal("expected ToDoFields.User to be initialized")
	}

	if ToDoFields.User.Auth == nil {
		t.Fatal("expected ToDoFields.User.Auth to be initialized")
	}

	if ToDoFields.User.Auth.Provider == nil || ToDoFields.User.Auth.Provider.BSONName() != "user.auth.provider" {
		t.Fatal("expected user.auth.provider bson name")
	}

	if ToDoFields.User.Auth.Scopes == nil || ToDoFields.User.Auth.Scopes.BSONName() != "user.auth.scopes" {
		t.Fatal("expected user.auth.scopes bson name")
	}
}

func TestDeepNestedQueryBuilders(t *testing.T) {
	providerQuery := ToDoFields.User.Auth.Provider.Eq("google")
	if !reflect.DeepEqual(providerQuery, bson.M{"user.auth.provider": "google"}) {
		t.Fatalf("unexpected provider query: %#v", providerQuery)
	}

	scopeQuery := ToDoFields.User.Auth.Scopes.Contains("email")
	if !reflect.DeepEqual(scopeQuery, bson.M{"user.auth.scopes": bson.M{"$in": []string{"email"}}}) {
		t.Fatalf("unexpected scope query: %#v", scopeQuery)
	}
}

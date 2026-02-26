package main

import (
	"reflect"
	"testing"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type userWithPointerScopes struct {
	Scopes *[]string `bson:"scopes,omitempty"`
}

type userWithPointerScopesSchema struct {
	Scopes *primitives.StringArrayField
}

func TestStringArrayFieldQueries(t *testing.T) {
	if !reflect.DeepEqual(ToDoFields.User.Auth.Scopes.Contains("email"), bson.M{"user.auth.scopes": bson.M{"$in": []string{"email"}}}) {
		t.Fatal("unexpected Contains query for StringArrayField")
	}

	if !reflect.DeepEqual(ToDoFields.User.Auth.Scopes.ContainsAll([]string{"email", "profile"}), bson.M{"user.auth.scopes": bson.M{"$all": []string{"email", "profile"}}}) {
		t.Fatal("unexpected ContainsAll query for StringArrayField")
	}

	if !reflect.DeepEqual(ToDoFields.User.Auth.Scopes.Size(2), bson.M{"user.auth.scopes": bson.M{"$size": 2}}) {
		t.Fatal("unexpected Size query for StringArrayField")
	}
}

func TestStringArrayFieldWithPointerSliceModel(t *testing.T) {
	fields := mongorm.FieldsOf[userWithPointerScopes, userWithPointerScopesSchema]()

	if fields.Scopes == nil {
		t.Fatal("expected Scopes field to be initialized")
	}

	if fields.Scopes.BSONName() != "scopes" {
		t.Fatalf("expected scopes BSON name, got %s", fields.Scopes.BSONName())
	}

	if !reflect.DeepEqual(fields.Scopes.Contains("admin"), bson.M{"scopes": bson.M{"$in": []string{"admin"}}}) {
		t.Fatal("unexpected query for pointer-slice string field")
	}
}

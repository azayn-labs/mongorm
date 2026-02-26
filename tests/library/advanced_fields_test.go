package main

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidateAdvancedFieldsOf(t *testing.T) {
	logger(t, "Validating advanced FieldsOf schema")

	if ToDoFields.Meta == nil || ToDoFields.Meta.BSONName() != "meta" {
		t.Fatal("expected Meta field with bson meta")
	}

	if ToDoFields.User == nil {
		t.Fatal("expected User nested schema to be initialized")
	}

	if ToDoFields.User.ID == nil || ToDoFields.User.ID.BSONName() != "user._id" {
		t.Fatal("expected User.ID nested field with bson user._id")
	}

	if ToDoFields.User.Email == nil || ToDoFields.User.Email.BSONName() != "user.email" {
		t.Fatal("expected User.Email nested field with bson user.email")
	}

	if ToDoFields.User.Auth == nil {
		t.Fatal("expected User.Auth nested schema to be initialized")
	}

	if ToDoFields.User.Auth.Provider == nil || ToDoFields.User.Auth.Provider.BSONName() != "user.auth.provider" {
		t.Fatal("expected User.Auth.Provider nested field with bson user.auth.provider")
	}

	if ToDoFields.User.Auth.Scopes == nil || ToDoFields.User.Auth.Scopes.BSONName() != "user.auth.scopes" {
		t.Fatal("expected User.Auth.Scopes nested field with bson user.auth.scopes")
	}

	if ToDoFields.Tags == nil || ToDoFields.Tags.BSONName() != "tags" {
		t.Fatal("expected Tags field with bson tags")
	}

	if ToDoMetaFields.Source == nil || ToDoMetaFields.Source.BSONName() != "meta.source" {
		t.Fatal("expected Source nested field with bson meta.source")
	}

	if ToDoMetaFields.Priority == nil || ToDoMetaFields.Priority.BSONName() != "meta.priority" {
		t.Fatal("expected Priority nested field with bson meta.priority")
	}
}

func ValidateAdvancedGenericQueries(t *testing.T) {
	logger(t, "Validating GenericField advanced query methods")

	if !reflect.DeepEqual(ToDoFields.Meta.Path("source").Eq("import"), bson.M{"meta.source": "import"}) {
		t.Fatal("unexpected Generic Path + Eq query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.Contains("urgent"), bson.M{"tags": bson.M{"$in": []any{"urgent"}}}) {
		t.Fatal("unexpected Generic Contains query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.Size(2), bson.M{"tags": bson.M{"$size": 2}}) {
		t.Fatal("unexpected Generic Size query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.ContainsAll([]any{"urgent", "backend"}), bson.M{"tags": bson.M{"$all": []any{"urgent", "backend"}}}) {
		t.Fatal("unexpected Generic ContainsAll query")
	}

	if !reflect.DeepEqual(ToDoFields.Meta.Exists(), bson.M{"meta": bson.M{"$exists": true}}) {
		t.Fatal("unexpected Generic Exists query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.ElemMatch(bson.M{"$eq": "urgent"}), bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "urgent"}}}) {
		t.Fatal("unexpected Generic ElemMatch query")
	}

	if !reflect.DeepEqual(ToDoMetaFields.Source.Eq("import"), bson.M{"meta.source": "import"}) {
		t.Fatal("unexpected nested typed string field query")
	}

	if !reflect.DeepEqual(ToDoMetaFields.Priority.Gte(2), bson.M{"meta.priority": bson.M{"$gte": int64(2)}}) {
		t.Fatal("unexpected nested typed int field query")
	}

	if !reflect.DeepEqual(ToDoFields.User.Email.Eq("john@example.com"), bson.M{"user.email": "john@example.com"}) {
		t.Fatal("unexpected nested user email query")
	}

	if !reflect.DeepEqual(ToDoFields.User.ID.Eq("507f1f77bcf86cd799439011"), bson.M{"user._id": "507f1f77bcf86cd799439011"}) {
		t.Fatal("unexpected nested user id query using schema override")
	}

	if !reflect.DeepEqual(ToDoFields.User.Auth.Provider.Eq("google"), bson.M{"user.auth.provider": "google"}) {
		t.Fatal("unexpected deep nested provider query")
	}

	if !reflect.DeepEqual(ToDoFields.User.Auth.Scopes.Contains("email"), bson.M{"user.auth.scopes": bson.M{"$in": []string{"email"}}}) {
		t.Fatal("unexpected deep nested typed scopes contains query")
	}
}

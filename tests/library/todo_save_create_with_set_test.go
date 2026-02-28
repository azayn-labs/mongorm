package main

import (
	"testing"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestSaveCreateWithSetMergesSchemaFields(t *testing.T) {
	email := "set-create@example.com"
	user := &ToDo{
		Text: mongorm.String(email),
	}

	model := mongorm.New(user)
	model.Set(&ToDo{
		Done:  mongorm.Bool(false),
		Count: 1,
	})

	if err := model.Save(t.Context()); err != nil {
		t.Fatalf("expected save to create document, got: %v", err)
	}
	defer DeleteLibraryTodoByID(t, user.ID)

	if user.ID == nil || *user.ID == (bson.ObjectID{}) {
		t.Fatal("expected created document id")
	}

	if user.Text == nil || *user.Text != email {
		t.Fatalf("expected existing schema field text to be preserved, got: %v", user.Text)
	}

	if user.Done == nil || *user.Done != false {
		t.Fatalf("expected set-applied field done=false, got: %v", user.Done)
	}

	if user.Count != 1 {
		t.Fatalf("expected set-applied field count=1, got: %d", user.Count)
	}

}

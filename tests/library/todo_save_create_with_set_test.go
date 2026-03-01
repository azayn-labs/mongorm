package main

import (
	"testing"
	"time"

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

func TestSaveCreateWithSetOnInsertMergesSchemaFields(t *testing.T) {
	email := "set-on-insert-create@example.com"
	user := &ToDo{
		Text: mongorm.String(email),
	}

	model := mongorm.New(user)
	model.SetOnInsert(&ToDo{
		Done:  mongorm.Bool(true),
		Count: 9,
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

	if user.Done == nil || *user.Done != true {
		t.Fatalf("expected setOnInsert-applied field done=true, got: %v", user.Done)
	}

	if user.Count != 9 {
		t.Fatalf("expected setOnInsert-applied field count=9, got: %d", user.Count)
	}
}

func TestSaveWithSetOnInsertUpsertInsertAndMatchNoop(t *testing.T) {
	text := "set-on-insert-upsert-" + time.Now().Format(time.RFC3339Nano)

	upserted := &ToDo{}
	err := mongorm.New(upserted).
		WhereBy(ToDoFields.Text, text).
		SetOnInsert(&ToDo{
			Text:  mongorm.String(text),
			Done:  mongorm.Bool(true),
			Count: 42,
		}).
		Save(t.Context())
	if err != nil {
		t.Fatalf("expected save upsert insert to succeed, got: %v", err)
	}
	defer DeleteLibraryTodoByID(t, upserted.ID)

	if upserted.ID == nil || *upserted.ID == (bson.ObjectID{}) {
		t.Fatal("expected upserted document id")
	}

	if upserted.Done == nil || *upserted.Done != true {
		t.Fatalf("expected upserted done=true, got: %v", upserted.Done)
	}

	if upserted.Count != 42 {
		t.Fatalf("expected upserted count=42, got: %d", upserted.Count)
	}

	if upserted.Text == nil || *upserted.Text != text {
		t.Fatalf("expected upserted text=%q, got: %v", text, upserted.Text)
	}

	err = mongorm.New(&ToDo{}).
		WhereBy(ToDoFields.Text, text).
		SetOnInsert(&ToDo{
			Done:  mongorm.Bool(false),
			Count: 7,
		}).
		Save(t.Context())
	if err != nil {
		t.Fatalf("expected save with matching filter to succeed, got: %v", err)
	}

	check := &ToDo{}
	if err := mongorm.New(check).WhereBy(ToDoFields.ID, *upserted.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if check.Done == nil || *check.Done != true {
		t.Fatalf("expected matched update to keep original done=true, got: %v", check.Done)
	}

	if check.Count != 42 {
		t.Fatalf("expected matched update to keep original count=42, got: %d", check.Count)
	}
}

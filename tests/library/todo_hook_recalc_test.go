package main

import (
	"errors"
	"testing"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type HookRecalcToDo struct {
	ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
	Text *string        `bson:"text,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type HookRecalcToDoSchema struct {
	ID   *primitives.ObjectIDField
	Text *primitives.StringField
}

var HookRecalcToDoFields = mongorm.FieldsOf[HookRecalcToDo, HookRecalcToDoSchema]()
var hookCreateAfterModified bool
var hookUpdateAfterModified bool
var hookFindTargetText string
var hookDeleteTargetID bson.ObjectID

func (t *HookRecalcToDo) BeforeSave(m *mongorm.MongORM[HookRecalcToDo], _ *bson.M) error {
	m.SetData(HookRecalcToDoFields.Text, "updated-by-hook")
	return nil
}

func TestBeforeSaveCanAddUpdateOperations(t *testing.T) {
	seed := &ToDo{Text: mongorm.String("hook-recalc-seed")}
	CreateLibraryTodo(t, seed)
	defer DeleteLibraryTodoByID(t, seed.ID)

	hooked := &HookRecalcToDo{}
	err := mongorm.New(hooked).
		WhereBy(HookRecalcToDoFields.ID, *seed.ID).
		Save(t.Context())
	if err != nil {
		t.Fatalf("expected save to succeed with hook-added update operations, got: %v", err)
	}

	verify := &ToDo{}
	if err := mongorm.New(verify).WhereBy(ToDoFields.ID, *seed.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if verify.Text == nil || *verify.Text != "updated-by-hook" {
		t.Fatalf("expected text updated by hook, got: %v", verify.Text)
	}
}

type HookCreateRecalcToDo struct {
	ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
	Text *string        `bson:"text,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type HookCreateRecalcToDoSchema struct {
	ID   *primitives.ObjectIDField
	Text *primitives.StringField
}

var HookCreateRecalcToDoFields = mongorm.FieldsOf[HookCreateRecalcToDo, HookCreateRecalcToDoSchema]()

func (t *HookCreateRecalcToDo) BeforeCreate(_ *mongorm.MongORM[HookCreateRecalcToDo]) error {
	t.Text = mongorm.String("created-by-hook")
	return nil
}

func (t *HookCreateRecalcToDo) AfterCreate(m *mongorm.MongORM[HookCreateRecalcToDo]) error {
	hookCreateAfterModified = m.IsModified(HookCreateRecalcToDoFields.Text)
	return nil
}

func TestBeforeCreateRebuildsModifiedAfterHook(t *testing.T) {
	hookCreateAfterModified = false

	hooked := &HookCreateRecalcToDo{}
	if err := mongorm.New(hooked).Save(t.Context()); err != nil {
		t.Fatalf("expected save(create) to succeed, got: %v", err)
	}
	defer DeleteLibraryTodoByID(t, hooked.ID)

	if !hookCreateAfterModified {
		t.Fatal("expected AfterCreate to observe hook-updated field as modified")
	}
}

type HookUpdateRecalcToDo struct {
	ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
	Text *string        `bson:"text,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type HookUpdateRecalcToDoSchema struct {
	ID   *primitives.ObjectIDField
	Text *primitives.StringField
}

var HookUpdateRecalcToDoFields = mongorm.FieldsOf[HookUpdateRecalcToDo, HookUpdateRecalcToDoSchema]()

func (t *HookUpdateRecalcToDo) BeforeUpdate(m *mongorm.MongORM[HookUpdateRecalcToDo], _ *bson.M, _ *bson.M) error {
	m.SetData(HookUpdateRecalcToDoFields.Text, "updated-in-before-update")
	return nil
}

func (t *HookUpdateRecalcToDo) AfterUpdate(m *mongorm.MongORM[HookUpdateRecalcToDo]) error {
	hookUpdateAfterModified = m.IsModified(HookUpdateRecalcToDoFields.Text)
	return nil
}

func TestBeforeUpdateRebuildsModifiedAfterHook(t *testing.T) {
	hookUpdateAfterModified = false

	seed := &ToDo{Text: mongorm.String("before-update-seed")}
	CreateLibraryTodo(t, seed)
	defer DeleteLibraryTodoByID(t, seed.ID)

	hooked := &HookUpdateRecalcToDo{}
	err := mongorm.New(hooked).
		WhereBy(HookUpdateRecalcToDoFields.ID, *seed.ID).
		SetData(HookUpdateRecalcToDoFields.Text, "initial-update").
		Save(t.Context())
	if err != nil {
		t.Fatalf("expected save(update) to succeed, got: %v", err)
	}

	if !hookUpdateAfterModified {
		t.Fatal("expected AfterUpdate to observe hook-updated field as modified")
	}

	verify := &ToDo{}
	if err := mongorm.New(verify).WhereBy(ToDoFields.ID, *seed.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if verify.Text == nil || *verify.Text != "updated-in-before-update" {
		t.Fatalf("expected text updated in BeforeUpdate hook, got: %v", verify.Text)
	}
}

type HookFindFilterToDo struct {
	ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
	Text *string        `bson:"text,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type HookFindFilterToDoSchema struct {
	ID   *primitives.ObjectIDField
	Text *primitives.StringField
}

var HookFindFilterToDoFields = mongorm.FieldsOf[HookFindFilterToDo, HookFindFilterToDoSchema]()

func (t *HookFindFilterToDo) BeforeFind(_ *mongorm.MongORM[HookFindFilterToDo], filter *bson.M) error {
	*filter = bson.M{"text": hookFindTargetText}
	return nil
}

func TestBeforeFindCanMutateFilter(t *testing.T) {
	target := &ToDo{Text: mongorm.String("before-find-target")}
	other := &ToDo{Text: mongorm.String("before-find-other")}
	CreateLibraryTodo(t, target)
	CreateLibraryTodo(t, other)
	defer DeleteLibraryTodoByID(t, target.ID)
	defer DeleteLibraryTodoByID(t, other.ID)

	hookFindTargetText = *target.Text

	found := &HookFindFilterToDo{}
	err := mongorm.New(found).
		WhereBy(HookFindFilterToDoFields.Text, *other.Text).
		First(t.Context())
	if err != nil {
		t.Fatalf("expected find to succeed, got: %v", err)
	}

	if found.Text == nil || *found.Text != *target.Text {
		t.Fatalf("expected hook-mutated filter to find target text %q, got: %v", *target.Text, found.Text)
	}
}

type HookDeleteFilterToDo struct {
	ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
	Text *string        `bson:"text,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type HookDeleteFilterToDoSchema struct {
	ID   *primitives.ObjectIDField
	Text *primitives.StringField
}

var HookDeleteFilterToDoFields = mongorm.FieldsOf[HookDeleteFilterToDo, HookDeleteFilterToDoSchema]()

func (t *HookDeleteFilterToDo) BeforeDelete(_ *mongorm.MongORM[HookDeleteFilterToDo], filter *bson.M) error {
	*filter = bson.M{"_id": hookDeleteTargetID}
	return nil
}

func TestBeforeDeleteCanMutateFilter(t *testing.T) {
	target := &ToDo{Text: mongorm.String("before-delete-target")}
	other := &ToDo{Text: mongorm.String("before-delete-other")}
	CreateLibraryTodo(t, target)
	CreateLibraryTodo(t, other)
	defer func() {
		err := mongorm.New(&ToDo{}).WhereBy(ToDoFields.ID, *target.ID).Delete(t.Context())
		if err != nil && !errors.Is(err, mongorm.ErrNotFound) {
			t.Fatalf("cleanup delete target failed: %v", err)
		}
	}()
	defer func() {
		err := mongorm.New(&ToDo{}).WhereBy(ToDoFields.ID, *other.ID).Delete(t.Context())
		if err != nil && !errors.Is(err, mongorm.ErrNotFound) {
			t.Fatalf("cleanup delete other failed: %v", err)
		}
	}()

	hookDeleteTargetID = *target.ID

	err := mongorm.New(&HookDeleteFilterToDo{}).
		WhereBy(HookDeleteFilterToDoFields.ID, *other.ID).
		Delete(t.Context())
	if err != nil {
		t.Fatalf("expected delete to succeed, got: %v", err)
	}

	deleted := &ToDo{}
	err = mongorm.New(deleted).WhereBy(ToDoFields.ID, *target.ID).First(t.Context())
	if !errors.Is(err, mongorm.ErrNotFound) {
		t.Fatalf("expected target to be deleted by hook-mutated filter, got: %v", err)
	}

	remaining := &ToDo{}
	if err := mongorm.New(remaining).WhereBy(ToDoFields.ID, *other.ID).First(t.Context()); err != nil {
		t.Fatalf("expected non-target document to remain, got: %v", err)
	}
}

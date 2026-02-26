package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ValidateLibraryOptimisticLocking(t *testing.T) {
	text := "optimistic-lock-" + time.Now().Format(time.RFC3339Nano)

	created := &ToDo{Text: mongorm.String(text), Done: mongorm.Bool(false), Count: 1}
	if err := mongorm.New(created).Save(t.Context()); err != nil {
		t.Fatal(err)
	}
	defer DeleteLibraryTodoByID(t, created.ID)

	if created.Version != 1 {
		t.Fatalf("expected initial version=1, got %d", created.Version)
	}

	current := &ToDo{}
	if err := mongorm.New(current).WhereBy(ToDoFields.ID, *created.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	staleVersion := current.Version

	updater := mongorm.New(&ToDo{ID: created.ID, Version: staleVersion})
	updater.Set(&ToDo{Text: mongorm.String(text + "-v2")})
	if err := updater.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	staleUpdater := mongorm.New(&ToDo{ID: created.ID, Version: staleVersion})
	staleUpdater.Set(&ToDo{Text: mongorm.String(text + "-stale")})
	err := staleUpdater.Save(t.Context())
	if err == nil {
		t.Fatal("expected optimistic lock conflict for stale version")
	}

	if !errors.Is(err, mongorm.ErrOptimisticLockConflict) {
		t.Fatalf("expected ErrOptimisticLockConflict, got: %v", err)
	}
}

func ValidateLibraryErrorTaxonomy(t *testing.T) {
	notFound := &ToDo{}
	err := mongorm.New(notFound).
		WhereBy(ToDoFields.ID, bson.NewObjectID()).
		First(t.Context())
	if err == nil {
		t.Fatal("expected not found error")
	}

	if !errors.Is(err, mongorm.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Fatalf("expected mongo.ErrNoDocuments compatibility, got: %v", err)
	}

	dupID := bson.NewObjectID()
	firstInsert := mongorm.NewBulkWriteBuilder[ToDo]().
		InsertOne(&ToDo{ID: &dupID, Text: mongorm.String(fmt.Sprintf("dup-a-%d", time.Now().UnixNano()))}).
		Models()

	if _, err := mongorm.New(&ToDo{}).BulkWrite(t.Context(), firstInsert); err != nil {
		t.Fatal(err)
	}
	defer DeleteLibraryTodoByID(t, &dupID)

	secondInsert := mongorm.NewBulkWriteBuilder[ToDo]().
		InsertOne(&ToDo{ID: &dupID, Text: mongorm.String(fmt.Sprintf("dup-b-%d", time.Now().UnixNano()))}).
		Models()

	_, err = mongorm.New(&ToDo{}).BulkWrite(t.Context(), secondInsert)
	if err == nil {
		t.Fatal("expected duplicate key error")
	}

	if !errors.Is(err, mongorm.ErrDuplicateKey) {
		t.Fatalf("expected ErrDuplicateKey, got: %v", err)
	}

	type invalidConfigModel struct {
		ID *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
	}

	bad := &invalidConfigModel{}
	err = mongorm.New(bad).First(t.Context())
	if err == nil {
		t.Fatal("expected invalid configuration error")
	}

	if !errors.Is(err, mongorm.ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got: %v", err)
	}
}

package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ValidateLibraryTransactions(t *testing.T) {
	rollbackText := "tx-rollback-" + time.Now().Format(time.RFC3339Nano)
	commitText := "tx-commit-" + time.Now().Format(time.RFC3339Nano)

	rollbackErr := errors.New("rollback transaction test")

	err := mongorm.New(&ToDo{}).WithTransaction(t.Context(), func(txCtx context.Context) error {
		toDo := &ToDo{Text: mongorm.String(rollbackText), Done: mongorm.Bool(false), Count: 1}
		if err := mongorm.New(toDo).Save(txCtx); err != nil {
			return err
		}

		return rollbackErr
	})

	if mongorm.IsTransactionUnsupported(err) {
		t.Skipf("transactions unsupported by current mongodb setup: %v", err)
	}

	if err == nil {
		t.Fatal("expected rollback error from transaction callback")
	}

	if !errors.Is(err, rollbackErr) {
		t.Fatalf("expected rollback sentinel error, got: %v", err)
	}

	rollbackVerify := &ToDo{}
	err = mongorm.New(rollbackVerify).
		WhereBy(ToDoFields.Text, rollbackText).
		First(t.Context())
	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Fatalf("expected no document after rollback, got: %v", err)
	}

	err = mongorm.New(&ToDo{}).WithTransaction(t.Context(), func(txCtx context.Context) error {
		toDo := &ToDo{Text: mongorm.String(commitText), Done: mongorm.Bool(true), Count: 2}
		return mongorm.New(toDo).Save(txCtx)
	})
	if err != nil {
		t.Fatal(err)
	}

	commitVerify := &ToDo{}
	err = mongorm.New(commitVerify).
		WhereBy(ToDoFields.Text, commitText).
		First(t.Context())
	if err != nil {
		t.Fatalf("expected committed document, got: %v", err)
	}

	DeleteAllLibraryTodoByText(t, commitText)
}

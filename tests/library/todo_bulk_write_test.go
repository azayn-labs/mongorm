package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ValidateLibraryBulkWrite(t *testing.T) {
	prefix := "bulk-write-" + time.Now().Format(time.RFC3339Nano)
	txPrefix := "bulk-write-tx-" + time.Now().Format(time.RFC3339Nano)

	insertModels := mongorm.NewBulkWriteBuilder[ToDo]().
		InsertOne(&ToDo{Text: mongorm.String(prefix + "-a"), Done: mongorm.Bool(false), Count: 1}).
		InsertOne(&ToDo{Text: mongorm.String(prefix + "-b"), Done: mongorm.Bool(false), Count: 2}).
		Models()

	res, err := mongorm.New(&ToDo{}).BulkWrite(
		t.Context(),
		insertModels,
		options.BulkWrite().SetOrdered(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	if res.InsertedCount != 2 {
		t.Fatalf("expected 2 inserted docs, got %d", res.InsertedCount)
	}

	updateDeleteModels := mongorm.NewBulkWriteBuilder[ToDo]().
		UpdateMany(
			ToDoFields.Text.Reg(fmt.Sprintf("^%s", regexp.QuoteMeta(prefix))),
			bson.M{"$set": bson.M{ToDoFields.Done.BSONName(): true}},
			false,
		).
		DeleteOne(ToDoFields.Text.Eq(prefix + "-a")).
		Models()

	res, err = mongorm.New(&ToDo{}).BulkWrite(
		t.Context(),
		updateDeleteModels,
		options.BulkWrite().SetOrdered(false),
	)
	if err != nil {
		t.Fatal(err)
	}

	if res.MatchedCount < 1 {
		t.Fatalf("expected at least 1 matched doc, got %d", res.MatchedCount)
	}

	if res.DeletedCount != 1 {
		t.Fatalf("expected 1 deleted doc, got %d", res.DeletedCount)
	}

	replaceModel := mongorm.NewBulkWriteBuilder[ToDo]().
		ReplaceOne(
			ToDoFields.Text.Eq(prefix+"-b"),
			&ToDo{Text: mongorm.String(prefix + "-b"), Done: mongorm.Bool(true), Count: 99},
			false,
		).
		Models()

	res, err = mongorm.New(&ToDo{}).BulkWrite(t.Context(), replaceModel)
	if err != nil {
		t.Fatal(err)
	}

	if res.MatchedCount != 1 {
		t.Fatalf("expected 1 matched doc for replace, got %d", res.MatchedCount)
	}

	verify := &ToDo{}
	err = mongorm.New(verify).
		WhereBy(ToDoFields.Text, prefix+"-b").
		First(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if verify.Count != 99 {
		t.Fatalf("expected replaced count=99, got %d", verify.Count)
	}

	txModels := mongorm.NewBulkWriteBuilder[ToDo]().
		InsertOne(&ToDo{Text: mongorm.String(txPrefix + "-a"), Done: mongorm.Bool(false), Count: 5}).
		UpdateOne(
			ToDoFields.Text.Eq(txPrefix+"-a"),
			bson.M{"$set": bson.M{ToDoFields.Done.BSONName(): true}},
			false,
		).
		Models()

	txRes, err := mongorm.New(&ToDo{}).BulkWriteInTransaction(
		t.Context(),
		txModels,
		options.BulkWrite().SetOrdered(true),
	)
	if mongorm.IsTransactionUnsupported(err) {
		t.Skipf("transactions unsupported by current mongodb setup: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}

	if txRes.InsertedCount != 1 {
		t.Fatalf("expected 1 inserted doc in tx bulk write, got %d", txRes.InsertedCount)
	}

	txVerify := &ToDo{}
	err = mongorm.New(txVerify).
		WhereBy(ToDoFields.Text, txPrefix+"-a").
		First(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if txVerify.Done == nil || !*txVerify.Done {
		t.Fatal("expected tx bulk write update to set done=true")
	}

	DeleteAllLibraryTodoByText(t, prefix+"-b")
	DeleteAllLibraryTodoByText(t, txPrefix+"-a")
}

package main

import (
	"testing"
	"time"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type findOneAndUpdateTimestampModel struct {
	ID        *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
	Text      *string        `bson:"text,omitempty"`
	Done      *bool          `bson:"done,omitempty"`
	CreatedAt *time.Time     `bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
	UpdatedAt *time.Time     `bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

func TestFindOneAndUpdateDoesNotInjectAutoTimestampsIntoFilter(t *testing.T) {
	prefix := "find-one-and-update-ts-filter-" + time.Now().Format(time.RFC3339Nano)

	first := &findOneAndUpdateTimestampModel{Text: mongorm.String(prefix), Done: mongorm.Bool(false)}
	second := &findOneAndUpdateTimestampModel{Text: mongorm.String(prefix), Done: mongorm.Bool(false)}

	if err := mongorm.New(first).Save(t.Context()); err != nil {
		t.Fatalf("failed creating first doc: %v", err)
	}
	if err := mongorm.New(second).Save(t.Context()); err != nil {
		t.Fatalf("failed creating second doc: %v", err)
	}
	defer func() {
		if first.ID != nil {
			_ = mongorm.New(&findOneAndUpdateTimestampModel{ID: first.ID}).Delete(t.Context())
		}
		if second.ID != nil {
			_ = mongorm.New(&findOneAndUpdateTimestampModel{ID: second.ID}).Delete(t.Context())
		}
	}()

	updater := &findOneAndUpdateTimestampModel{Text: mongorm.String(prefix)}
	err := mongorm.New(updater).
		Where(bson.M{"text": prefix}).
		Set(&findOneAndUpdateTimestampModel{Done: mongorm.Bool(true)}).
		FindOneAndUpdate(
			t.Context(),
			options.FindOneAndUpdate().SetSort(bson.D{{Key: "_id", Value: 1}}),
		)
	if err != nil {
		t.Fatalf("expected FindOneAndUpdate to succeed without timestamp filter injection, got: %v", err)
	}

	updatedCount, err := mongorm.New(&findOneAndUpdateTimestampModel{}).
		Where(bson.M{"text": prefix, "done": true}).
		Count(t.Context())
	if err != nil {
		t.Fatalf("failed counting updated docs: %v", err)
	}
	if updatedCount != 1 {
		t.Fatalf("expected exactly one updated doc, got %d", updatedCount)
	}
}

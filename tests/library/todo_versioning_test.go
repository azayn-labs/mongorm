package main

import (
	"errors"
	"testing"
	"time"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidateLibraryVersioningConfigSafety(t *testing.T) {
	t.Run("Unexported version field returns invalid config", func(t *testing.T) {
		type invalidUnexportedVersionModel struct {
			ID *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`

			version int64 `bson:"_version,omitempty" mongorm:"version"`

			connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
			database         *string `mongorm:"orm-test,connection:database"`
			collection       *string `mongorm:"todo_library,connection:collection"`
		}

		id := bson.NewObjectID()
		model := &invalidUnexportedVersionModel{ID: &id}
		err := mongorm.New(model).Save(t.Context())
		if err == nil {
			t.Fatal("expected invalid configuration error for unexported version field")
		}

		if !errors.Is(err, mongorm.ErrInvalidConfig) {
			t.Fatalf("expected ErrInvalidConfig, got: %v", err)
		}
	})

	t.Run("Pointer version field initializes safely", func(t *testing.T) {
		type pointerVersionModel struct {
			ID      *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
			Text    *string        `bson:"text,omitempty"`
			Version *int64         `bson:"_version,omitempty" mongorm:"version"`

			connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
			database         *string `mongorm:"orm-test,connection:database"`
			collection       *string `mongorm:"todo_library,connection:collection"`
		}

		id := bson.NewObjectID()
		text := "pointer-version-" + time.Now().Format(time.RFC3339Nano)
		model := &pointerVersionModel{
			ID:   &id,
			Text: mongorm.String(text),
		}

		err := mongorm.New(model).Save(t.Context())
		if err != nil {
			t.Fatalf("unexpected save error: %v", err)
		}
		defer func() {
			_ = mongorm.New(&ToDo{}).
				WhereBy(ToDoFields.ID, *model.ID).
				Delete(t.Context())
		}()

		if model.Version == nil || *model.Version != 1 {
			t.Fatalf("expected initialized version=1, got: %+v", model.Version)
		}
	})
}

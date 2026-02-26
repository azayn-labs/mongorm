package main

import (
	"testing"
	"time"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidateLibraryTimestamps(t *testing.T) {
	t.Run("Only updated_at is refreshed on save", func(t *testing.T) {
		type updatedAtOnlyModel struct {
			ID        *bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" mongorm:"primary"`
			Text      *string        `bson:"text,omitempty"`
			UpdatedAt *time.Time     `bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`

			connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
			database         *string `mongorm:"orm-test,connection:database"`
			collection       *string `mongorm:"todo_library,connection:collection"`
		}

		text := "ts-updated-at-only-" + time.Now().Format(time.RFC3339Nano)
		model := &updatedAtOnlyModel{Text: mongorm.String(text)}
		if err := mongorm.New(model).Save(t.Context()); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if model.ID == nil {
				return
			}

			_ = mongorm.New(&updatedAtOnlyModel{ID: model.ID}).Delete(t.Context())
		}()

		if model.UpdatedAt == nil || model.UpdatedAt.IsZero() {
			t.Fatal("expected UpdatedAt to be set on insert")
		}

		firstUpdatedAt := *model.UpdatedAt
		time.Sleep(2 * time.Millisecond)

		updater := mongorm.New(model)
		updater.Set(&updatedAtOnlyModel{Text: mongorm.String(text + "-v2")})
		if err := updater.Save(t.Context()); err != nil {
			t.Fatal(err)
		}

		if model.UpdatedAt == nil || !model.UpdatedAt.After(firstUpdatedAt) {
			t.Fatalf("expected UpdatedAt to be refreshed; first=%v current=%v", firstUpdatedAt, model.UpdatedAt)
		}
	})

	t.Run("Only created_at is initialized on insert", func(t *testing.T) {
		type createdAtOnlyModel struct {
			ID        *bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" mongorm:"primary"`
			Text      *string        `bson:"text,omitempty"`
			CreatedAt *time.Time     `bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`

			connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
			database         *string `mongorm:"orm-test,connection:database"`
			collection       *string `mongorm:"todo_library,connection:collection"`
		}

		text := "ts-created-at-only-" + time.Now().Format(time.RFC3339Nano)
		model := &createdAtOnlyModel{Text: mongorm.String(text)}
		if err := mongorm.New(model).Save(t.Context()); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if model.ID == nil {
				return
			}

			_ = mongorm.New(&createdAtOnlyModel{ID: model.ID}).Delete(t.Context())
		}()

		if model.CreatedAt == nil || model.CreatedAt.IsZero() {
			t.Fatal("expected CreatedAt to be set on insert")
		}
	})

	t.Run("No timestamp tags does not auto-populate time fields", func(t *testing.T) {
		type noTimestampTagsModel struct {
			ID        *bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" mongorm:"primary"`
			Text      *string        `bson:"text,omitempty"`
			CreatedAt *time.Time     `bson:"createdAt,omitempty"`

			connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
			database         *string `mongorm:"orm-test,connection:database"`
			collection       *string `mongorm:"todo_library,connection:collection"`
		}

		model := &noTimestampTagsModel{Text: mongorm.String("ts-no-tags-" + time.Now().Format(time.RFC3339Nano))}
		if err := mongorm.New(model).Save(t.Context()); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if model.ID == nil {
				return
			}

			_ = mongorm.New(&noTimestampTagsModel{ID: model.ID}).Delete(t.Context())
		}()

		if model.CreatedAt != nil {
			t.Fatalf("expected CreatedAt to remain nil without timestamp tags, got: %v", model.CreatedAt)
		}
	})
}

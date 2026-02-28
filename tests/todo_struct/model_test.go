package todostruct

import (
	"fmt"
	"time"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
	ID   *bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" mongorm:"primary"`
	Text *string        `json:"text,omitempty" bson:"text,omitempty"`
	Tags []string       `json:"tags,omitempty" bson:"tags,omitempty"`
	Meta *ToDoMeta      `json:"meta,omitempty" bson:"meta,omitempty"`

	// Timestamps
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`

	// MongORM options
	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo,connection:collection"`
}

type ToDoMeta struct {
	Source   *string `json:"source,omitempty" bson:"source,omitempty"`
	Priority *int64  `json:"priority,omitempty" bson:"priority,omitempty"`
}

type ToDoMetaSchema struct {
	Source   *primitives.StringField
	Priority *primitives.Int64Field
}

type ToDoSchema struct {
	ID        *primitives.ObjectIDField
	Text      *primitives.StringField
	Tags      *primitives.GenericField
	Meta      *primitives.GenericField
	CreatedAt *primitives.TimestampField
	UpdatedAt *primitives.TimestampField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()
var ToDoMetaFields = mongorm.NestedFieldsOf[ToDoMeta, ToDoMetaSchema](ToDoFields.Meta)

// Hooks

func (t *ToDo) BeforeFind(m *mongorm.MongORM[ToDo], query *bson.M) error {
	fmt.Printf("[HOOK] before finding a document using query: %+v\n", *query)

	return nil
}

func (t *ToDo) AfterFind(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] after finding a document")

	return nil
}

func (t *ToDo) BeforeSave(m *mongorm.MongORM[ToDo], query *bson.M) error {
	if query != nil {
		fmt.Printf("[HOOK] before saving a document with query %+v\n", *query)
	} else {
		fmt.Println("[HOOK] before saving a new document")
	}

	return nil
}

func (t *ToDo) AfterSave(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] after saving a document")

	return nil
}

func (t *ToDo) BeforeCreate(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] before creating a document")

	return nil
}

func (t *ToDo) AfterCreate(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] after creating a document")

	return nil
}

func (t *ToDo) BeforeUpdate(m *mongorm.MongORM[ToDo], filter *bson.M, update *bson.M) error {
	fmt.Printf("[HOOK] before updating a document with filter: %+v and update: %+v\n", *filter, *update)

	return nil
}

func (t *ToDo) AfterUpdate(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] after updating a document")

	return nil
}

func (t *ToDo) BeforeDelete(m *mongorm.MongORM[ToDo], filter *bson.M) error {
	fmt.Printf("[HOOK] before deleting a document with filter: %+v\n", *filter)

	return nil
}

func (t *ToDo) AfterDelete(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] after deleting a document")

	return nil
}

func (t *ToDo) BeforeFinalize(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] before finalizing a document")

	return nil
}

func (t *ToDo) AfterFinalize(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] after finalizing a document")

	return nil
}

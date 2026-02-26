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

	// Timestamps
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`

	// MongORM options
	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo,connection:collection"`
}

type ToDoSchema struct {
	ID        *primitives.ObjectIDField
	Text      *primitives.StringField
	CreatedAt *primitives.TimestampField
	UpdatedAt *primitives.TimestampField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

// Hooks

func (t *ToDo) BeforeFind(m *mongorm.MongORM[ToDo], query *bson.M) error {
	fmt.Printf("[HOOK] before finding a document using query: %+v\n", *query)

	return nil
}

func (t *ToDo) AfterFind(m *mongorm.MongORM[ToDo]) error {
	fmt.Println("[HOOK] after finding a document")

	return nil
}

func (t *ToDo) BeforeSave(m *mongorm.MongORM[ToDo], id *bson.ObjectID, query *bson.M) error {
	if id != nil {
		fmt.Printf("[HOOK] before saving a document with ID %+v and query %+v\n", id.Hex(), *query)
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

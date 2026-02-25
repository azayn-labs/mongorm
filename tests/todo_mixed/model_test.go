package main

import (
	"time"

	"github.com/CdTgr/mongorm"
	"github.com/CdTgr/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
	ID   *bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" mongorm:"primary"`
	Text *string        `json:"text,omitempty" bson:"text,omitempty"`

	// Timestamps
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`

	// Mongorm options
	timestamps *bool   `mongorm:"true,timestamps"`
	database   *string `mongorm:"orm-test,connection:database"`
	collection *string `mongorm:"todo,connection:collection"`
}

type ToDoSchema struct {
	ID        *primitives.ObjectIDField
	Text      *primitives.StringField
	CreatedAt *primitives.TimestampField
	UpdatedAt *primitives.TimestampField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

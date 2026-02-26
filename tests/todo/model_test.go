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
	CreatedAt *time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`
}

type ToDoSchema struct {
	ID        *primitives.ObjectIDField
	Text      *primitives.StringField
	CreatedAt *primitives.TimestampField
	UpdatedAt *primitives.TimestampField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

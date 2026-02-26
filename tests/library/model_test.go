package main

import (
	"time"

	"github.com/CdTgr/mongorm"
	"github.com/CdTgr/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
	ID        *bson.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty" mongorm:"primary"`
	Text      *string           `json:"text,omitempty" bson:"text,omitempty"`
	Done      *bool             `json:"done,omitempty" bson:"done,omitempty"`
	Count     int64             `json:"count,omitempty" bson:"count,omitempty"`
	Location  *mongorm.GeoPoint `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt *time.Time        `json:"createdAt,omitempty" bson:"createdAt,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type ToDoSchema struct {
	ID        *primitives.ObjectIDField
	Text      *primitives.StringField
	Done      *primitives.BoolField
	Count     *primitives.Int64Field
	Location  *primitives.GeoField
	CreatedAt *primitives.TimestampField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

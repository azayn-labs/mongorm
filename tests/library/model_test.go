package main

import (
	"time"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
	ID        *bson.ObjectID    `bson:"_id,omitempty" json:"_id,omitempty" mongorm:"primary"`
	Text      *string           `json:"text,omitempty" bson:"text,omitempty"`
	User      *User             `json:"user,omitempty" bson:"user,omitempty"`
	Done      *bool             `json:"done,omitempty" bson:"done,omitempty"`
	Count     int64             `json:"count,omitempty" bson:"count,omitempty"`
	Tags      []string          `json:"tags,omitempty" bson:"tags,omitempty"`
	Meta      *ToDoMeta         `json:"meta,omitempty" bson:"meta,omitempty"`
	Version   int64             `json:"version,omitempty" bson:"_version,omitempty" mongorm:"version"`
	Location  *mongorm.GeoPoint `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt *time.Time        `json:"createdAt,omitempty" bson:"createdAt,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
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
	User      *UserSchema
	Done      *primitives.BoolField
	Count     *primitives.Int64Field
	Tags      *primitives.GenericField
	Meta      *primitives.GenericField
	Version   *primitives.Int64Field
	Location  *primitives.GeoField
	CreatedAt *primitives.TimestampField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()
var ToDoMetaFields = mongorm.NestedFieldsOf[ToDoMeta, ToDoMetaSchema](ToDoFields.Meta)

type User struct {
	ID    *bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email *string        `json:"email,omitempty" bson:"email,omitempty"`
	Auth  *UserAuth      `json:"auth,omitempty" bson:"auth,omitempty"`
}

type UserAuth struct {
	Provider *string  `json:"provider,omitempty" bson:"provider,omitempty"`
	Scopes   []string `json:"scopes,omitempty" bson:"scopes,omitempty"`
}

type UserAuthSchema struct {
	Provider *primitives.StringField
	Scopes   *primitives.StringArrayField
}

type UserSchema struct {
	ID    *primitives.StringField
	Email *primitives.StringField
	Auth  *UserAuthSchema
}

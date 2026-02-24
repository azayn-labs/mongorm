package orm

import "go.mongodb.org/mongo-driver/v2/bson"

type BaseModel struct {
	ID *bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty" mongorm:"primary"`
}

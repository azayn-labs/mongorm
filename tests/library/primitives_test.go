package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/CdTgr/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidatePrimitiveQueries(t *testing.T) {
	logger(t, "Validating primitives query methods")

	text := primitives.StringType("text")
	if !reflect.DeepEqual(text.Eq("todo"), bson.M{"text": "todo"}) {
		t.Fatal("unexpected String Eq query")
	}
	if !reflect.DeepEqual(text.Reg("^to"), bson.M{"text": bson.M{"$regex": "^to"}}) {
		t.Fatal("unexpected String Reg query")
	}

	done := primitives.BoolType("done")
	if !reflect.DeepEqual(done.Eq(true), bson.M{"done": true}) {
		t.Fatal("unexpected Bool Eq query")
	}
	if !reflect.DeepEqual(done.IsNotNull(), bson.M{"done": bson.M{"$ne": nil}}) {
		t.Fatal("unexpected Bool IsNotNull query")
	}

	count := primitives.Int64Type("count")
	if !reflect.DeepEqual(count.Gt(2), bson.M{"count": bson.M{"$gt": int64(2)}}) {
		t.Fatal("unexpected Int64 Gt query")
	}

	price := primitives.Float64Type("price")
	if !reflect.DeepEqual(price.Lte(9.5), bson.M{"price": bson.M{"$lte": 9.5}}) {
		t.Fatal("unexpected Float64 Lte query")
	}

	id := primitives.ObjectIDType("_id")
	oid := bson.NewObjectID()
	if !reflect.DeepEqual(id.Eq(oid), bson.M{"_id": oid}) {
		t.Fatal("unexpected ObjectID Eq query")
	}

	createdAt := primitives.TimestampType("createdAt")
	now := time.Now().UTC().Truncate(time.Millisecond)
	if !reflect.DeepEqual(createdAt.Gte(now), bson.M{"createdAt": bson.M{"$gte": now}}) {
		t.Fatal("unexpected Timestamp Gte query")
	}
}

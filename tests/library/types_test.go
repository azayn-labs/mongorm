package main

import (
	"testing"
	"time"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidateTypesHelpers(t *testing.T) {
	logger(t, "Validating pointer helpers")

	text := mongorm.String("hello")
	if text == nil || *text != "hello" {
		t.Fatal("failed creating string pointer")
	}

	if v := mongorm.StringVal(text); v != "hello" {
		t.Fatal("failed reading string pointer")
	}

	if v := mongorm.StringVal(nil); v != "" {
		t.Fatal("expected empty string for nil pointer")
	}

	done := mongorm.Bool(true)
	if done == nil || !*done {
		t.Fatal("failed creating bool pointer")
	}

	if v := mongorm.BoolVal(done); !v {
		t.Fatal("failed reading bool pointer")
	}

	if v := mongorm.BoolVal(nil); v {
		t.Fatal("expected false for nil bool pointer")
	}

	count := mongorm.Int64(15)
	if count == nil || *count != 15 {
		t.Fatal("failed creating int64 pointer")
	}

	if v := mongorm.Int64Val(count); v != 15 {
		t.Fatal("failed reading int64 pointer")
	}

	if v := mongorm.Int64Val(nil); v != 0 {
		t.Fatal("expected 0 for nil int64 pointer")
	}

	price := mongorm.Float64(19.99)
	if price == nil || *price != 19.99 {
		t.Fatal("failed creating float64 pointer")
	}

	if v := mongorm.Float64Val(price); v != 19.99 {
		t.Fatal("failed reading float64 pointer")
	}

	if v := mongorm.Float64Val(nil); v != 0 {
		t.Fatal("expected 0 for nil float64 pointer")
	}

	amount, err := bson.ParseDecimal128("123.45")
	if err != nil {
		t.Fatalf("failed to parse decimal128: %v", err)
	}

	amountPtr := mongorm.Decimal128(amount)
	if amountPtr == nil || *amountPtr != amount {
		t.Fatal("failed creating decimal128 pointer")
	}

	if v := mongorm.Decimal128Val(amountPtr); v != amount {
		t.Fatal("failed reading decimal128 pointer")
	}

	if v := mongorm.Decimal128Val(nil); v != (bson.Decimal128{}) {
		t.Fatal("expected zero decimal128 for nil pointer")
	}

	now := time.Now().UTC().Truncate(time.Millisecond)
	stamp := mongorm.Timestamp(now)
	if stamp == nil || !stamp.Equal(now) {
		t.Fatal("failed creating timestamp pointer")
	}

	if v := mongorm.TimestampVal(stamp); !v.Equal(now) {
		t.Fatal("failed reading timestamp pointer")
	}

	if v := mongorm.TimestampVal(nil); !v.IsZero() {
		t.Fatal("expected zero time for nil timestamp pointer")
	}
}

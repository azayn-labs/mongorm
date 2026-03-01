package main

import (
	"reflect"
	"testing"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type decimalAutoModel struct {
	Amount bson.Decimal128 `bson:"amount,omitempty"`
}

type decimalAutoSchema struct {
	Amount *primitives.Decimal128Field
}

type decimalExplicitModel struct {
	Amount float64 `bson:"amount,omitempty"`
}

type decimalExplicitSchema struct {
	Amount *primitives.Decimal128Field
}

func TestFieldsOfMapsDecimal128ModelToDecimal128Field(t *testing.T) {
	fields := mongorm.FieldsOf[decimalAutoModel, decimalAutoSchema]()
	if fields.Amount == nil {
		t.Fatal("expected decimal128 field to be mapped")
	}

	if fields.Amount.BSONName() != "amount" {
		t.Fatalf("expected BSON name amount, got %s", fields.Amount.BSONName())
	}

	amount, err := bson.ParseDecimal128("10.01")
	if err != nil {
		t.Fatalf("failed to parse decimal128: %v", err)
	}

	query := fields.Amount.Gte(amount)
	if !reflect.DeepEqual(query, bson.M{"amount": bson.M{"$gte": amount}}) {
		t.Fatalf("unexpected decimal128 query: %#v", query)
	}
}

func TestFieldsOfSupportsExplicitDecimal128SchemaType(t *testing.T) {
	fields := mongorm.FieldsOf[decimalExplicitModel, decimalExplicitSchema]()
	if fields.Amount == nil {
		t.Fatal("expected explicit decimal128 schema field to be mapped")
	}

	amount, err := bson.ParseDecimal128("20.5")
	if err != nil {
		t.Fatalf("failed to parse decimal128: %v", err)
	}

	query := fields.Amount.Lt(amount)
	if !reflect.DeepEqual(query, bson.M{"amount": bson.M{"$lt": amount}}) {
		t.Fatalf("unexpected explicit decimal128 query: %#v", query)
	}
}

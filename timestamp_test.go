package mongorm

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type onlyUpdatedTimestampModel struct {
	UpdatedAt *time.Time `bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`
}

type onlyCreatedTimestampModel struct {
	CreatedAt *time.Time `bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
}

type noTimestampModel struct {
	Name *string `bson:"name,omitempty"`
}

func TestSetTimestampRequirementsFromSchema_EnablesWithSingleTimestampField(t *testing.T) {
	orm := &MongORM[onlyUpdatedTimestampModel]{
		schema:  &onlyUpdatedTimestampModel{},
		options: &MongORMOptions{},
	}

	if err := orm.setTimestampRequirementsFromSchema(); err != nil {
		t.Fatalf("setTimestampRequirementsFromSchema failed: %v", err)
	}

	if !orm.options.Timestamps {
		t.Fatal("expected timestamps option to be enabled for model with one timestamp field")
	}
}

func TestSetTimestampRequirementsFromSchema_DoesNotEnableWithoutTimestampFields(t *testing.T) {
	orm := &MongORM[noTimestampModel]{
		schema:  &noTimestampModel{},
		options: &MongORMOptions{},
	}

	if err := orm.setTimestampRequirementsFromSchema(); err != nil {
		t.Fatalf("setTimestampRequirementsFromSchema failed: %v", err)
	}

	if orm.options.Timestamps {
		t.Fatal("expected timestamps option to remain disabled when no timestamp fields exist")
	}
}

func TestApplyTimestamps_WithOnlyUpdatedAt(t *testing.T) {
	orm := &MongORM[onlyUpdatedTimestampModel]{
		schema:  &onlyUpdatedTimestampModel{},
		options: &MongORMOptions{Timestamps: true},
		operations: &MongORMOperations{
			update: bson.M{},
		},
	}

	orm.applyTimestamps()

	if orm.schema.UpdatedAt == nil || orm.schema.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}

	setDoc, ok := orm.operations.update["$set"].(bson.M)
	if !ok {
		t.Fatal("expected $set update document to be created")
	}
	if _, ok := setDoc["updatedAt"]; !ok {
		t.Fatal("expected updatedAt to be present in $set update document")
	}
}

func TestApplyTimestamps_WithOnlyCreatedAt(t *testing.T) {
	orm := &MongORM[onlyCreatedTimestampModel]{
		schema:  &onlyCreatedTimestampModel{},
		options: &MongORMOptions{Timestamps: true},
		operations: &MongORMOperations{
			update: bson.M{},
		},
	}

	orm.applyTimestamps()

	if orm.schema.CreatedAt == nil || orm.schema.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}

	if _, ok := orm.operations.update["$set"]; ok {
		t.Fatal("did not expect $set update document when UpdatedAt field is absent")
	}
}

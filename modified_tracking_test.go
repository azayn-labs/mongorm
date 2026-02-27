package mongorm

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type modifiedProfile struct {
	Provider *string `bson:"provider,omitempty"`
}

type modifiedArrayItem struct {
	Name *string `bson:"name,omitempty"`
}

type modifiedModel struct {
	ID      *bson.ObjectID       `bson:"_id,omitempty" mongorm:"primary"`
	Email   *string              `bson:"email,omitempty"`
	Profile *modifiedProfile     `bson:"profile,omitempty"`
	Items   *[]modifiedArrayItem `bson:"items,omitempty"`
}

func newModifiedTrackingModel() *MongORM[modifiedModel] {
	return &MongORM[modifiedModel]{
		schema: &modifiedModel{},
		info: &MongORMInfo{
			fields: buildFields[modifiedModel](),
		},
		options: &MongORMOptions{},
		operations: &MongORMOperations{
			query:  bson.M{},
			update: bson.M{},
		},
		modified: map[string]struct{}{},
	}
}

func TestIsModifiedTracksSetField(t *testing.T) {
	m := newModifiedTrackingModel()
	m.Set(&modifiedModel{Email: String("john@example.com")})

	if !m.IsModified("email") {
		t.Fatal("expected email to be marked as modified")
	}

	if !m.IsModifed("email") {
		t.Fatal("expected IsModifed alias to return true")
	}
}

func TestIsModifiedMatchesNestedPathWhenParentUpdated(t *testing.T) {
	m := newModifiedTrackingModel()
	m.Set(&modifiedModel{Profile: &modifiedProfile{Provider: String("google")}})

	if !m.IsModified("profile") {
		t.Fatal("expected profile to be marked as modified")
	}

	if !m.IsModified("profile.provider") {
		t.Fatal("expected nested profile.provider to be treated as modified when profile is updated")
	}
}

func TestModifiedFieldsFromSchemaIncludesDeepNestedValues(t *testing.T) {
	m := newModifiedTrackingModel()
	m.schema = &modifiedModel{
		Email: String("john@example.com"),
		Profile: &modifiedProfile{
			Provider: String("github"),
		},
		Items: &[]modifiedArrayItem{{Name: String("item-a")}},
	}

	m.rebuildModifiedFromSchema()

	if !m.IsModified("email") {
		t.Fatal("expected email to be marked as modified")
	}

	if !m.IsModified("profile.provider") {
		t.Fatal("expected profile.provider to be marked as modified")
	}

	if !m.IsModified("items") {
		t.Fatal("expected items to be marked as modified")
	}
}

func TestModifiedFieldsFromUpdateTracksDottedPaths(t *testing.T) {
	m := newModifiedTrackingModel()
	m.rebuildModifiedFromUpdate(bson.M{
		"$set": bson.M{
			"profile.provider": "google",
		},
	})

	if !m.IsModified("profile.provider") {
		t.Fatal("expected profile.provider to be marked as modified")
	}

	if !m.IsModified("profile") {
		t.Fatal("expected parent profile to be treated as modified")
	}
}

func TestModifiedFieldsSortedOutput(t *testing.T) {
	m := newModifiedTrackingModel()
	m.markModified("z")
	m.markModified("a")

	if !reflect.DeepEqual(m.ModifiedFields(), []string{"a", "z"}) {
		t.Fatal("expected sorted modified fields")
	}
}

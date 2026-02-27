package mongorm

import (
	"reflect"
	"testing"

	"github.com/azayn-labs/mongorm/primitives"
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

type modifiedProfileSchema struct {
	Provider *primitives.StringField
}

type modifiedModelSchema struct {
	Email   *primitives.StringField
	Profile *modifiedProfileSchema
}

var modifiedFields = FieldsOf[modifiedModel, modifiedModelSchema]()

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

func TestSetDataTracksDirectField(t *testing.T) {
	m := newModifiedTrackingModel()
	m.SetData(modifiedFields.Email, "john@example.com")

	set, ok := m.operations.update["$set"].(bson.M)
	if !ok {
		t.Fatal("expected $set update document")
	}

	if !reflect.DeepEqual(set, bson.M{"email": "john@example.com"}) {
		t.Fatalf("unexpected $set payload: %#v", set)
	}

	if !m.IsModified("email") {
		t.Fatal("expected email to be marked as modified")
	}
}

func TestSetDataTracksNestedField(t *testing.T) {
	m := newModifiedTrackingModel()
	m.SetData(modifiedFields.Profile.Provider, "google")

	set, ok := m.operations.update["$set"].(bson.M)
	if !ok {
		t.Fatal("expected $set update document")
	}

	if !reflect.DeepEqual(set, bson.M{"profile.provider": "google"}) {
		t.Fatalf("unexpected $set payload: %#v", set)
	}

	if !m.IsModified("profile.provider") {
		t.Fatal("expected profile.provider to be marked as modified")
	}

	if !m.IsModified("profile") {
		t.Fatal("expected parent profile to be treated as modified")
	}
}

func TestUnsetDataTracksDirectField(t *testing.T) {
	m := newModifiedTrackingModel()
	m.UnsetData(modifiedFields.Email)

	unset, ok := m.operations.update["$unset"].(bson.M)
	if !ok {
		t.Fatal("expected $unset update document")
	}

	if !reflect.DeepEqual(unset, bson.M{"email": 1}) {
		t.Fatalf("unexpected $unset payload: %#v", unset)
	}

	if !m.IsModified("email") {
		t.Fatal("expected email to be marked as modified")
	}
}

func TestUnsetDataTracksNestedField(t *testing.T) {
	m := newModifiedTrackingModel()
	m.UnsetData(modifiedFields.Profile.Provider)

	unset, ok := m.operations.update["$unset"].(bson.M)
	if !ok {
		t.Fatal("expected $unset update document")
	}

	if !reflect.DeepEqual(unset, bson.M{"profile.provider": 1}) {
		t.Fatalf("unexpected $unset payload: %#v", unset)
	}

	if !m.IsModified("profile.provider") {
		t.Fatal("expected profile.provider to be marked as modified")
	}

	if !m.IsModified("profile") {
		t.Fatal("expected parent profile to be treated as modified")
	}
}

func TestSetDataSupportsPositionalArrayPath(t *testing.T) {
	m := newModifiedTrackingModel()

	path := FieldPath(PositionalFiltered(RawField("items"), "item"), "name")
	m.SetData(path, "updated-name")

	set, ok := m.operations.update["$set"].(bson.M)
	if !ok {
		t.Fatal("expected $set update document")
	}

	expected := bson.M{"items.$[item].name": "updated-name"}
	if !reflect.DeepEqual(set, expected) {
		t.Fatalf("unexpected positional $set payload: %#v", set)
	}

	if !m.IsModified("items.$[item].name") {
		t.Fatal("expected positional nested path to be marked as modified")
	}
}

func TestUnsetDataSupportsPositionalArrayPath(t *testing.T) {
	m := newModifiedTrackingModel()

	path := FieldPath(Positional(RawField("items")), "name")
	m.UnsetData(path)

	unset, ok := m.operations.update["$unset"].(bson.M)
	if !ok {
		t.Fatal("expected $unset update document")
	}

	expected := bson.M{"items.$.name": 1}
	if !reflect.DeepEqual(unset, expected) {
		t.Fatalf("unexpected positional $unset payload: %#v", unset)
	}

	if !m.IsModified("items.$.name") {
		t.Fatal("expected positional unset path to be marked as modified")
	}
}

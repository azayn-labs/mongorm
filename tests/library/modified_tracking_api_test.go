package main

import (
	"reflect"
	"testing"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type trackingProfile struct {
	Provider *string `bson:"provider,omitempty"`
}

type trackingArrayItem struct {
	Name *string `bson:"name,omitempty"`
}

type trackingModel struct {
	ID      *bson.ObjectID       `bson:"_id,omitempty" mongorm:"primary"`
	Email   *string              `bson:"email,omitempty"`
	Secret  *string              `bson:"secret,omitempty" mongorm:"readonly"`
	Profile *trackingProfile     `bson:"profile,omitempty"`
	Items   *[]trackingArrayItem `bson:"items,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type trackingProfileSchema struct {
	Provider *primitives.StringField
}

type trackingModelSchema struct {
	ID      *primitives.ObjectIDField
	Email   *primitives.StringField
	Secret  *primitives.StringField
	Profile *trackingProfileSchema
}

var TrackingFields = mongorm.FieldsOf[trackingModel, trackingModelSchema]()

func newTrackingORM() *mongorm.MongORM[trackingModel] {
	return mongorm.New(&trackingModel{})
}

func TestModifiedTracksSetFieldAndAlias(t *testing.T) {
	m := newTrackingORM()
	m.Set(&trackingModel{Email: mongorm.String("john@example.com")})

	if !m.IsModified("email") {
		t.Fatal("expected email to be marked as modified")
	}

	if !m.IsModifed("email") {
		t.Fatal("expected IsModifed alias to return true")
	}
}

func TestModifiedMatchesNestedPathWhenParentUpdated(t *testing.T) {
	m := newTrackingORM()
	m.Set(&trackingModel{Profile: &trackingProfile{Provider: mongorm.String("google")}})

	if !m.IsModified("profile") {
		t.Fatal("expected profile to be marked as modified")
	}

	if !m.IsModified("profile.provider") {
		t.Fatal("expected nested profile.provider to be treated as modified when profile is updated")
	}
}

func TestModifiedFieldsSortedOutput(t *testing.T) {
	m := newTrackingORM()
	m.SetData(mongorm.RawField("z"), 1)
	m.SetData(mongorm.RawField("a"), 2)

	if !reflect.DeepEqual(m.ModifiedFields(), []string{"a", "z"}) {
		t.Fatal("expected sorted modified fields")
	}
}

func TestSetDataAndUnsetDataTrackNestedField(t *testing.T) {
	m := newTrackingORM()

	m.SetData(TrackingFields.Profile.Provider, "google")
	if !m.IsModified("profile.provider") {
		t.Fatal("expected profile.provider to be marked as modified")
	}
	if !m.IsModified("profile") {
		t.Fatal("expected parent profile to be treated as modified")
	}

	m.UnsetData(TrackingFields.Profile.Provider)
	if !m.IsModified("profile.provider") {
		t.Fatal("expected unset nested field to remain marked as modified")
	}
}

func TestSetUnsetDataSkipPrimaryAndReadonly(t *testing.T) {
	m := newTrackingORM()

	m.SetData(TrackingFields.ID, bson.NewObjectID())
	m.SetData(TrackingFields.Secret, "top-secret")
	m.UnsetData(TrackingFields.ID)
	m.UnsetData(TrackingFields.Secret)

	if m.IsModified("_id") || m.IsModified("secret") {
		t.Fatal("expected protected fields not to be marked as modified")
	}
}

func TestSetUnsetDataSupportPositionalArrayPath(t *testing.T) {
	m := newTrackingORM()

	setPath := mongorm.FieldPath(mongorm.PositionalFiltered(mongorm.RawField("items"), "item"), "name")
	m.SetData(setPath, "updated-name")
	if !m.IsModified("items.$[item].name") {
		t.Fatal("expected positional set path to be marked as modified")
	}

	unsetPath := mongorm.FieldPath(mongorm.Positional(mongorm.RawField("items")), "name")
	m.UnsetData(unsetPath)
	if !m.IsModified("items.$.name") {
		t.Fatal("expected positional unset path to be marked as modified")
	}
}

func TestIncAndDecrementDataTrackNestedField(t *testing.T) {
	m := newTrackingORM()

	m.IncData(TrackingFields.Profile.Provider, 1)
	if !m.IsModified("profile.provider") {
		t.Fatal("expected profile.provider to be marked as modified after inc")
	}
	if !m.IsModified("profile") {
		t.Fatal("expected parent profile to be treated as modified after inc")
	}

	m.DecData(TrackingFields.Profile.Provider, 2)
	if !m.IsModified("profile.provider") {
		t.Fatal("expected profile.provider to stay marked as modified after decrement")
	}
}

func TestIncAndDecrementDataSkipPrimaryAndReadonly(t *testing.T) {
	m := newTrackingORM()

	m.IncData(TrackingFields.ID, 1)
	m.IncData(TrackingFields.Secret, 1)
	m.DecData(TrackingFields.ID, 1)
	m.DecData(TrackingFields.Secret, 1)

	if m.IsModified("_id") || m.IsModified("secret") {
		t.Fatal("expected protected fields not to be marked as modified")
	}
}

func TestArrayUpdateDataTrackFieldPaths(t *testing.T) {
	m := newTrackingORM()

	itemsPath := mongorm.RawField("items")
	if itemsPath == nil {
		t.Fatal("expected non-nil items path")
	}

	m.PushData(itemsPath, bson.M{"name": "first"})
	m.PushEachData(itemsPath, []any{bson.M{"name": "second"}, bson.M{"name": "third"}})
	m.AddToSetData(itemsPath, bson.M{"name": "unique"})
	m.AddToSetEachData(itemsPath, []any{"a", "b"})
	m.PullData(itemsPath, bson.M{"name": "first"})
	m.PopFirstData(itemsPath)
	m.PopLastData(itemsPath)

	if !m.IsModified("items") {
		t.Fatal("expected items to be marked as modified by array update APIs")
	}
}

func TestArrayUpdateDataSkipPrimaryAndReadonly(t *testing.T) {
	m := newTrackingORM()

	m.PushData(TrackingFields.ID, 1)
	m.PushData(TrackingFields.Secret, 1)
	m.AddToSetData(TrackingFields.ID, 1)
	m.AddToSetData(TrackingFields.Secret, 1)
	m.PullData(TrackingFields.ID, 1)
	m.PullData(TrackingFields.Secret, 1)
	m.PopData(TrackingFields.ID, 1)
	m.PopData(TrackingFields.Secret, -1)

	if m.IsModified("_id") || m.IsModified("secret") {
		t.Fatal("expected protected fields not to be marked as modified")
	}
}

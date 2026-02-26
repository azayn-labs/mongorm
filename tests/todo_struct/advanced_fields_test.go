package todostruct

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidateAdvancedFieldsOf(t *testing.T) {
	logger(t, "Validating advanced FieldsOf schema")

	if ToDoFields.Meta == nil || ToDoFields.Meta.BSONName() != "meta" {
		t.Fatal("expected Meta field with bson meta")
	}

	if ToDoFields.Tags == nil || ToDoFields.Tags.BSONName() != "tags" {
		t.Fatal("expected Tags field with bson tags")
	}

	if ToDoMetaFields.Source == nil || ToDoMetaFields.Source.BSONName() != "meta.source" {
		t.Fatal("expected Source nested field with bson meta.source")
	}

	if ToDoMetaFields.Priority == nil || ToDoMetaFields.Priority.BSONName() != "meta.priority" {
		t.Fatal("expected Priority nested field with bson meta.priority")
	}
}

func ValidateAdvancedGenericQueries(t *testing.T) {
	logger(t, "Validating GenericField advanced query methods")

	if !reflect.DeepEqual(ToDoFields.Meta.Path("source").Eq("import"), bson.M{"meta.source": "import"}) {
		t.Fatal("unexpected Generic Path + Eq query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.Contains("urgent"), bson.M{"tags": bson.M{"$in": []any{"urgent"}}}) {
		t.Fatal("unexpected Generic Contains query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.Size(2), bson.M{"tags": bson.M{"$size": 2}}) {
		t.Fatal("unexpected Generic Size query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.ContainsAll([]any{"urgent", "backend"}), bson.M{"tags": bson.M{"$all": []any{"urgent", "backend"}}}) {
		t.Fatal("unexpected Generic ContainsAll query")
	}

	if !reflect.DeepEqual(ToDoFields.Meta.Exists(), bson.M{"meta": bson.M{"$exists": true}}) {
		t.Fatal("unexpected Generic Exists query")
	}

	if !reflect.DeepEqual(ToDoFields.Tags.ElemMatch(bson.M{"$eq": "urgent"}), bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "urgent"}}}) {
		t.Fatal("unexpected Generic ElemMatch query")
	}

	if !reflect.DeepEqual(ToDoMetaFields.Source.Eq("import"), bson.M{"meta.source": "import"}) {
		t.Fatal("unexpected nested typed string field query")
	}

	if !reflect.DeepEqual(ToDoMetaFields.Priority.Gte(2), bson.M{"meta.priority": bson.M{"$gte": int64(2)}}) {
		t.Fatal("unexpected nested typed int field query")
	}
}

package main

import (
	"strings"
	"testing"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestGetRawQueryReturnsCopy(t *testing.T) {
	model := mongorm.New(&ToDo{})
	model.WhereBy(ToDoFields.Count, int64(7))

	rawQuery := model.GetRawQuery()
	encodedQuery, err := bson.MarshalExtJSON(rawQuery, true, false)
	if err != nil {
		t.Fatalf("expected query to be encodable, got: %v", err)
	}
	if !strings.Contains(string(encodedQuery), "\"$and\"") {
		t.Fatalf("expected $and query conditions, got: %s", string(encodedQuery))
	}

	rawQuery["$and"] = bson.A{bson.M{"count": int64(999)}}
	freshQuery := model.GetRawQuery()
	freshEncodedQuery, err := bson.MarshalExtJSON(freshQuery, true, false)
	if err != nil {
		t.Fatalf("expected internal query to be encodable, got: %v", err)
	}

	jsonQuery := string(freshEncodedQuery)
	if strings.Contains(jsonQuery, "999") {
		t.Fatalf("expected internal query copy isolation, got: %s", jsonQuery)
	}
	if !strings.Contains(jsonQuery, "\"7\"") {
		t.Fatalf("expected original count filter value in internal query, got: %s", jsonQuery)
	}
}

func TestGetRawUpdateReturnsCopy(t *testing.T) {
	model := mongorm.New(&ToDo{})
	model.Set(&ToDo{Count: 5})

	rawUpdate := model.GetRawUpdate()
	encodedUpdate, err := bson.MarshalExtJSON(rawUpdate, true, false)
	if err != nil {
		t.Fatalf("expected update to be encodable, got: %v", err)
	}
	if !strings.Contains(string(encodedUpdate), "\"$set\"") {
		t.Fatalf("expected $set in update doc, got: %s", string(encodedUpdate))
	}

	rawUpdate["$set"] = bson.M{"count": int64(999)}
	freshUpdate := model.GetRawUpdate()
	freshEncodedUpdate, err := bson.MarshalExtJSON(freshUpdate, true, false)
	if err != nil {
		t.Fatalf("expected internal update to be encodable, got: %v", err)
	}

	jsonUpdate := string(freshEncodedUpdate)
	if strings.Contains(jsonUpdate, "999") {
		t.Fatalf("expected internal update copy isolation, got: %s", jsonUpdate)
	}
	if !strings.Contains(jsonUpdate, "\"5\"") {
		t.Fatalf("expected original count update value in internal update, got: %s", jsonUpdate)
	}
}

func TestGetResolvedRawQueryIncludesSchemaFields(t *testing.T) {
	done := false
	model := mongorm.New(&ToDo{Text: mongorm.String("job-a")})
	model.WhereBy(ToDoFields.Count, int64(7)).WhereBy(ToDoFields.Done, done)

	resolved, err := model.GetResolvedRawQuery()
	if err != nil {
		t.Fatalf("expected resolved query without error, got: %v", err)
	}

	encoded, err := bson.MarshalExtJSON(resolved, true, false)
	if err != nil {
		t.Fatalf("expected resolved query to be encodable, got: %v", err)
	}

	jsonQuery := string(encoded)
	if !strings.Contains(jsonQuery, "\"text\":\"job-a\"") {
		t.Fatalf("expected schema field text to be present in resolved query, got: %s", jsonQuery)
	}
	if !strings.Contains(jsonQuery, "\"count\"") {
		t.Fatalf("expected where field count to be present in resolved query, got: %s", jsonQuery)
	}
}

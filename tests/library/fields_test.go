package main

import "testing"

func ValidateFieldsOf(t *testing.T) {
	logger(t, "Validating FieldsOf schema")

	if ToDoFields.ID == nil || ToDoFields.ID.BSONName() != "_id" {
		t.Fatal("expected ID field with bson _id")
	}

	if ToDoFields.Text == nil || ToDoFields.Text.BSONName() != "text" {
		t.Fatal("expected Text field with bson text")
	}

	if ToDoFields.Done == nil || ToDoFields.Done.BSONName() != "done" {
		t.Fatal("expected Done field with bson done")
	}

	if ToDoFields.Count == nil || ToDoFields.Count.BSONName() != "count" {
		t.Fatal("expected Count field with bson count")
	}

	if ToDoFields.Location == nil || ToDoFields.Location.BSONName() != "location" {
		t.Fatal("expected Location field with bson location")
	}

	if ToDoFields.CreatedAt == nil || ToDoFields.CreatedAt.BSONName() != "createdAt" {
		t.Fatal("expected CreatedAt field with bson createdAt")
	}
}

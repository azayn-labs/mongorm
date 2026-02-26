package todostruct

import (
	"fmt"
	"testing"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UpdateToDoByID(t *testing.T, id *bson.ObjectID, update *ToDo) {
	logger(t, fmt.Sprintf("[TODO] Using id %s for update\n", id.Hex()))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.ID.Eq(*id)).Set(update)

	if err := todoModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, fmt.Sprintf("[TODO] Updated with ID %s: %+v\n", id.Hex(), toDo))
}

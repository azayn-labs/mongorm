package todostruct

import (
	"fmt"
	"testing"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func FindByIDToDo(t *testing.T, id *bson.ObjectID) {
	logger(t, fmt.Sprintf("[TODO] Finding by id %s\n", id.Hex()))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.ID.Eq(*id))

	if err := todoModel.First(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, fmt.Sprintf("[TODO] Found using ID %s: %+v\n", id.Hex(), toDo))
}
